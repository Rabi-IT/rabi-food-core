package gateway

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type GoTrueGatewayAdapter struct {
	baseURL    string
	serviceKey string
	httpClient *http.Client
	db         *database.PgxAdapter
}

func NewGoTrue(baseURL, serviceKey string, db *database.PgxAdapter) AuthGateway {
	return &GoTrueGatewayAdapter{
		baseURL:    baseURL,
		serviceKey: serviceKey,
		httpClient: &http.Client{},
		db:         db,
	}
}

func wrapResponseError(base error, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: status=%d body_read_error=%w", base, resp.StatusCode, err)
	}

	bodyText := strings.TrimSpace(string(body))
	if bodyText == "" {
		return fmt.Errorf("%w: status=%d", base, resp.StatusCode)
	}

	return fmt.Errorf("%w: status=%d body=%s", base, resp.StatusCode, bodyText)
}
