package services

import (
	"context"
	"errors"
	"time"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/ehsanshah/campaign-services/src/internal/core/port"
	"github.com/google/uuid"
)

type CampaignService struct {
	repo port.ICampaignRepository
}

func NewCampaignServiceMta(repo port.ICampaignRepository) port.ICampaignService {
	return &CampaignService{
		repo: repo,
	}
}

// CreateCampaign: ایجاد کمپین اولیه (همیشه با وضعیت Draft)
func (s *CampaignService) CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	// ۱. تنظیم مقادیر پیش‌فرض
	campaign.ID = uuid.New().String()
	campaign.Status = domain.StatusDraft
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()

	// پیش‌فرض‌های منطقی
	campaign.CanBeScheduled = true
	campaign.IsStopped = false

	// ۲. ذخیره در دیتابیس
	err := s.repo.Create(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	// ۱. بررسی وجود کمپین
	existing, err := s.repo.GetByID(ctx, campaign.ID, campaign.AccountID)
	if err != nil {
		return nil, err
	}

	// ۲. قانون: کمپینی که ارسال شده یا در حال ارسال است، نباید کامل ادیت شود
	// (مگر اینکه لاجیک خاصی داشته باشید، ولی معمولاً قفل می‌شود)
	if existing.Status != domain.StatusDraft && existing.Status != domain.StatusScheduled {
		return nil, errors.New("cannot update campaign that is already processing or sent")
	}

	// ۳. آپدیت فیلدها
	campaign.UpdatedAt = time.Now()
	// نکته: اینجا باید فیلدهای خالی را مدیریت کنیم تا نال نشوند (Merge Logic)
	// اما برای سادگی فرض می‌کنیم کلاینت آبجکت کامل را فرستاده است.

	err = s.repo.Update(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *CampaignService) GetCampaign(ctx context.Context, id string, accountID string) (*domain.Campaign, error) {
	return s.repo.GetByID(ctx, id, accountID)
}

func (s *CampaignService) ListCampaigns(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Campaign, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(ctx, accountID, limit, offset)
}

// ScheduleCampaign: حساس‌ترین متد بیزنس لاجیک
func (s *CampaignService) ScheduleCampaign(ctx context.Context, id string, accountID string, sendAt time.Time) (*domain.Campaign, error) {
	// ۱. واکشی کمپین
	campaign, err := s.repo.GetByID(ctx, id, accountID)
	if err != nil {
		return nil, err
	}

	// ۲. اعتبارسنجی وضعیت (فقط Draft می‌تواند زمان‌بندی شود)
	if campaign.Status != domain.StatusDraft {
		return nil, errors.New("only draft campaigns can be scheduled")
	}

	// ۳. اعتبارسنجی محتوا (آیا گیرنده دارد؟ آیا محتوا دارد؟)
	if len(campaign.Recipients.ListIDs) == 0 && len(campaign.Recipients.SegmentIDs) == 0 {
		return nil, errors.New("campaign must have at least one recipient list or segment")
	}
	if len(campaign.EmailIDs) == 0 {
		return nil, errors.New("campaign must have content linked to it")
	}

	// ۴. اعمال تغییرات
	campaign.ScheduledFor = &sendAt
	campaign.Status = domain.StatusScheduled
	campaign.UpdatedAt = time.Now()

	// ۵. ذخیره تغییرات
	err = s.repo.Update(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *CampaignService) CancelCampaign(ctx context.Context, id string, accountID string) (*domain.Campaign, error) {
	campaign, err := s.repo.GetByID(ctx, id, accountID)
	if err != nil {
		return nil, err
	}

	// فقط کمپین‌های Scheduled یا Processing را می‌توان کنسل کرد
	if campaign.Status == domain.StatusSent || campaign.Status == domain.StatusFailed {
		return nil, errors.New("cannot cancel a finished campaign")
	}

	campaign.Status = domain.StatusCancelled
	campaign.IsStopped = true
	now := time.Now()
	campaign.StoppedAt = &now
	campaign.UpdatedAt = now

	err = s.repo.Update(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, id string, accountID string) error {
	// معمولاً Soft Delete پیشنهاد می‌شود، اما طبق متد Repository فعلاً Hard Delete می‌کنیم
	return s.repo.Delete(ctx, id, accountID)
}
