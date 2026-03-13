package fixtures

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func DefaultHTTP(t *testing.T) *httpexpect.Expect {
	t.Helper()

	c := httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  AppURL,
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})

	return c
}
