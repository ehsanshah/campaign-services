DROP TABLE IF EXISTS campaigns_ad;

CREATE TABLE campaigns_ad (
          id UUID PRIMARY KEY,
          organization_id VARCHAR(50) NOT NULL,
          name VARCHAR(255) NOT NULL,
          status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',

-- ذخیره مقدار Enum پروتو (CampaignAdType)
-- می‌توانیم عدد (int) یا رشته ذخیره کنیم. اینجا عدد ذخیره می‌کنیم که دقیقاً منطبق با پروتو باشد.
          campaign_type INTEGER NOT NULL,

-- تمام جزئیات (AdDetails یا EmailDetails) اینجا به صورت JSON ذخیره می‌شود
          details JSONB NOT NULL,

          scheduled_at BIGINT,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_campaign_org ON campaigns_ad(organization_id);
CREATE INDEX idx_campaign_status ON campaigns_ad(status);
CREATE INDEX idx_campaign_type ON campaigns_ad(campaign_type);
CREATE INDEX idx_campaign_details ON campaigns_ad USING gin (details);

