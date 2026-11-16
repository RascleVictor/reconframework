package registry

import (
	"fmt"
	"sync"

	"github.com/RascleVictor/recon-engine/internal/plugin/manifest"
)

type Registry struct {
	mu      sync.RWMutex
	plugins map[string]manifest.Manifest
}

func New() *Registry {
	return &Registry{
		plugins: make(map[string]manifest.Manifest),
	}
}

func (r *Registry) Register(plugin manifest.Manifest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := plugin.Plugin.ID

	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("plugin already registered: %s", id)
	}

	r.plugins[id] = plugin

	return nil
}

func (r *Registry) Get(id string) (manifest.Manifest, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[id]

	return plugin, exists
}

func (r *Registry) List() []manifest.Manifest {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]manifest.Manifest, 0, len(r.plugins))

	for _, plugin := range r.plugins {
		result = append(result, plugin)
	}

	return result
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.plugins)
}
