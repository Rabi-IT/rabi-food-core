package gateway

import "net/http"

type GoTrueGatewayAdapter struct {
	baseURL    string
	serviceKey string
	httpClient *http.Client
}

func NewGoTrue(baseURL, serviceKey string) AuthGateway {
	return &GoTrueGatewayAdapter{
		baseURL:    baseURL,
		serviceKey: serviceKey,
		httpClient: &http.Client{},
	}
}
