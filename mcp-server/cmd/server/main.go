package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/app"
	"github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("config: %+v", conf)

	app, err := app.NewApplication(conf)
	if err != nil {
		log.Fatalf("failed to create application instance: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		log.Fatalf("application failed while stopped: %v", err)
	}

	log.Printf("application finished gracefully")
}
