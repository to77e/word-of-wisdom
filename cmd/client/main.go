package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/to77e/word-of-wisdom/internal/client"
	"github.com/to77e/word-of-wisdom/internal/config"
	"github.com/to77e/word-of-wisdom/internal/proofofwork"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	configPath     = "config.yaml"
	defaultClients = 100
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

	wg := sync.WaitGroup{}
	wg.Add(defaultClients)
	for i := 1; i <= defaultClients; i++ {
		go func(n int) {
			start := time.Now()

			tcpClient, err := client.New(fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port))
			if err != nil {
				slog.With("error", err.Error()).Error("new client")
				return
			}
			defer func() {
				wg.Done()
				if err = tcpClient.Close(); err != nil {
					slog.With("error", err.Error()).Error("close tcp connection")
				}
			}()

			slog.With("number", n).Debug("start")
			if err = tcpClient.WriteChallengeRequest(); err != nil {
				slog.With("error", err.Error()).Error("write challenge request")
				return
			}

			challenge, err := tcpClient.ReadChallengeResponse()
			if err != nil {
				slog.With("error", err.Error()).Error("read challenge response")
				return
			}

			// 4. solve
			pow := proofofwork.New(challenge.Difficulty)
			pow.SetChallenge(challenge.Content)
			if err = pow.ComputeSolution(); err != nil {
				slog.With("error", err.Error()).Error("compute solution")
				return
			}

			if err = tcpClient.WriteSolutionRequest(pow.GetSolution()); err != nil {
				slog.With("error", err.Error()).Error("write solution request")
				return
			}

			result, err := tcpClient.ReadSolutionResponse()
			if err != nil {
				slog.With("error", err.Error()).Error("read solution response")
				return
			}

			fmt.Printf("%d\n%s\n\n", n, result)
			slog.With("number", n, "seconds", time.Since(start).Seconds()).Debug("time elapsed")
		}(i)
	}
	wg.Wait()
}
