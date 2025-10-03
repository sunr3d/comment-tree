package main

import (
	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/comment-tree/internal/config"
	"github.com/sunr3d/comment-tree/internal/entrypoint"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("запуск приложения...")

	zlog.Logger.Info().Msg("загрузка конфига...")
	cfg, err := config.GetConfig("config.yml")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("config.GetConfig")
	}
	zlog.Logger.Info().Msg("конфиг успешно загружен...")

	if err := entrypoint.Run(cfg); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("entrypoint.Run")
	}
}
