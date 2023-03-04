package logger

import (
	"os"
	"testing"

	"github.com/sugimoto-ne/go_practice.git/config"
)

func TestNewLogger(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}

	testLogger, err := NewLogger("logs", cfg.Env, "test", "json")
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}

	testLogger.Logger.Info("test")

	filePath := "../logs/test/test.json"
	dirPath := "../logs/test"
	err = os.Remove(filePath)
	if err != nil {
		t.Fatalf("want to delete %s, but got error: %v", filePath, err)
	}

	err = os.Remove(dirPath)
	if err != nil {
		t.Fatalf("want to delete %s, but got error: %v", dirPath, err)
	}
}
