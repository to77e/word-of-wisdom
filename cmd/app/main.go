package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/to77e/word-of-wisdom/internal/server"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	defaultAddress = "localhost:11001"
	logLevel       = slog.LevelInfo
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	tcpServer, err := server.New(defaultAddress)
	if err != nil {
		slog.With("error", err.Error()).Error("new tcp server")
		return
	}

	validator.New()

	tcpServer.Start()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	tcpServer.Stop()
}
