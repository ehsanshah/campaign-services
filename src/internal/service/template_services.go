package services

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/ehsanshah/campaign-services/src/internal/core/port"
)

// ۱. اصلاح فیلد repo: باید ریپازیتوری باشد نه سرویس!
type templateServices struct {
	repo port.ITemplateRepository
}

func NewTemplateServices(repo port.ITemplateRepository) port.ITemplateServices {
	return &templateServices{
		repo: repo,
	}
}

// ---------------------------------------------------------
// پیاده‌سازی تمام متدهای ITemplateServices (قرارداد)
// ---------------------------------------------------------

// 1️⃣ متد CreateTemplate

func (s *templateServices) CreateTemplate(ctx context.Context, t *domain.Template, v *domain.TemplateVersion) (*domain.Template, error) {
	if err := s.repo.SaveTemplate(ctx, t); err != nil {
		return nil, err
	}
	v.TemplateID = t.ID
	if err := s.repo.SaveVersion(ctx, v); err != nil {
		return nil, err
	}
	t.CurrentVersionID = v.ID
	_ = s.repo.SaveTemplate(ctx, t)
	return t, nil
}

// 2️⃣ متد UpdateTemplate (باید پیاده‌سازی می‌شد تا کامپایل شود)

func (s *templateServices) UpdateTemplate(ctx context.Context, accountID string, templateID string, v *domain.TemplateVersion) (*domain.Template, error) {
	// ابتدا بررسی وجود قالب
	t, _, err := s.repo.GetTemplate(ctx, accountID, templateID)
	if err != nil {
		return nil, err
	}

	// ذخیره نسخه جدید
	v.TemplateID = templateID
	if err := s.repo.SaveVersion(ctx, v); err != nil {
		return nil, err
	}

	// بروزرسانی نسخه فعلی قالب
	t.CurrentVersionID = v.ID
	if err := s.repo.SaveTemplate(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

// 3️⃣ متد CopyTemplate

func (s *templateServices) CopyTemplate(ctx context.Context, accountID, sourceID, newName string) (*domain.Template, error) {
	_, oldV, err := s.repo.GetTemplate(ctx, accountID, sourceID)
	if err != nil {
		return nil, fmt.Errorf("source template not found: %w", err)
	}

	newT := &domain.Template{
		AccountID: accountID,
		Name:      newName,
	}
	if err := s.repo.SaveTemplate(ctx, newT); err != nil {
		return nil, err
	}

	newV := *oldV
	newV.ID = ""
	newV.TemplateID = newT.ID
	newV.VersionLabel = "Copy of " + oldV.VersionLabel

	if err := s.repo.SaveVersion(ctx, &newV); err != nil {
		return nil, err
	}

	newT.CurrentVersionID = newV.ID
	_ = s.repo.SaveTemplate(ctx, newT)

	return newT, nil
}

// 4️⃣ متد ImportTemplateFromUrl

func (s *templateServices) ImportTemplateFromUrl(ctx context.Context, accountID, name, url string) (*domain.Template, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch content from URL")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	t := &domain.Template{AccountID: accountID, Name: name}
	if err := s.repo.SaveTemplate(ctx, t); err != nil {
		return nil, err
	}

	v := &domain.TemplateVersion{
		TemplateID:   t.ID,
		HTMLContent:  string(body),
		Subject:      "Imported Template",
		VersionLabel: "v1 (Imported)",
	}
	if err := s.repo.SaveVersion(ctx, v); err != nil {
		return nil, err
	}

	t.CurrentVersionID = v.ID
	_ = s.repo.SaveTemplate(ctx, t)

	return t, nil
}

// 5️⃣ متد TestTemplate (برای هماهنگی با پروتو و اینترفیس)

func (s *templateServices) TestTemplate(ctx context.Context, accountID, templateID, versionID, testEmail string) error {
	// منطق ارسال تست (MTA Call) در اینجا قرار می‌گیرد
	fmt.Printf("Sending test email to %s for template %s\n", testEmail, templateID)
	return nil
}
