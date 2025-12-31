package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// انواع اتوماسیون
type AutomationType string

const (
	AUTO_WELCOME     AutomationType = "welcome"
	AUTO_ABANDONED   AutomationType = "abandoned_cart"
	AUTO_BIRTHDAY    AutomationType = "birthday"
	AUTO_ANNIVERSARY AutomationType = "anniversary"
	AUTO_DRIP        AutomationType = "drip_campaign"
	AUTO_FOLLOWUP    AutomationType = "follow_up"
	AUTO_REACTIVATE  AutomationType = "reactivation"
	AUTO_CUSTOM      AutomationType = "custom"
)

// اتوماسیون ایمیل
type EmailAutomation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        AutomationType     `bson:"type" json:"type"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`

	// رویدادهای شروع‌کننده
	TriggerEvents []TriggerEvent `bson:"trigger_events" json:"trigger_events"`

	// آیا اتوماسیون یکبار برای هر مشترک اجرا شود یا چندبار
	OneTimeOnly bool `bson:"one_time_only" json:"one_time_only"`

	// مراحل اتوماسیون
	Steps []AutomationStep `bson:"steps" json:"automation_steps"`

	// تنظیمات اتوماسیون
	Settings AutomationSettings `bson:"settings" json:"settings"`

	// آمار
	Statistics AutomationStats `bson:"statistics" json:"statistics"`
}

// رویداد شروع‌کننده اتوماسیون
type TriggerEvent struct {
	EventType       string                 `bson:"event_type" json:"event_type"`
	EventConditions map[string]interface{} `bson:"event_conditions" json:"event_conditions"`
	Delay           time.Duration          `bson:"delay" json:"delay"`
}

// مرحله اتوماسیون
type AutomationStep struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StepNumber int                `bson:"step_number" json:"step_number"`
	StepType   string             `bson:"step_type" json:"step_type"` // "email", "delay", "condition", "split", "exit"
	Name       string             `bson:"name" json:"name"`

	// برای مرحله ایمیل
	MessageID primitive.ObjectID `bson:"message_id" json:"message_id"`
	Subject   string             `bson:"subject" json:"subject"`
	FromName  string             `bson:"from_name" json:"from_name"`
	FromEmail string             `bson:"from_email" json:"from_email"`

	// برای مرحله تأخیر
	DelayTime int    `bson:"delay_time" json:"delay_time"` // مقدار به روز
	DelayUnit string `bson:"delay_unit" json:"delay_unit"` // "minutes", "hours", "days", "weeks"

	// برای مرحله شرط
	Conditions     []Condition        `bson:"conditions" json:"conditions"`
	PositiveStepID primitive.ObjectID `bson:"positive_step_id" json:"positive_step_id"`
	NegativeStepID primitive.ObjectID `bson:"negative_step_id" json:"negative_step_id"`

	// برای مرحله تقسیم
	SplitRatio   []float64            `bson:"split_ratio" json:"split_ratio"`
	SplitStepIDs []primitive.ObjectID `bson:"split_step_ids" json:"split_step_ids"`

	// آمار مرحله
	Statistics StepStats `bson:"statistics" json:"statistics"`
}

// شرط
type Condition struct {
	Field    string      `bson:"field" json:"field"`
	Operator string      `bson:"operator" json:"operator"` // "equals", "not_equals", "contains", "not_contains", "greater_than", "less_than", etc.
	Value    interface{} `bson:"value" json:"value"`
}

// تنظیمات اتوماسیون
type AutomationSettings struct {
	TrackingSettings map[string]bool   `bson:"tracking_settings" json:"tracking_settings"`
	SendingHours     []int             `bson:"sending_hours" json:"sending_hours"`
	SendingDays      []int             `bson:"sending_days" json:"sending_days"`
	Timezone         string            `bson:"timezone" json:"timezone"`
	UTMParameters    map[string]string `bson:"utm_parameters" json:"utm_parameters"`
}

// آمار اتوماسیون
type AutomationStats struct {
	EntriesCount     int       `bson:"entries_count" json:"entries_count"`
	CompletionsCount int       `bson:"completions_count" json:"completions_count"`
	ActiveNow        int       `bson:"active_now" json:"active_now"`
	ExitedCount      int       `bson:"exited_count" json:"exited_count"`
	RevenueGenerated float64   `bson:"revenue_generated" json:"revenue_generated"`
	LastUpdatedAt    time.Time `bson:"last_updated_at" json:"last_updated_at"`
}

