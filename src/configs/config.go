// internal/config/config.go
// 📦 پکیج config مسئول بارگذاری پیکربندی سیستم از فایل یا ENV است
// تعریف نام پکیج برای استفاده در بقیه پروژه

package configs

import ( // ایمپورت‌های مورد نیاز برای پیکربندی
	"fmt"
	"github.com/spf13/viper" // کتابخانه viper برای خواندن config از فایل و ENV
	"log"                    // برای لاگ‌گرفتن خطاها و هشدارها
	"time"
) // پایان ایمپورت‌ها

// ✅ تنظیمات gRPC

type GrpcConfig struct { // ساختار تنظیمات grpc
	Address string `mapstructure:"address"` // آدرس لیسن gRPC
	Port    string `mapstructure:"port"`    // پورت gRPC
} // پایان GrpcConfig

// ✅ تنظیمات سرور HTTP / TLS
type ServerConfig struct { // ساختار تنظیمات سرور
	Address  string `mapstructure:"address"`   // آدرس لیسن سرور
	Port     string `mapstructure:"port"`      // پورت سرور
	CertFile string `mapstructure:"cert_file"` // مسیر فایل cert برای TLS
	KeyFile  string `mapstructure:"key_file"`  // مسیر فایل key برای TLS
} // پایان ServerConfig

// ✅ تنظیمات احراز هویت

type AuthConfig struct { // ساختار تنظیمات auth
	JWTSecret   string `mapstructure:"jwt_secret"`   // کلید امضای JWT
	OPAEndpoint string `mapstructure:"opa_endpoint"` // آدرس سرویس OPA برای policy
} // پایان AuthConfig

// ✅ ساختار تنظیمات گوگل OAuth

type GoogleConfig struct { // ساختار تنظیمات گوگل
	ClientID     string `mapstructure:"client_id"`     // client id گوگل
	ClientSecret string `mapstructure:"client_secret"` // client secret گوگل
	RedirectURL  string `mapstructure:"redirect_url"`  // آدرس redirect بعد از لاگین
} // پایان GoogleConfig

// ✅ تنظیمات PostgreSQL (جدید)
// این struct جایی است که مشخصات اتصال دیتابیس در آن نگه داشته می‌شود

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"` // 👈 این خط باید باشد

	// 👇 فیلدهای جدید برای Pool
	MaxConns        int32         `mapstructure:"max_conns"`         // pgx از int32 استفاده می‌کند
	MinConns        int32         `mapstructure:"min_conns"`         // pgx از int32 استفاده می‌کند
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"` // وایپر خودش رشته "1h" را به Duration تبدیل می‌کند
}

// DSN متدی برای ساختن رشته اتصال استاندارد PostgreSQL است

func (c PostgresConfig) DSN() string { // متد عضو روی PostgresConfig
	// در اینجا عمداً از fmt استفاده نکردیم تا ساده بماند؛
	// اگر خواستی می‌توانی به نسخه fmt.Sprintf برگردی.
	return "host=" + c.Host +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" port=" + fmt.Sprint(c.Port) +
		" sslmode=" + c.SSLMode +
		" TimeZone=" + c.TimeZone
} // پایان متد DSN

// ✅ ساختار نهایی Config که همه تنظیمات را کنار هم نگه می‌دارد

type Config struct { // ساختار تجمیع کل تنظیمات
	Server      ServerConfig   `mapstructure:"server"`       // تنظیمات سرور HTTP
	Auth        AuthConfig     `mapstructure:"auth"`         // تنظیمات احراز هویت
	GoogleOAuth GoogleConfig   `mapstructure:"google_oauth"` // تنظیمات گوگل OAuth
	Grpc        GrpcConfig     `mapstructure:"grpc"`         // تنظیمات gRPC
	Postgres    PostgresConfig `mapstructure:"postgresdb"`   // 🔴 تنظیمات Postgres (بخش جدید)
} // پایان Config

// Load وظیفه دارد config.yaml را بخواند و در struct Config قرار دهد

func Load() (*Config, error) { // تابع بارگذاری کانفیگ
	viper.SetConfigName("config")   // نام فایل کانفیگ بدون پسوند
	viper.SetConfigType("yaml")     // نوع فایل کانفیگ
	viper.AddConfigPath(".")        // مسیر فعلی
	viper.AddConfigPath("./config") // مسیر پوشه config
	viper.AutomaticEnv()            // خواندن مقادیر از ENV در صورت وجود

	if err := viper.ReadInConfig(); err != nil { // تلاش برای خواندن فایل config
		log.Println("⚠️ config file not found, relying on ENV variables") // هشدار در صورت نبود فایل
	} // پایان if

	var cfg Config                                // تعریف متغیر برای نگه داشتن کانفیگ نهایی
	if err := viper.Unmarshal(&cfg); err != nil { // تبدیل config خوانده شده به struct
		log.Fatalf("failed to load config: %v", err) // اگر خطا بود برنامه را متوقف کن
	} // پایان if

	return &cfg, nil // بازگشت اشاره‌گر به config (خروجی دوم را اگر جای دیگری استفاده می‌کردی، می‌توانیم بعداً سفارشی کنیم)
} // پایان Load
