package postgres

import (
	"context"
	"database/sql"
	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	_ "github.com/jackc/pgx/v5/pgxpool" // برای تبدیل []string به آرایه Postgres
	"github.com/lib/pq"
)

type templateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) *templateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) SaveTemplate(ctx context.Context, t *domain.Template) error {
	query := `INSERT INTO templates (id, account_id, name, current_version_id, updated_at)
	          VALUES (COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, NULLIF($4, '')::uuid, NOW())
	          ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, current_version_id = EXCLUDED.current_version_id, updated_at = NOW()
	          RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, t.ID, t.AccountID, t.Name, t.CurrentVersionID).Scan(&t.ID, &t.CreatedAt)
}

func (r *templateRepository) SaveVersion(ctx context.Context, v *domain.TemplateVersion) error {
	query := `INSERT INTO template_versions (template_id, version_label, subject, html_content, plain_text, language, tags, metadata)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		v.TemplateID, v.VersionLabel, v.Subject, v.HTMLContent, v.PlainText,
		v.Language, pq.Array(v.Tags), v.Metadata,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *templateRepository) GetTemplate(ctx context.Context, accountID, templateID string) (*domain.Template, *domain.TemplateVersion, error) {
	var t domain.Template
	var v domain.TemplateVersion
	query := `SELECT t.id, t.name, t.current_version_id, tv.subject, tv.html_content, tv.tags, tv.metadata
	          FROM templates t 
	          LEFT JOIN template_versions tv ON t.current_version_id = tv.id
	          WHERE t.account_id = $1 AND t.id = $2`

	err := r.db.QueryRowContext(ctx, query, accountID, templateID).Scan(
		&t.ID, &t.Name, &t.CurrentVersionID, &v.Subject, &v.HTMLContent, pq.Array(&v.Tags), &v.Metadata,
	)
	return &t, &v, err
}
