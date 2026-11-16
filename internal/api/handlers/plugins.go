package handlers

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/RascleVictor/recon-engine/internal/plugin/executor"
	"github.com/RascleVictor/recon-engine/internal/plugin/registry"
)

type PluginHandler struct {
	registry *registry.Registry
	executor executor.Executor
}

func NewPluginHandler(
	r *registry.Registry,
	e executor.Executor,
) *PluginHandler {
	return &PluginHandler{
		registry: r,
		executor: e,
	}
}

func (h *PluginHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
		"count":   h.registry.Count(),
		"plugins": h.registry.List(),
	}

	_ = json.NewEncoder(w).Encode(response)
}

func (h *PluginHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	plugin, exists := h.registry.Get(id)
	if !exists {
		http.Error(w, "plugin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(plugin)
}

func (h *PluginHandler) Execute(w http.ResponseWriter, r *http.Request) {
	expectedKey := os.Getenv("RECON_API_KEY")
	providedKey := r.Header.Get("X-API-Key")

	if expectedKey == "" {
		http.Error(w, "plugin execution is not configured", http.StatusServiceUnavailable)
		return
	}

	if subtle.ConstantTimeCompare(
		[]byte(providedKey),
		[]byte(expectedKey),
	) != 1 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	plugin, exists := h.registry.Get(id)
	if !exists {
		http.Error(w, "plugin not found", http.StatusNotFound)
		return
	}

	result, err := h.executor.Execute(r.Context(), plugin)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusGatewayTimeout
		}

		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":  err.Error(),
			"result": result,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
