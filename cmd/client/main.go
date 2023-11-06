package main

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/to77e/word-of-wisdom/internal/client"
	"github.com/to77e/word-of-wisdom/internal/proofofwork"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	logLevel       = slog.LevelInfo
	defaultAddress = "localhost:11001"
	defaultClients = 10
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	validator.New()

	wg := sync.WaitGroup{}
	wg.Add(defaultClients)
	for i := 1; i <= defaultClients; i++ {
		go func(n int) {
			start := time.Now()
			tcpClient, err := client.New(defaultAddress)
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
			pow := proofofwork.New()
			pow.SetChallenge(challenge)
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

			slog.Info(result)
			slog.With("number", n, "seconds", time.Since(start).Seconds()).Debug("time elapsed")
		}(i)
	}
	wg.Wait()
}
