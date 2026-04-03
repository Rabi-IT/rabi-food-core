package fixtures

import (
	"context"
	"testing"

	category_model "github.com/Rabi-IT/rabi-food-core/features/category/model"
	order_model "github.com/Rabi-IT/rabi-food-core/features/order/model"
	product_model "github.com/Rabi-IT/rabi-food-core/features/product/model"
	subscription_model "github.com/Rabi-IT/rabi-food-core/features/subscription/model"
	tenant_model "github.com/Rabi-IT/rabi-food-core/features/tenant/model"
	user_model "github.com/Rabi-IT/rabi-food-core/features/user/model"

	"github.com/stretchr/testify/require"
)

var tables = []string{
	user_model.User{}.TableName(),
	tenant_model.Tenant{}.TableName(),
	category_model.Category{}.TableName(),
	product_model.Product{}.TableName(),
	order_model.Order{}.TableName(),
	subscription_model.SubscriptionDelivery{}.TableName(),
	subscription_model.SubscriptionConfig{}.TableName(),
	subscription_model.Subscription{}.TableName(),
}

func CleanDatabase(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	if testDB.Conn == nil {
		err := testDB.Connect(ctx)
		require.NoError(t, err)
	}

	for _, table := range tables {
		err := testDB.Conn.Exec("TRUNCATE " + table + " CASCADE").Error
		require.NoError(t, err)
	}
}
