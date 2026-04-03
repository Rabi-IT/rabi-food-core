package fixtures

import (
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/di"
	"github.com/Rabi-IT/rabi-food-core/libs/http"

	"github.com/samber/do"
)

var (
	testInjector   = di.NewTest()
	testDB         = do.MustInvoke[*database.GormAdapter](testInjector)
	testHTTPServer = do.MustInvoke[http.HTTPServer](testInjector)
)
