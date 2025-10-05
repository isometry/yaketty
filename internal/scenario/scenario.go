package scenario

import (
	"go.yaml.in/yaml/v4"

	"github.com/isometry/yaketty/internal/library"
)

type Scenario struct {
	Scenario      string    `mapstructure:"scenario"`
	Roles         [2]string `mapstructure:"roles"`
	OpeningPrompt string    `mapstructure:"opening_prompt" default:"Start the conversation with an appropriate greeting or opening statement for this scenario"`
}

func (s *Scenario) LoadFromFile(filePath string) error {
	// Use library.ReadFileOrPath which tries:
	// 1. Local filesystem at the given path
	// 2. Embedded filesystem at the given path
	scenarioData, err := library.ReadFileOrPath(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(scenarioData, s)
}
