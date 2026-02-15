package port

import (
	"context"
	"time"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
)

// ICampaignRepository: لایه دیتابیس (Postgres)
type ICampaignRepository interface {
	// CRUD پایه
	Create(ctx context.Context, campaign *domain.Campaign) error
	Update(ctx context.Context, campaign *domain.Campaign) error
	GetByID(ctx context.Context, id string, accountID string) (*domain.Campaign, error)
	Delete(ctx context.Context, id string, accountID string) error

	// لیست کردن با صفحه بندی (طبق ListCampaignsRequest)
	// فیلترهای لیست (مثل جستجو بر اساس نام) را می‌توان در options پاس داد
	List(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Campaign, error)

	// متد اختصاصی برای تغییر وضعیت سریع
	UpdateStatus(ctx context.Context, id string, status string) error
}

// ICampaignService: لایه بیزنس (UseCase)
// این متدها دقیقاً متناظر با RPCهای فایل پروتو هستند
type ICampaignService interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)

	UpdateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)

	GetCampaign(ctx context.Context, id string, accountID string) (*domain.Campaign, error)

	ListCampaigns(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Campaign, error)

	// نگاشت ScheduleCampaignRequest
	ScheduleCampaign(ctx context.Context, id string, accountID string, sendAt time.Time) (*domain.Campaign, error)

	// نگاشت CancelCampaignRequest
	CancelCampaign(ctx context.Context, id string, accountID string) (*domain.Campaign, error)

	// نگاشت DeleteCampaignRequest
	DeleteCampaign(ctx context.Context, id string, accountID string) error
}
