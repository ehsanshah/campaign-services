package services

import (
	"context"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/ehsanshah/campaign-services/src/internal/core/port"
)

type CampaignAdService struct {
	repo port.CampaignAdRepository
}

func NewCampaignAdService(repo port.CampaignAdRepository) *CampaignAdService {
	return &CampaignAdService{repo: repo}
}

func (s *CampaignAdService) Create(ctx context.Context, c *domain.CampaignAd) error {
	if err := c.Validate(); err != nil {
		return err
	}
	return s.repo.Save(ctx, c)
}

func (s *CampaignAdService) Get(ctx context.Context, id string) (*domain.CampaignAd, error) {
	return s.repo.GetByID(ctx, id)
}
