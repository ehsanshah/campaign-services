package port

import (
	"context"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
)

type ITemplateRepository interface {
	SaveTemplate(ctx context.Context, t *domain.Template) error
	SaveVersion(ctx context.Context, v *domain.TemplateVersion) error
	GetTemplate(ctx context.Context, accountID, template_id string) (*domain.Template, *domain.TemplateVersion, error)
	ListTemplates(ctx context.Context, accountID string, limit, offset int32) ([]*domain.Template, int32, error)
	DeleteTemplate(ctx context.Context, accountID, template_id string) error
}

type ITemplateServices interface {
	CreateTemplate(ctx context.Context, t *domain.Template, v *domain.TemplateVersion) (*domain.Template, error)
	UpdateTemplate(ctx context.Context, accountID string, templateID string, v *domain.TemplateVersion) (*domain.Template, error)
	CopyTemplate(ctx context.Context, accountID, sourceID, newName string) (*domain.Template, error)
	ImportTemplateFromUrl(ctx context.Context, accountID, name, url string) (*domain.Template, error)
	TestTemplate(ctx context.Context, accountID, templateID, versionID, testEmail string) error // ✅ اضافه شد
}
