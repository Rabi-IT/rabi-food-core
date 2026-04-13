package fixtures

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var tables = []string{
	"commerce.orders",
	"subscription.subscription_deliveries",
	"subscription.subscription_configs",
	"subscription.subscriptions",
	"catalog.products",
	"catalog.categories",
	"iam.tenant_customers",
	"iam.tenants",
	"auth.users",
}

func CleanDatabase(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		_, err := testDB.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE %s CASCADE", table))
		require.NoError(t, err)
	}
}
