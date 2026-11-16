package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RascleVictor/recon-engine/internal/plugin/registry"
)

type PluginHandler struct {
	registry *registry.Registry
}

func NewPluginHandler(r *registry.Registry) *PluginHandler {
	return &PluginHandler{
		registry: r,
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
