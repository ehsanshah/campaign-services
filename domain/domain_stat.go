package domain

// DomainStat - مدل آمار دامنه‌ها
type DomainStat struct {
	Domain         string  `bson:"domain" json:"domain"`
	Count          int     `bson:"count" json:"count"`
	DeliveryRate   float64 `bson:"delivery_rate" json:"deliveryRate"`
	OpenRate       float64 `bson:"open_rate" json:"openRate"`
	ClickRate      float64 `bson:"click_rate" json:"clickRate"`
	BounceRate     float64 `bson:"bounce_rate" json:"bounceRate"`
	SpamReportRate float64 `bson:"spam_report_rate" json:"spamReportRate"`
}
