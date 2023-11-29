package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/to77e/word-of-wisdom/internal/config"
	"github.com/to77e/word-of-wisdom/internal/server"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	configPath = "config.yaml"
)

func main() {
	if err := config.ReadFile(configPath); err != nil {
		slog.With("error", err.Error()).Error("init configuration")
		os.Exit(1)
	}

	cfg := config.GetInstance()

	validator.New()

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.Level(cfg.Project.LogLevel),
				},
			),
		),
	)

	tcpServer, err := server.New(fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port))
	if err != nil {
		slog.With("error", err.Error()).Error("new tcp server")
		return
	}

	tcpServer.Start()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	tcpServer.Stop()
}
