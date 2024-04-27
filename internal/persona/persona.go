package persona

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/isometry/yaketty/internal/options"
)

type Persona struct {
	Model   string               `mapstructure:"model" default:"llama3"`
	Name    string               `mapstructure:"name"`
	Persona string               `mapstructure:"persona"`
	Prompts []string             `mapstructure:"prompts"`
	Options options.ModelOptions `mapstructure:"options"`
}

func (p *Persona) LoadFromFile(filePath string) error {
	personaData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(personaData, p)
}
