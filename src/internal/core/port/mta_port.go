package port

import "context"

// IMtaService اینترفیس ارتباطی با میکروسرویس ارسال ایمیل
type IMtaService interface {
	SendImmediate(ctx context.Context, to string, subject string, html string) error
}
