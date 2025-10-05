package persona

import (
	"go.yaml.in/yaml/v4"

	"github.com/isometry/yaketty/internal/library"
	"github.com/isometry/yaketty/internal/options"
)

type Persona struct {
	Model   string               `mapstructure:"model" default:"gemma3"`
	Name    string               `mapstructure:"name"`
	Persona string               `mapstructure:"persona"`
	Prompts []string             `mapstructure:"prompts"`
	Options options.ModelOptions `mapstructure:"options"`
}

func (p *Persona) LoadFromFile(filePath string) error {
	// Use library.ReadFileOrPath which tries:
	// 1. Local filesystem at the given path
	// 2. Embedded filesystem at the given path
	personaData, err := library.ReadFileOrPath(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(personaData, p)
}
