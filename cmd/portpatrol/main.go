package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/containeroo/portpatrol/internal/checker"
	"github.com/containeroo/portpatrol/internal/config"
	"github.com/containeroo/portpatrol/internal/logger"
	"github.com/containeroo/portpatrol/internal/runner"
)

const version = "0.0.10"

// run is the main function of the application
func run(ctx context.Context, getenv func(string) string, output io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseConfig(getenv)
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	cfg.Version = version

	logger := logger.SetupLogger(cfg, output)

	targetChecker, err := checker.NewChecker(cfg.CheckType, cfg.TargetName, cfg.TargetAddress, cfg.DialTimeout, getenv)
	if err != nil {
		return fmt.Errorf("failed to initialize checker: %w", err)
	}

	return runner.LoopUntilReady(ctx, cfg.Interval, targetChecker, logger)
}

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Getenv, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