// آمار مرحله اتوماسیون
type StepStats struct {
	EntriesCount     int       `bson:"entries_count" json:"entries_count"`
	CompletionsCount int       `bson:"completions_count" json:"completions_count"`
	DropOffs         int       `bson:"drop_offs" json:"drop_offs"`
	EmailsSent       int       `bson:"emails_sent" json:"emails_sent"`
	EmailsOpened     int       `bson:"emails_opened" json:"emails_opened"`
	EmailsClicked    int       `bson:"emails_clicked" json:"emails_clicked"`
	ConversionRate   float64   `bson:"conversion_rate" json:"conversion_rate"`
	LastUpdatedAt    time.Time `bson:"last_updated_at" json:"last_updated_at"`
}

// جریان اتوماسیون
type AutomationJourney struct {
	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	AutomationID      primitive.ObjectID     `bson:"automation_id" json:"automation_id"`
	SubscriberID      primitive.ObjectID     `bson:"subscriber_id" json:"subscriber_id"`
	SubscriberEmail   string                 `bson:"subscriber_email" json:"subscriber_email"`
	CurrentStepID     primitive.ObjectID     `bson:"current_step_id" json:"current_step_id"`
	CurrentStepNumber int                    `bson:"current_step_number" json:"current_step_number"`
	Status            string                 `bson:"status" json:"status"` // "active", "completed", "paused", "exited"
	EnteredAt         time.Time              `bson:"entered_at" json:"entered_at"`
	LastUpdatedAt     time.Time              `bson:"last_updated_at" json:"last_updated_at"`
	CompletedAt       time.Time              `bson:"completed_at" json:"completed_at"`
	NextStepAt        time.Time              `bson:"next_step_at" json:"next_step_at"`
	StepHistory       []JourneyStep          `bson:"step_history" json:"step_history"`
	Variables         map[string]interface{} `bson:"variables" json:"variables"`
}

// مرحله جری
// مرحله طی شده در جریان اتوماسیون
type JourneyStep struct {
	StepID      primitive.ObjectID `bson:"step_id" json:"step_id"`
	StepNumber  int                `bson:"step_number" json:"step_number"`
	StepType    string             `bson:"step_type" json:"step_type"` // "email", "delay", "condition", "split", "exit"
	EnteredAt   time.Time          `bson:"entered_at" json:"entered_at"`
	CompletedAt time.Time          `bson:"completed_at" json:"completed_at"`
	Status      string             `bson:"status" json:"status"` // "completed", "skipped", "failed"

	// برای مرحله‌های ایمیل
	MessageID      primitive.ObjectID `bson:"message_id,omitempty" json:"message_id,omitempty"`
	DeliveryID     primitive.ObjectID `bson:"delivery_id,omitempty" json:"delivery_id,omitempty"`
	EmailStatus    string             `bson:"email_status,omitempty" json:"email_status,omitempty"` // "sent", "delivered", "opened", "clicked", "bounced"
	EmailOpenedAt  time.Time          `bson:"email_opened_at,omitempty" json:"email_opened_at,omitempty"`
	EmailClickedAt time.Time          `bson:"email_clicked_at,omitempty" json:"email_clicked_at,omitempty"`

	// برای مرحله‌های شرطی
	ConditionResult bool                   `bson:"condition_result,omitempty" json:"condition_result,omitempty"`
	ConditionData   map[string]interface{} `bson:"condition_data,omitempty" json:"condition_data,omitempty"`

	// برای مرحله‌های تقسیم (split)
	SplitGroup int `bson:"split_group,omitempty" json:"split_group,omitempty"`

	// برای مرحله‌های تأخیر
	DelayDuration  time.Duration `bson:"delay_duration,omitempty" json:"delay_duration,omitempty"`
	ScheduledEndAt time.Time     `bson:"scheduled_end_at,omitempty" json:"scheduled_end_at,omitempty"`

	// برای ذخیره اطلاعات متغیر در هر مرحله
	Variables map[string]interface{} `bson:"variables,omitempty" json:"variables,omitempty"`

	// اطلاعات خطا در صورت وجود
	ErrorMessage string `bson:"error_message,omitempty" json:"error_message,omitempty"`
}
