package api

import (
	"net/http"

	"github.com/RascleVictor/recon-engine/internal/api/handlers"
	"github.com/RascleVictor/recon-engine/internal/plugin/registry"
)

func NewRouter(pluginRegistry *registry.Registry) http.Handler {
	mux := http.NewServeMux()

	pluginHandler := handlers.NewPluginHandler(pluginRegistry)

	mux.HandleFunc("GET /health", handlers.Health)

	mux.HandleFunc("GET /api/v1/plugins", pluginHandler.List)
	mux.HandleFunc("GET /api/v1/plugins/{id}", pluginHandler.Get)

	return mux
}
