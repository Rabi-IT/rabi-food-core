package docproxy

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// Forward executa a requisição HTTP para o endpoint real e devolve status, body e erro.
// headers pode ser nil; Content-Type: application/json é definido quando body != nil.
func Forward(
	ctx context.Context,
	client *http.Client,
	baseURL, method, path string,
	body []byte,
	headers map[string]string,
) (statusCode int, respBody []byte, err error) {
	if client == nil {
		client = http.DefaultClient
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, respBody, nil
}

// Op descreve como montar a requisição e a resposta para uma operação da doc.
// Use com Wrap para que o "Try it" chame o endpoint real em vez do usecase.
type Op[Input, Output any] struct {
	Method string
	Path   string
	// PathFunc sobrescreve Path quando definido (ex.: "/category/"+input.ID para PATCH/GET).
	PathFunc func(Input) string
	// Body extrai o body da requisição a partir do input (nil = sem body).
	Body func(Input) ([]byte, error)
	// BuildOutput monta o output a partir do status e do body da resposta.
	BuildOutput func(statusCode int, body []byte) (Output, error)
	// BuildError é chamado quando statusCode não é 2xx; o erro é retornado pelo handler.
	BuildError func(statusCode int, body []byte) error
}

// Wrap devolve um handler que, quando baseURL != "", faz o request para o endpoint;
// caso contrário chama fallback. Assim a doc pode "chamar o endpoint" sem repetir lógica.
func Wrap[Input, Output any](
	baseURL string,
	op Op[Input, Output],
	fallback func(context.Context, Input) (Output, error),
) func(context.Context, Input) (Output, error) {
	return func(ctx context.Context, input Input) (Output, error) {
		var zero Output
		if baseURL == "" {
			return fallback(ctx, input)
		}

		var body []byte
		if op.Body != nil {
			var err error
			body, err = op.Body(input)
			if err != nil {
				return zero, err
			}
		}

		path := op.Path
		if op.PathFunc != nil {
			path = op.PathFunc(input)
		}
		statusCode, respBody, err := Forward(ctx, nil, baseURL, op.Method, path, body, nil)
		if err != nil {
			return zero, err
		}

		out, err := op.BuildOutput(statusCode, respBody)
		if err != nil {
			return zero, err
		}

		if statusCode < 200 || statusCode >= 300 {
			return zero, op.BuildError(statusCode, respBody)
		}
		return out, nil
	}
}
