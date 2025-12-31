package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// وضعیت‌های پردازش صف
type QueueStatus string

const (
	QUEUE_PENDING    QueueStatus = "pending"
	QUEUE_PROCESSING QueueStatus = "processing"
	QUEUE_COMPLETED  QueueStatus = "completed"
	QUEUE_FAILED     QueueStatus = "failed"
	QUEUE_CANCELLED  QueueStatus = "cancelled"
)

// اولویت صف
type QueuePriority int

const (
	PRIORITY_LOW    QueuePriority = 1
	PRIORITY_NORMAL QueuePriority = 5
	PRIORITY_HIGH   QueuePriority = 10
	PRIORITY_URGENT QueuePriority = 20
)

// آیتم صف پیام
type QueueItem struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CampaignID     primitive.ObjectID     `bson:"campaign_id" json:"campaign_id"`
	MessageID      primitive.ObjectID     `bson:"message_id" json:"message_id"`
	RecipientEmail string                 `bson:"recipient_email" json:"recipient_email"`
	RecipientData  map[string]interface{} `bson:"recipient_data" json:"recipient_data"`
	Status         QueueStatus            `bson:"status" json:"status"`
	Priority       QueuePriority          `bson:"priority" json:"priority"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	ScheduledAt    time.Time              `bson:"scheduled_at" json:"scheduled_at"`
	ProcessedAt    time.Time              `bson:"processed_at" json:"processed_at"`
	CompletedAt    time.Time              `bson:"completed_at" json:"completed_at"`
	RetryCount     int                    `bson:"retry_count" json:"retry_count"`
	MaxRetries     int                    `bson:"max_retries" json:"max_retries"`
	LastError      string                 `bson:"last_error" json:"last_error"`
	ExternalID     string                 `bson:"external_id" json:"external_id"`

	// اطلاعات RabbitMQ
	QueueName  string `bson:"queue_name" json:"queue_name"`
	Exchange   string `bson:"exchange" json:"exchange"`
	RoutingKey string `bson:"routing_key" json:"routing_key"`
	Server     string `bson:"server" json:"server"`
}

// مدیر پردازش صف
type QueueBatch struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CampaignID     primitive.ObjectID `bson:"campaign_id" json:"campaign_id"`
	Status         QueueStatus        `bson:"status" json:"status"`
	ItemCount      int                `bson:"item_count" json:"item_count"`
	ProcessedCount int                `bson:"processed_count" json:"processed_count"`
	SuccessCount   int                `bson:"success_count" json:"success_count"`
	FailedCount    int                `bson:"failed_count" json:"failed_count"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	StartedAt      time.Time          `bson:"started_at" json:"started_at"`
	CompletedAt    time.Time          `bson:"completed_at" json:"completed_at"`
	WorkerID       string             `bson:"worker_id" json:"worker_id"`
}
