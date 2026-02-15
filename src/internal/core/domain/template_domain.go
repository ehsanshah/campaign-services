package domain

import "time"

// Template موجودیت اصلی قالب
type Template struct {
	ID               string    `json:"id"`
	AccountID        string    `json:"account_id"` // پی‌نوشت ۱۰
	Name             string    `json:"name"`
	CurrentVersionID string    `json:"current_version_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TemplateVersion جزئیات هر نسخه از قالب
type TemplateVersion struct {
	ID           string    `json:"id"`
	TemplateID   string    `json:"template_id"`
	VersionLabel string    `json:"version_label"`
	Subject      string    `json:"subject"`
	HTMLContent  string    `json:"html_content"`
	PlainText    string    `json:"plain_text"`
	Language     string    `json:"language"`
	Tags         []string  `json:"tags"`
	Metadata     string    `json:"metadata"` // JSON string
	CreatedAt    time.Time `json:"created_at"`
}
