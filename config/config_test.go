package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	wantPort := 8088
	t.Setenv("PORT", fmt.Sprint(wantPort))
	cfg, err := NewConfig()

	if err != nil {
		t.Errorf("want no error, but got %v", err)
	}

	wantEnv := "dev"
	if cfg.Env != wantEnv {
		t.Errorf("env: want %v, but got %v", wantEnv, cfg.Env)
	}

	if cfg.Port != wantPort {
		t.Errorf("port: want %d, but got %d", wantPort, cfg.Port)
	}
}
