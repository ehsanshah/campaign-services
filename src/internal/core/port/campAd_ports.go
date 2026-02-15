package port

import (
	"context"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
)

type CampaignAdRepository interface {
	Save(ctx context.Context, c *domain.CampaignAd) error
	GetByID(ctx context.Context, id string) (*domain.CampaignAd, error)
}

type CampaignAdService interface {
	Create(ctx context.Context, c *domain.CampaignAd) error
	Get(ctx context.Context, id string) (*domain.CampaignAd, error)
}
