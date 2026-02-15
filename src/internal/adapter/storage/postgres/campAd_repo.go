package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool" // درایور جدید
)

type CampaignRepo struct {
	db *pgxpool.Pool // تایپ تغییر کرد
}

func NewCampaignRepo(db *pgxpool.Pool) *CampaignRepo {
	return &CampaignRepo{db: db}
}

func (r *CampaignRepo) Save(ctx context.Context, c *domain.CampaignAd) error {
	detailsJSON, err := c.DetailsToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}

	query := `
		INSERT INTO campaigns_ad (id, organization_id, name, status, campaign_type, details, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	// در pgx متد ExecContext وجود ندارد، خود Exec کانترکست می‌گیرد
	_, err = r.db.Exec(ctx, query,
		c.ID, c.OrganizationID, c.Name, c.Status, c.Type, detailsJSON, c.ScheduledAt)

	return err
}

func (r *CampaignRepo) GetByID(ctx context.Context, id string) (*domain.CampaignAd, error) {
	query := `SELECT id, organization_id, name, status, campaign_type, details, scheduled_at FROM campaigns_ad WHERE id = $1`

	var c domain.CampaignAd
	var detailsBytes []byte

	// در pgx هم QueryRow کانکست می‌گیرد
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&c.ID, &c.OrganizationID, &c.Name, &c.Status, &c.Type, &detailsBytes, &c.ScheduledAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(detailsBytes, &c.Details); err != nil {
		return nil, err
	}

	return &c, nil
}
