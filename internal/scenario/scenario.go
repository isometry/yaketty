package scenario

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Scenario struct {
	Scenario string    `mapstructure:"scenario"`
	Roles    [2]string `mapstructure:"roles"`
	Opening  string    `mapstructure:"opening" default:"Welcome"`
}

func (s *Scenario) LoadFromFile(filePath string) error {
	scenarioData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(scenarioData, s)
}
