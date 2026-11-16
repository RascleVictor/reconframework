package manifest

type Manifest struct {
	APIVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Plugin     PluginInfo `yaml:"plugin" json:"plugin"`
	Runtime    Runtime    `yaml:"runtime" json:"runtime"`
	Inputs     []Input    `yaml:"inputs" json:"inputs"`
	Output     Output     `yaml:"output" json:"output"`
}

type PluginInfo struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
	Category    string `yaml:"category" json:"category"`
	Description string `yaml:"description" json:"description"`
}

type Runtime struct {
	Type    string   `yaml:"type" json:"type"`
	Image   string   `yaml:"image" json:"image"`
	Command []string `yaml:"command" json:"command"`
	Timeout string   `yaml:"timeout" json:"timeout"`
}

type Input struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
	Default  any    `yaml:"default" json:"default"`
}

type Output struct {
	Format string `yaml:"format" json:"format"`
	Source string `yaml:"source" json:"source"`
}
