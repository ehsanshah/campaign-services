package domain

import (
	"errors"
	"time"
)

// وضعیت‌های کمپین (بر اساس لاجیک شما)
const (
	StatusDraft      = "DRAFT"
	StatusScheduled  = "SCHEDULED"
	StatusProcessing = "PROCESSING"
	StatusSent       = "SENT"
	StatusCancelled  = "CANCELLED"
	StatusFailed     = "FAILED"
)

// Campaign: مدل اصلی دقیقاً منطبق با message Campaign در پروتو

type Campaign struct {
	ID        string `json:"id" db:"id"`
	AccountID string `json:"account_id" db:"account_id"`
	Name      string `json:"name" db:"name"`
	Status    string `json:"status" db:"status"`

	TypeForHumans string `json:"type_for_humans" db:"type_for_humans"`

	// اشیاء تو در تو (Nested Objects) که در دیتابیس Postgres به صورت JSONB ذخیره می‌شوند
	Recipients CampaignRecipient `json:"recipients" db:"-"` // db:- یعنی مستقیم مپ نمیشه، در Repository هندل میشه
	Options    CampaignOptions   `json:"options" db:"-"`
	Stats      CampaignStats     `json:"stats" db:"-"`

	// فیلترها (آرایه‌ای از شرط‌ها)
	Filters []FilterCondition `json:"filters" db:"-"`

	// زمان‌بندی‌ها
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	ScheduledFor     *time.Time `json:"scheduled_for" db:"scheduled_for"` // Pointer for nullable
	QueuedAt         *time.Time `json:"queued_at" db:"queued_at"`
	StartedAt        *time.Time `json:"started_at" db:"started_at"`
	FinishedAt       *time.Time `json:"finished_at" db:"finished_at"`
	StoppedAt        *time.Time `json:"stopped_at" db:"stopped_at"`
	WinnerSelectedAt *time.Time `json:"winner_selected_at" db:"winner_selected_at"`

	// فلگ‌های کنترلی
	IsStopped          bool `json:"is_stopped" db:"is_stopped"`
	IsCurrentlySending bool `json:"is_currently_sending_out" db:"is_currently_sending_out"`
	CanBeScheduled     bool `json:"can_be_scheduled" db:"can_be_scheduled"`
	HasWinner          bool `json:"has_winner" db:"has_winner"`

	// اطلاعات A/B Test و خروجی انسانی
	WinnerVersionForHuman      string `json:"winner_version_for_human" db:"winner_version_for_human"`
	WinnerSendingTimeForHumans string `json:"winner_sending_time_for_humans" db:"winner_sending_time_for_humans"`

	// محتوا
	EmailIDs       []string `json:"email_ids" db:"email_ids"` // Postgres Array
	DefaultEmailID string   `json:"default_email_id" db:"default_email_id"`

	Warnings []string `json:"warnings" db:"warnings"`

	UsedInAutomations bool           `json:"used_in_automations" db:"used_in_automations"`
	ExtraFields       map[string]any `json:"extra_fields" db:"-"` // JSONB
}

// ---------------------------------------------
// مدل‌های زیرمجموعه (Sub-structs)
// ---------------------------------------------

type CampaignRecipient struct {
	ListIDs      []string `json:"list_ids"`
	SegmentIDs   []string `json:"segment_ids"`
	ListNames    []string `json:"list_names"`
	SegmentNames []string `json:"segment_names"`
}

type CampaignOptions struct {
	DeliveryOptimization string `json:"delivery_optimization"`
	TrackOpens           bool   `json:"track_opens"`
	TrackClicks          bool   `json:"track_clicks"`
	UseGoogleAnalytics   bool   `json:"use_google_analytics"`
	EcommerceTracking    bool   `json:"ecommerce_tracking"`
	TriggerFrequency     int32  `json:"trigger_frequency"`
	TriggerCount         int32  `json:"trigger_count"`
	UsesSurvey           bool   `json:"uses_survey"`
}

type FilterCondition struct {
	Operator string `json:"operator"`
	Args     []any  `json:"args"` // نگاشت google.protobuf.Value
}

// StatsRate مدل کمکی برای نرخ‌ها (طبق فایل پروتو: double float, string string)
type StatsRate struct {
	Float  float64 `json:"float"`  // نام فیلد در پروتو float بود
	String string  `json:"string"` // نام فیلد در پروتو string بود
}

type CampaignStats struct {
	Sent             int64     `json:"sent"`
	OpensCount       int64     `json:"opens_count"`
	UniqueOpensCount int64     `json:"unique_opens_count"`
	OpenRate         StatsRate `json:"open_rate"`

	ClicksCount       int64     `json:"clicks_count"`
	UniqueClicksCount int64     `json:"unique_clicks_count"`
	ClickRate         StatsRate `json:"click_rate"`

	UnsubscribesCount int64     `json:"unsubscribes_count"`
	UnsubscribeRate   StatsRate `json:"unsubscribe_rate"`

	SpamCount int64     `json:"spam_count"`
	SpamRate  StatsRate `json:"spam_rate"`

	HardBouncesCount int64     `json:"hard_bounces_count"`
	HardBounceRate   StatsRate `json:"hard_bounce_rate"`

	SoftBouncesCount int64     `json:"soft_bounces_count"`
	SoftBounceRate   StatsRate `json:"soft_bounce_rate"`

	DeliveryRate float64 `json:"delivery_rate"`
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
