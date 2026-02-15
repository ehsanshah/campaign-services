package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool" // درایور جدید

	"github.com/ehsanshah/campaign-services/src/configs"
	grpcHandler "github.com/ehsanshah/campaign-services/src/internal/adapter/handler/grpc"
	"github.com/ehsanshah/campaign-services/src/internal/adapter/storage/postgres"
	"github.com/ehsanshah/campaign-services/src/internal/service"

	// مسیر کدهای جنریت شده پروتو
	pb "github.com/ehsanshah/campaign-services/src/pkg/pb/camp/v1"

	// پکیج اتصال دیتابیس که ساختیم
	pkgPostgres "github.com/ehsanshah/campaign-services/src/internal/adapter/storage/postgres"
)

// App تمام وابستگی‌های سطح بالای سرویس را نگه می‌دارد
type App struct {
	Cfg        *configs.Config
	GRPCServer *grpc.Server
	DB         *pgxpool.Pool // استفاده از Pool قدرتمند pgx
}

// NewApp وظیفه سیم‌کشی (Wiring) و Dependency Injection را دارد
func NewApp(cfg *configs.Config) (*App, error) {

	// 1. اتصال به دیتابیس (با استفاده از پکیج pkg/postgres)
	dbPool, err := pkgPostgres.NewConnection(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}

	// 2. راه‌اندازی لایه‌ها (Repo -> Service -> Handler)

	// مخزن داده (Repository)
	campaignRepo := postgres.NewCampaignRepo(dbPool)

	// بیزینس لاجیک (Service)
	campaignService := service.NewCampaignService(campaignRepo)

	// هندلر gRPC
	campaignHandler := grpcHandler.NewServer(campaignService)

	// 3. راه‌اندازی سرور gRPC
	grpcServer := grpc.NewServer()

	// ثبت سرویس با نام جدید CampaignServiceAd
	pb.RegisterCampaignServiceAdServer(grpcServer, campaignHandler)

	// فعال‌سازی Reflection (برای ابزارهایی مثل Postman/gRPCurl)
	reflection.Register(grpcServer)

	// بازگرداندن ساختار App
	return &App{
		Cfg:        cfg,
		GRPCServer: grpcServer,
		DB:         dbPool,
	}, nil
}

// Run سرور را روی پورت مشخص شده اجرا می‌کند (Blocking)
func (a *App) Run() error {
	// ساخت آدرس پورت (مثلا :50054)
	port := fmt.Sprintf(":%s", a.Cfg.Grpc.Port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	log.Printf("🚀 Campaign Service (Ad/MTA) is running on port %s", port)

	// شروع سرویس‌دهی
	return a.GRPCServer.Serve(lis)
}

// Shutdown منابع را به صورت امن آزاد می‌کند
func (a *App) Shutdown() {
	log.Println("🛑 Stopping gRPC Server...")
	a.GRPCServer.GracefulStop()

	log.Println("🔌 Closing Database Connection Pool...")
	a.DB.Close() // بستن کانکشن‌های pgx
}
