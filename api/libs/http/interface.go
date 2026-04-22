package http

// HTTPServer defines the interface for an HTTP server with start and stop capabilities.
type HTTPServer interface {
	Start() error
	Stop() error
}
