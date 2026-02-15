package domain

import (
	"errors"
	"time"
)

// وضعیت‌های کمپین (بر اساس لاجیک شما)
const (
	StatusDraft      = "draft"
	StatusScheduled  = "scheduled"
	StatusProcessing = "processing"
	StatusSent       = "sent"
	StatusCancelled  = "cancelled"
	StatusFailed     = "failed"
)

// Campaign: مدل اصلی دقیقاً منطبق با message Campaign در پروتو
type Campaign struct {
	ID        string `json:"id" bson:"_id"`
	AccountID string `json:"account_id" bson:"account_id"`
	Name      string `json:"name" bson:"name"`
	Status    string `json:"status" bson:"status"`

	TypeForHumans string `json:"type_for_humans" bson:"type_for_humans"`

	// اشیاء تو در تو (Nested Objects)
	Recipients CampaignRecipient `json:"recipients" bson:"recipients"`
	Options    CampaignOptions   `json:"options" bson:"options"`
	Stats      CampaignStats     `json:"stats" bson:"stats"`

	// فیلترها (آرایه‌ای از شرط‌ها)
	Filters []FilterCondition `json:"filters" bson:"filters"`

	// زمان‌بندی‌ها
	CreatedAt        time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" bson:"updated_at"`
	ScheduledFor     *time.Time `json:"scheduled_for" bson:"scheduled_for"` // Pointer for nullable
	QueuedAt         *time.Time `json:"queued_at" bson:"queued_at"`
	StartedAt        *time.Time `json:"started_at" bson:"started_at"`
	FinishedAt       *time.Time `json:"finished_at" bson:"finished_at"`
	StoppedAt        *time.Time `json:"stopped_at" bson:"stopped_at"`
	WinnerSelectedAt *time.Time `json:"winner_selected_at" bson:"winner_selected_at"`

	// فلگ‌های کنترلی
	IsStopped          bool `json:"is_stopped" bson:"is_stopped"`
	IsCurrentlySending bool `json:"is_currently_sending_out" bson:"is_currently_sending_out"`
	CanBeScheduled     bool `json:"can_be_scheduled" bson:"can_be_scheduled"`
	HasWinner          bool `json:"has_winner" bson:"has_winner"`

	// اطلاعات A/B Test و خروجی انسانی
	WinnerVersionForHuman      string `json:"winner_version_for_human" bson:"winner_version_for_human"`
	WinnerSendingTimeForHumans string `json:"winner_sending_time_for_humans" bson:"winner_sending_time_for_humans"`

	// محتوا
	EmailIDs       []string `json:"email_ids" bson:"email_ids"`
	DefaultEmailID string   `json:"default_email_id" bson:"default_email_id"`

	Warnings []string `json:"warnings" bson:"warnings"`

	UsedInAutomations bool           `json:"used_in_automations" bson:"used_in_automations"`
	ExtraFields       map[string]any `json:"extra_fields" bson:"extra_fields"` // JSONB
}

// ---------------------------------------------
// مدل‌های زیرمجموعه (Sub-structs)
// ---------------------------------------------

type CampaignRecipient struct {
	ListIDs      []string `json:"list_ids" bson:"list_ids"`
	SegmentIDs   []string `json:"segment_ids" bson:"segment_ids"`
	ListNames    []string `json:"list_names" bson:"list_names"`
	SegmentNames []string `json:"segment_names" bson:"segment_names"`
}

type CampaignOptions struct {
	DeliveryOptimization string `json:"delivery_optimization" bson:"delivery_optimization"`
	TrackOpens           bool   `json:"track_opens" bson:"track_opens"`
	TrackClicks          bool   `json:"track_clicks" bson:"track_clicks"`
	UseGoogleAnalytics   bool   `json:"use_google_analytics" bson:"use_google_analytics"`
	EcommerceTracking    bool   `json:"ecommerce_tracking" bson:"ecommerce_tracking"`
	TriggerFrequency     int32  `json:"trigger_frequency" bson:"trigger_frequency"`
	TriggerCount         int32  `json:"trigger_count" bson:"trigger_count"`
	UsesSurvey           bool   `json:"uses_survey" bson:"uses_survey"`
}

type FilterCondition struct {
	Operator string `json:"operator" bson:"operator"`
	Args     []any  `json:"args" bson:"args"` // نگاشت google.protobuf.Value
}

// StatsRate مدل کمکی برای نرخ‌ها
type StatsRate struct {
	Value float64 `json:"value" bson:"value"` // نگاشت به float در پروتو
	Text  string  `json:"text" bson:"text"`   // نگاشت به string در پروتو
}

type CampaignStats struct {
	Sent             int64     `json:"sent" bson:"sent"`
	OpensCount       int64     `json:"opens_count" bson:"opens_count"`
	UniqueOpensCount int64     `json:"unique_opens_count" bson:"unique_opens_count"`
	OpenRate         StatsRate `json:"open_rate" bson:"open_rate"`

	ClicksCount       int64     `json:"clicks_count" bson:"clicks_count"`
	UniqueClicksCount int64     `json:"unique_clicks_count" bson:"unique_clicks_count"`
	ClickRate         StatsRate `json:"click_rate" bson:"click_rate"`

	UnsubscribesCount int64     `json:"unsubscribes_count" bson:"unsubscribes_count"`
	UnsubscribeRate   StatsRate `json:"unsubscribe_rate" bson:"unsubscribe_rate"`

	SpamCount int64     `json:"spam_count" bson:"spam_count"`
	SpamRate  StatsRate `json:"spam_rate" bson:"spam_rate"`

	HardBouncesCount int64     `json:"hard_bounces_count" bson:"hard_bounces_count"`
	HardBounceRate   StatsRate `json:"hard_bounce_rate" bson:"hard_bounce_rate"`

	SoftBouncesCount int64     `json:"soft_bounces_count" bson:"soft_bounces_count"`
	SoftBounceRate   StatsRate `json:"soft_bounce_rate" bson:"soft_bounce_rate"`

	DeliveryRate float64 `json:"delivery_rate" bson:"delivery_rate"`
}

// ---------------------------------------------
// منطق‌های دامین (Domain Logic)
// ---------------------------------------------

// ValidateTransition بررسی می‌کند آیا تغییر وضعیت مجاز است؟
func (c *Campaign) ValidateTransition(newStatus string) error {
	// قوانین ساده شده برای شروع:
	if c.IsStopped && newStatus == StatusProcessing {
		return errors.New("cannot restart a stopped campaign directly")
	}

	// نمی‌توان کمپین تکمیل شده را دوباره درفت کرد
	if c.Status == StatusSent && newStatus == StatusDraft {
		return errors.New("cannot revert sent campaign to draft")
	}

	return nil
}
