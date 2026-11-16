package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RascleVictor/recon-engine/internal/plugin/manifest"
)

var pluginIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

func Validate(m *manifest.Manifest) error {
	if m.APIVersion != "recon.framework/v1" {
		return fmt.Errorf("unsupported API version: %s", m.APIVersion)
	}

	if !pluginIDPattern.MatchString(m.Plugin.ID) {
		return fmt.Errorf("invalid plugin ID: %s", m.Plugin.ID)
	}

	if strings.TrimSpace(m.Plugin.Name) == "" {
		return fmt.Errorf("plugin name is required")
	}

	if strings.TrimSpace(m.Plugin.Version) == "" {
		return fmt.Errorf("plugin version is required")
	}

	if m.Runtime.Type != "docker" {
		return fmt.Errorf("unsupported runtime type: %s", m.Runtime.Type)
	}

	if strings.TrimSpace(m.Runtime.Image) == "" {
		return fmt.Errorf("runtime image is required")
	}

	if len(m.Runtime.Command) == 0 {
		return fmt.Errorf("runtime command is required")
	}

	if m.Output.Format == "" {
		return fmt.Errorf("output format is required")
	}

	if m.Output.Source != "stdout" && m.Output.Source != "file" {
		return fmt.Errorf("unsupported output source: %s", m.Output.Source)
	}

	for _, input := range m.Inputs {
		if strings.TrimSpace(input.Name) == "" {
			return fmt.Errorf("input name is required")
		}

		switch input.Type {
		case "string", "integer", "number", "boolean", "domain", "ip":
		default:
			return fmt.Errorf(
				"unsupported input type %q for %s",
				input.Type,
				input.Name,
			)
		}
	}

	return nil
}
