-- migrations/campaign/000002_templates_init.up.sql

CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    current_version_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS template_versions (
                                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES templates(id) ON DELETE CASCADE,
    version_label VARCHAR(50),
    subject VARCHAR(255),
    html_content TEXT,
    plain_text TEXT,
    language VARCHAR(10) DEFAULT 'en', -- ✅ زبان قالب
    tags TEXT[] DEFAULT '{}',          -- ✅ تگ‌ها (مانند: ["marketing", "newsletter"])
    metadata JSONB DEFAULT '{}',       -- ✅ اطلاعات Extra (مانند: کد رنگ‌ها، نام طراح و غیره)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ایجاد ایندکس GIN برای جستجوی سریع روی تگ‌ها و متادیتا
CREATE INDEX idx_templates_tags ON template_versions USING GIN (tags);
CREATE INDEX idx_templates_metadata ON template_versions USING GIN (metadata);