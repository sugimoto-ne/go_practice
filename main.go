package main

import (
	"context"
	"log"
	"os"

	"github.com/sugimoto-ne/go_practice.git/config"
	"github.com/sugimoto-ne/go_practice.git/logger"
)

func main() {

	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}

}

func run(ctx context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	systemLogger, err := logger.NewLogger("./logs", cfg.Env, "sys", "json")
	if err != nil {
		return err
	}
	systemLogger.Logger.Info("start run")

	if err != nil {
		return err
	}

	mux, err := NewMux(cfg)
	if err != nil {
		return err
	}

	server := NewServer(cfg, mux)

	serverStopErr := server.Run(ctx)
	systemLogger.Logger.Error("shutdown error", serverStopErr)
	return serverStopErr
}
