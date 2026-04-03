package fixtures

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/di"
	"rabi-food-core/libs/http"

	"github.com/samber/do"
)

var (
	testInjector   = di.NewTest()
	testDB         = do.MustInvoke[*database.GormAdapter](testInjector)
	testHTTPServer = do.MustInvoke[http.HTTPServer](testInjector)
)
