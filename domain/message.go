package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// مدل پیام بهینه‌شده
type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID   primitive.ObjectID `bson:"client_id" json:"client_id"`
	Name       string             `bson:"name" json:"name"`
	Subject    string             `bson:"subject" json:"subject"`
	PreHeader  string             `bson:"pre_header" json:"pre_header"`
	BodyHTML   string             `bson:"body_html" json:"body_html"`
	BodyText   string             `bson:"body_text" json:"body_text"`
	TemplateID string             `bson:"template_id" json:"template_id"`
	IsTemplate bool               `bson:"is_template" json:"is_template"`
	Language   string             `bson:"language" json:"language"`
	FolderID   primitive.ObjectID `bson:"folder_id" json:"folder_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy  primitive.ObjectID `bson:"updated_by" json:"updated_by"`

	// تنظیمات طراحی
	BuiltWithEditor bool   `bson:"built_with_editor" json:"built_with_editor"`
	LinkAlignment   string `bson:"link_alignment" json:"link_alignment"`
	TopLink         bool   `bson:"top_link" json:"top_link"`

	// متغیرهای شخصی‌سازی
	Variables map[string]string `bson:"variables" json:"variables"`

	// اطلاعات متاداده
	Metadata map[string]string `bson:"metadata" json:"metadata"`
}

// پوشه پیام‌ها
type MessageFolder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ParentID    primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
