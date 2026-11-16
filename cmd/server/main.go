package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RascleVictor/recon-engine/internal/api"
	"github.com/RascleVictor/recon-engine/internal/config"
	"github.com/RascleVictor/recon-engine/internal/plugin/executor"
	"github.com/RascleVictor/recon-engine/internal/plugin/loader"
	"github.com/RascleVictor/recon-engine/internal/plugin/registry"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		logger.Error("configuration error", "error", err)
		os.Exit(1)
	}

	pluginRegistry := registry.New()

	if cfg.Plugins.AutoDiscovery {
		plugins, err := loader.Discover(cfg.Plugins.Directory)
		if err != nil {
			logger.Error("plugin discovery failed", "error", err)
			os.Exit(1)
		}

		for _, plugin := range plugins {
			if err := pluginRegistry.Register(plugin); err != nil {
				logger.Error("plugin registration failed", "error", err)
				os.Exit(1)
			}

			logger.Info(
				"plugin registered",
				"id", plugin.Plugin.ID,
				"name", plugin.Plugin.Name,
			)
		}
	}
	dockerExecutor, err := executor.NewDockerExecutor()
	if err != nil {
		logger.Error("docker executor initialization failed", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(pluginRegistry, dockerExecutor)

	address := fmt.Sprintf(
		"%s:%d",
		cfg.Server.Host,
		cfg.Server.Port,
	)

	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info("starting recon engine", "address", address)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}

	logger.Info("recon engine stopped")
}
