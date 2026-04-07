package fixtures

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var tables = []string{
	"users",
	"tenants",
	"categories",
	"products",
	"orders",
	"subscription_deliveries",
	"subscription_configs",
	"subscriptions",
}

func CleanDatabase(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		_, err := testDB.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE %s CASCADE", table))
		require.NoError(t, err)
	}
}
