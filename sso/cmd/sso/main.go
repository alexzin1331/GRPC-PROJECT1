package main

import (
	"log/slog"
	"os"
	"os/signal"
	app "sso_module/internal/app"
	"sso_module/internal/config"
	"syscall"
)

/*
Для запуска:
перейти в sso
go run cmd/sso/main.go --config=./config/local.yaml
*/

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("start",
		slog.Any("cfg", cfg),
	) // потом убрать
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTl)
	go application.GRPCSrv.MustRun()
	//инициализировать объект конфига /у
	//инициализировать логгер /у
	//инициализировать приложение (app)
	//запустить gRPC-сервер приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Info("stopping application", slog.String("signal: ", sign.String()))
	application.GRPCSrv.Stop()

	log.Info("application stopped")
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
