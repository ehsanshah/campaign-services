-- فعال‌سازی افزونه UUID برای تولید ID خودکار (اگر از سمت Go تولید نکنیم)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE campaigns (
                           id UUID PRIMARY KEY,
                           account_id UUID NOT NULL,
                           name VARCHAR(255) NOT NULL,
                           status VARCHAR(50) NOT NULL,
                           type_for_humans VARCHAR(100),

    -- ⚡ فیلدهای JSONB برای داده‌های تو در تو
                           recipients JSONB DEFAULT '{}',
                           options JSONB DEFAULT '{}',
                           stats JSONB DEFAULT '{}',         -- اسنپ‌شات آمار
                           filters JSONB DEFAULT '[]',       -- آرایه فیلترها
                           extra_fields JSONB DEFAULT '{}',

    -- زمان‌بندی‌ها
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           scheduled_for TIMESTAMP WITH TIME ZONE,
                           started_at TIMESTAMP WITH TIME ZONE,
                           finished_at TIMESTAMP WITH TIME ZONE,
                           stopped_at TIMESTAMP WITH TIME ZONE,
                           winner_selected_at TIMESTAMP WITH TIME ZONE,

    -- فلگ‌های کنترلی (Boolean)
                           is_stopped BOOLEAN DEFAULT FALSE,
                           is_currently_sending_out BOOLEAN DEFAULT FALSE,
                           can_be_scheduled BOOLEAN DEFAULT TRUE,
                           has_winner BOOLEAN DEFAULT FALSE,

    -- اطلاعات A/B Testing
                           winner_version_for_human VARCHAR(100),
                           winner_sending_time_for_humans VARCHAR(100),

    -- آرایه‌های متنی (در Postgres آرایه نیتیو داریم)
                           email_ids TEXT[],
                           default_email_id UUID,
                           warnings TEXT[],
                           used_in_automations BOOLEAN DEFAULT FALSE
);

-- ایجاد ایندکس برای سرعت بالا در کوئری‌ها
CREATE INDEX idx_campaigns_account_id ON campaigns(account_id);
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_scheduled_for ON campaigns(scheduled_for);