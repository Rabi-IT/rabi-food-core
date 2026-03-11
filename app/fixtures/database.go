package fixtures

import (
	"context"
	"rabi-food-core/libs/database/gorm_adapter/models"
	"testing"

	"github.com/stretchr/testify/require"
)

var tables = []string{
	models.User{}.TableName(),
	models.Tenant{}.TableName(),
	models.Category{}.TableName(),
	models.Product{}.TableName(),
	models.Order{}.TableName(),
	models.SubscriptionDelivery{}.TableName(),
	models.SubscriptionConfig{}.TableName(),
	models.Subscription{}.TableName(),
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
