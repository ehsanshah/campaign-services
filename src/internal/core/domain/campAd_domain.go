package domain

import (
	"encoding/json"
	"errors"
)

// CampaignAdType دقیقاً منطبق با Enum در فایل campaign_ad.proto
type CampaignAdType int32

const (
	CampaignTypeUnspecified CampaignAdType = 0
	CampaignTypeAd          CampaignAdType = 1 // تبلیغاتی
	CampaignTypeMTA         CampaignAdType = 2 // ایمیلی
)

type CampaignAd struct {
	ID             string
	OrganizationID string
	Name           string
	Status         string
	Type           CampaignAdType // فیلد تایپ با نام جدید
	ScheduledAt    int64

	// جزئیات متغیر (AdDetails یا EmailDetails)
	Details map[string]interface{}
}

// متد کمکی برای تبدیل جزئیات به JSON

func (c *CampaignAd) DetailsToJSON() ([]byte, error) {
	return json.Marshal(c.Details)
}

func (c *CampaignAd) Validate() error {
	if c.Name == "" {
		return errors.New("campaign name is required")
	}
	if c.Type == CampaignTypeUnspecified {
		return errors.New("campaign type must be specified")
	}
	return nil
}
