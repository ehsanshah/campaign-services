package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// وضعیت‌های ممکن برای کمپین
type CampaignStatus string

const (
	STATUS_DRAFT     CampaignStatus = "draft"
	STATUS_SCHEDULED CampaignStatus = "scheduled"
	STATUS_SENDING   CampaignStatus = "sending"
	STATUS_PAUSED    CampaignStatus = "paused"
	STATUS_COMPLETED CampaignStatus = "completed"
	STATUS_CANCELLED CampaignStatus = "cancelled"
	STATUS_FAILED    CampaignStatus = "failed"
)

// انواع کمپین
type CampaignType string

const (
	TYPE_REGULAR       CampaignType = "regular"
	TYPE_AUTOMATED     CampaignType = "automated"
	TYPE_AB_TEST       CampaignType = "ab_test"
	TYPE_TRIGGERED     CampaignType = "triggered"
	TYPE_TRANSACTIONAL CampaignType = "transactional"
)

// مدل سگمنت کمپین
type CampaignSegment struct {
	SegmentID         string                 `bson:"segment_id" json:"segment_id"`
	SegmentName       string                 `bson:"segment_name" json:"segment_name"`
	FilterCriteria    map[string]interface{} `bson:"filter_criteria,omitempty" json:"filter_criteria,omitempty"`
	EstimatedCount    int                    `bson:"estimated_count" json:"estimated_count"`
	ListIDs           []string               `bson:"list_ids,omitempty" json:"list_ids,omitempty"`
	IncludeSegmentIDs []string               `bson:"include_segment_ids,omitempty" json:"include_segment_ids,omitempty"`
	ExcludeSegmentIDs []string               `bson:"exclude_segment_ids,omitempty" json:"exclude_segment_ids,omitempty"`
}

// تنظیمات کمپین
type CampaignSettings struct {
	OpenTracking        bool      `bson:"open_tracking" json:"open_tracking"`
	ClickTracking       bool      `bson:"click_tracking" json:"click_tracking"`
	Sandbox             bool      `bson:"sandbox" json:"sandbox"`
	InlineCss           bool      `bson:"inline_css" json:"inline_css"`
	IpPool              string    `bson:"ip_pool" json:"ip_pool"`
	UnsubscribeTracking bool      `bson:"unsubscribe_tracking" json:"unsubscribe_tracking"`
	GoogleAnalytics     bool      `bson:"google_analytics" json:"google_analytics"`
	UTMParameters       UTMParams `bson:"utm_parameters" json:"utm_parameters"`
}

// پارامترهای UTM برای ردیابی
type UTMParams struct {
	Source   string `bson:"source" json:"source"`
	Medium   string `bson:"medium" json:"medium"`
	Campaign string `bson:"campaign" json:"campaign"`
	Term     string `bson:"term" json:"term"`
	Content  string `bson:"content" json:"content"`
}

// مدل اصلی کمپین
type Campaign struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID    primitive.ObjectID `bson:"client_id" json:"client_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        CampaignType       `bson:"type" json:"type"`
	Status      CampaignStatus     `bson:"status" json:"status"`
	Token       string             `bson:"token" json:"token"`

	// زمان‌بندی
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	ScheduledAt time.Time `bson:"scheduled_at" json:"scheduled_at"`
	SentAt      time.Time `bson:"sent_at" json:"sent_at"`
	CompletedAt time.Time `bson:"completed_at" json:"completed_at"`
	Timezone    string    `bson:"timezone" json:"timezone"`

	// محتوا
	MessageID    primitive.ObjectID `bson:"message_id" json:"message_id"`
	Subject      string             `bson:"subject" json:"subject"`
	PreHeader    string             `bson:"pre_header" json:"pre_header"`
	SenderName   string             `bson:"sender_name" json:"sender_name"`
	SenderEmail  string             `bson:"sender_email" json:"sender_email"`
	ReplyToEmail string             `bson:"reply_to_email" json:"reply_to_email"`

	// مخاطبین
	SegmentIDs      []string  `bson:"segment_ids" json:"segment_ids"`
	ListIDs         []string  `bson:"list_ids" json:"list_ids"`
	ExcludedListIDs []string  `bson:"excluded_list_ids" json:"excluded_list_ids"`
	TestEmails      []string  `bson:"test_emails" json:"test_emails"`
	Messages        []Message `bson:"messages" json:"messages"` // پیام‌های مرتبط با این کمپین
	// تنظیمات
	Settings   CampaignSettings  `bson:"settings" json:"settings"`
	Categories []string          `bson:"categories" json:"categories"`
	Tags       []string          `bson:"tags" json:"tags"`
	Segments   []CampaignSegment `bson:"segments" json:"segments"` // سگمنت‌های مخاطبان این کمپین

	// آمار
	//Statistics      CampaignStats      `bson:"statistics" json:"statistics"`
}
