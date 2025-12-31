package domain

// Ab - مدل آمار دامنه‌ها

// ABTestVariation - مدل نسخه‌های آزمایش A/B
type ABTestVariation struct {
	VariationID         string `bson:"variation_id" json:"variationId"`
	Name                string `bson:"name" json:"name"`
	Subject             string `bson:"subject" json:"subject"`
	PreviewText         string `bson:"preview_text" json:"previewText"`
	HTMLContent         string `bson:"html_content" json:"htmlContent"`
	DistributionPercent int    `bson:"distribution_percent" json:"distributionPercent"`
	Sent                int    `bson:"sent" json:"sent"`
	Delivered           int    `bson:"delivered" json:"delivered"`
	Opened              int    `bson:"opened" json:"opened"`
	Clicked             int    `bson:"clicked" json:"clicked"`
	IsWinner            bool   `bson:"is_winner" json:"isWinner"`
}
