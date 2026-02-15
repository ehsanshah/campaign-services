package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ehsanshah/campaign-services/src/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg configs.PostgresConfig) (*pgxpool.Pool, error) {
	// 1. هندل کردن تایم‌زون
	timeZone := cfg.TimeZone
	if timeZone == "" {
		timeZone = "UTC"
	}

	// 2. ساخت رشته اتصال (DSN)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		timeZone,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing db config: %w", err)
	}

	// 3. ✅ اصلاح شده: استفاده مستقیم از مقدار Duration (چون در کانفیگ تبدیل شده است)
	if cfg.MaxConnLifetime != 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	} else {
		poolConfig.MaxConnLifetime = time.Hour // مقدار پیش‌فرض
	}

	// سایر تنظیمات Pool
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns

	// اتصال و پینگ
	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating db pool: %w", err)
	}

	if err := connPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	return connPool, nil
}
