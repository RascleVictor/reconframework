package loader

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RascleVictor/recon-engine/internal/plugin/manifest"
	"github.com/RascleVictor/recon-engine/internal/plugin/validator"
	"gopkg.in/yaml.v3"
)

func Discover(directory string) ([]manifest.Manifest, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read plugin directory: %w", err)
	}

	var plugins []manifest.Manifest

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(
			directory,
			entry.Name(),
			"plugin.yaml",
		)

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var plugin manifest.Manifest

		decoder := yaml.NewDecoder(
			bytes.NewReader(data),
		)
		decoder.KnownFields(true)

		if err := decoder.Decode(&plugin); err != nil {
			return nil, fmt.Errorf(
				"decode plugin %s: %w",
				entry.Name(),
				err,
			)
		}

		if err := validator.Validate(&plugin); err != nil {
			return nil, fmt.Errorf(
				"validate plugin %s: %w",
				entry.Name(),
				err,
			)
		}

		if plugin.Plugin.ID != entry.Name() {
			return nil, fmt.Errorf(
				"plugin ID %q does not match directory %q",
				plugin.Plugin.ID,
				entry.Name(),
			)
		}

		plugins = append(plugins, plugin)
	}

	return plugins, nil
}
