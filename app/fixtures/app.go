package fixtures

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Rabi-IT/rabi-food-core/config"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/http"

	"github.com/stretchr/testify/require"
)

var (
	apiHost = "localhost:" + config.ApiPort
	ApiURL  = "http://" + apiHost
)

type Api struct {
	http     http.HTTPServer
	database database.Database
}

func NewApi() *Api {
	time.Local = time.UTC

	return &Api{
		http:     testHTTPServer,
		database: testDB,
	}
}

func (a *Api) Start(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	err := a.database.Start(ctx)
	if err != nil {
		require.NoError(t, err, "could not start the database")
	}

	go func() {
		err := a.http.Start()
		if err != nil {
			panic(fmt.Sprintf("Could not start the server: %v", err))
		}
	}()

	err = waitForServer()
	require.NoError(t, err)
}

func (a *Api) Stop(t *testing.T) {
	t.Helper()
	err := a.http.Stop()
	require.NoError(t, err)

	err = a.database.Stop()
	require.NoError(t, err)
}

//nolint:mnd
func waitForServer() error {
	timeout := 15 * time.Second
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", apiHost, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()

			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	//nolint:err113
	return fmt.Errorf("server %s did not start within %s", apiHost, timeout)
}
