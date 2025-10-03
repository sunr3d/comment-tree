package entrypoint

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/comment-tree/internal/config"
	httphandlers "github.com/sunr3d/comment-tree/internal/handlers"
	"github.com/sunr3d/comment-tree/internal/infra/postgres"
)

func Run(cfg *config.Config) error {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инфраслой
	repo, err := postgres.New(appCtx, cfg.DB)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("posgres.New")
		return fmt.Errorf("postgres.New(): %w", err)
	}
	// Сервисный слой
	// TODO: svc := ...

	// REST API (HTTP) + Middleware
	h := httphandlers.New(svc)
	engine := h.RegisterHandlers()

	// Server
	return engine.Run(":" + cfg.HTTPPort)
}
