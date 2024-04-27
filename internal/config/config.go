package config

import (
	"cmp"
	"log/slog"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"

	"github.com/isometry/yaketty/internal/options"
	"github.com/isometry/yaketty/internal/persona"
	"github.com/isometry/yaketty/internal/scenario"
	"github.com/isometry/yaketty/internal/utils"
)

type Config struct {
	scenario.Scenario `mapstructure:",squash"`
	ExtraPrompts      []string             `mapstructure:"prompts"`
	Persona1          persona.Persona      `mapstructure:"persona1"`
	Persona2          persona.Persona      `mapstructure:"persona2"`
	Options           options.ModelOptions `mapstructure:"options"`
}

func Load(path, name string) (*Config, error) {
	// if name is a valid file path, then use it, else use the default config file name
	filename := filepath.Join(path, name)
	slog.Debug("loading config", "filename", filename)

	if _, err := os.Stat(filename); err == nil {
		viper.SetConfigFile(filename)
	} else {
		viper.AddConfigPath(path)
		viper.SetConfigName(name)
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config", "error", err)
		return nil, err
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("failed to unmarshal config", "error", err)
		return nil, err
	}

	slog.Debug("config before defaults", slog.Any("config", config))

	defaults.SetDefaults(&config)

	slog.Debug("config after defaults", slog.Any("config", config))

	// override config values with command line flags
	config.Scenario.Scenario = cmp.Or[string](viper.GetString("scenario"), config.Scenario.Scenario)

	scenarioLibrary := viper.GetString("scenarios")
	personaLibrary := viper.GetString("personas")

	scenarioFilePath := filepath.Join(scenarioLibrary, config.Scenario.Scenario)
	if utils.IsFilePath(scenarioFilePath) {
		slog.Debug("loading scenario from file", slog.String("scenario", config.Scenario.Scenario))
		if err := config.Scenario.LoadFromFile(scenarioFilePath); err != nil {
			slog.Warn("error loading scenario", slog.Any("error", err))
			return nil, err
		}
	}

	filePathPersona1 := filepath.Join(personaLibrary, config.Persona1.Persona)
	if utils.IsFilePath(filePathPersona1) {
		if err := config.Persona1.LoadFromFile(filePathPersona1); err != nil {
			slog.Warn("error loading persona1", slog.Any("error", err))
			return nil, err
		}
	}

	filePathPersona2 := filepath.Join(personaLibrary, config.Persona2.Persona)
	if utils.IsFilePath(filePathPersona2) {
		if err := config.Persona2.LoadFromFile(filePathPersona2); err != nil {
			slog.Warn("error loading persona2", slog.Any("error", err))
			return nil, err
		}
	}

	// set default names
	config.Persona1.Name = cmp.Or[string](config.Persona1.Name, "Jane")
	config.Persona2.Name = cmp.Or[string](config.Persona2.Name, "John")

	// merge shared options into persona options
	if err := mergo.Merge(&config.Persona1.Options, &config.Options); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&config.Persona2.Options, &config.Options); err != nil {
		return nil, err
	}

	return &config, nil
}
