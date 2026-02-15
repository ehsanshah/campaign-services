package port

import (
	"context"
)

// IContentClient پورت خروجی برای ارتباط با میکروسرویس کانتنت

type IContentClient interface {

	// فریز کردن محتوا برای شروع کمپین

	CreateSnapshot(ctx context.Context, originalID string, campaignID string) (string, error)
}
