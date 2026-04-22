package usecases

import (
	"context"

	"github.com/Rabi-IT/rabi-food-core/domain/tenant"
	g "github.com/Rabi-IT/rabi-food-core/features/tenant/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
)

func (c *TenantCase) Create(ctx context.Context, name, ownerUserID string) (string, error) {
	input := g.CreateInput{
		Name:     name,
		Slug:     slugify(name),
		Language: "pt-BR",
		InitialMember: g.InitialMemberInput{
			UserID: ownerUserID,
			Role:   tenant.Owner,
		},
	}

	id, err := c.gateway.Create(ctx, input)
	logger.GetWideEvent(ctx).SetTenantID(id)

	return id, err
}
