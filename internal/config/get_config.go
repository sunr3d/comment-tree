package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/wb-go/wbf/config"
)

func GetConfig(path string) (*Config, error) {
	cfg := config.New()
	if err := cfg.Load(path); err != nil {
	}

	cfg.SetDefault("HTTP_PORT", "8080")
	cfg.SetDefault("BASE_URL", "http://localhost:8080")
	cfg.SetDefault("LOG_LEVEL", "info")

	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cfg.Unmarshal: %w", err)
	}

	// fallback (костыль)
	if strings.TrimSpace(c.DB.DSN) == "" {
		if dsn := strings.TrimSpace(os.Getenv("DB_DSN")); dsn != "" {
			c.DB.DSN = dsn
		}
	}
	if strings.TrimSpace(c.DB.DSN) == "" {
		return nil, fmt.Errorf("DB.DSN не может быть пустым")
	}

	return &c, nil
}
