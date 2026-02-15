package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ehsanshah/campaign-services/src/configs"
	"github.com/ehsanshah/campaign-services/src/internal/app"
)

func main() {
	// =========================================================================
	// 1. Load Configuration
	// =========================================================================
	// بارگذاری تنظیمات از فایل config.yaml یا متغیرهای محیطی
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configurations: %v", err)
	}
	log.Printf("✅ Config loaded successfully")

	// =========================================================================
	// 2. Initialize Application (Wiring)
	// =========================================================================
	// ساخت کل ساختار برنامه (دیتابیس، سرویس‌ها، هندلرها) در لایه app
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize app: %v", err)
	}

	// =========================================================================
	// 3. Run Server
	// =========================================================================
	// اجرای سرور در یک Goroutine جداگانه تا برنامه بلاک نشود
	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("❌ Server runtime error: %v", err)
		}
	}()

	// =========================================================================
	// 4. Graceful Shutdown
	// =========================================================================
	// منتظر سیگنال سیستم عامل (مثل CTRL+C یا سیگنال داکر) می‌مانیم
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// برنامه اینجا متوقف می‌شود تا سیگنال برسد
	sig := <-quit
	log.Printf("⚠️ Signal received: %v. Shutting down...", sig)

	// فراخوانی متد Shutdown در لایه app برای بستن کانکشن‌ها
	application.Shutdown()

	log.Println("👋 Server stopped gracefully")
}
