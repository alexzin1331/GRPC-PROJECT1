package sso

import (
	"log/slog"
	"os"
	app "sso_module/internal/app"
	"sso_module/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("start",
		slog.Any("cfg", cfg),
	) // потом убрать
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)

	application.GRPCSrv.MustRun()
	//инициализировать объект конфига /у
	//инициализировать логгер /у
	//инициализировать приложение (app)
	//запустить gRPC-сервер приложения

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

/*Тут передаем переменную окружения и фиксируем логи:
 */
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
