package config

import (
	"cmp"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"

	"github.com/isometry/yaketty/internal/library"
	"github.com/isometry/yaketty/internal/options"
	"github.com/isometry/yaketty/internal/persona"
	"github.com/isometry/yaketty/internal/scenario"
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

	scenarioLibrary := viper.GetString("scenarios")
	personaLibrary := viper.GetString("personas")

	// Load scenario file first (if scenario override is specified via flag, use that)
	scenarioToLoad := cmp.Or[string](viper.GetString("scenario"), config.Scenario.Scenario)

	if library.IsDirectPath(scenarioToLoad) {
		// Direct path - use as-is
		slog.Debug("loading scenario from direct path", slog.String("path", scenarioToLoad))
		if err := config.LoadFromFile(scenarioToLoad); err != nil {
			slog.Warn("error loading scenario from path", slog.Any("error", err))
			return nil, err
		}
	} else if scenarioToLoad != "" {
		// Library reference - add .yaml if not present
		scenarioName := scenarioToLoad
		if !strings.HasSuffix(scenarioName, ".yaml") {
			scenarioName = scenarioName + ".yaml"
		}
		scenarioFilePath := filepath.Join(scenarioLibrary, scenarioName)
		if library.FileExists(scenarioLibrary, scenarioName) {
			slog.Debug("loading scenario from library", slog.String("scenario", scenarioToLoad))
			if err := config.LoadFromFile(scenarioFilePath); err != nil {
				slog.Warn("error loading scenario", slog.Any("error", err))
				return nil, err
			}
		} else if viper.GetString("scenario") != "" {
			// Not a file, use as scenario text
			config.Scenario.Scenario = scenarioToLoad
		}
	}

	// Load persona files from config or scenario
	if config.Persona1.Persona != "" {
		if library.IsDirectPath(config.Persona1.Persona) {
			// Direct path - use as-is
			slog.Debug("loading persona1 from direct path", slog.String("path", config.Persona1.Persona))
			if err := config.Persona1.LoadFromFile(config.Persona1.Persona); err != nil {
				slog.Warn("error loading persona1 from path", slog.Any("error", err))
				return nil, err
			}
		} else {
			// Library reference - add .yaml if not present
			personaName := config.Persona1.Persona
			if !strings.HasSuffix(personaName, ".yaml") {
				personaName = personaName + ".yaml"
			}
			filePathPersona1 := filepath.Join(personaLibrary, personaName)
			if library.FileExists(personaLibrary, personaName) {
				slog.Debug("loading persona1 from library", slog.String("persona", config.Persona1.Persona))
				if err := config.Persona1.LoadFromFile(filePathPersona1); err != nil {
					slog.Warn("error loading persona1", slog.Any("error", err))
					return nil, err
				}
			}
		}
	}

	if config.Persona2.Persona != "" {
		if library.IsDirectPath(config.Persona2.Persona) {
			// Direct path - use as-is
			slog.Debug("loading persona2 from direct path", slog.String("path", config.Persona2.Persona))
			if err := config.Persona2.LoadFromFile(config.Persona2.Persona); err != nil {
				slog.Warn("error loading persona2 from path", slog.Any("error", err))
				return nil, err
			}
		} else {
			// Library reference - add .yaml if not present
			personaName := config.Persona2.Persona
			if !strings.HasSuffix(personaName, ".yaml") {
				personaName = personaName + ".yaml"
			}
			filePathPersona2 := filepath.Join(personaLibrary, personaName)
			if library.FileExists(personaLibrary, personaName) {
				slog.Debug("loading persona2 from library", slog.String("persona", config.Persona2.Persona))
				if err := config.Persona2.LoadFromFile(filePathPersona2); err != nil {
					slog.Warn("error loading persona2", slog.Any("error", err))
					return nil, err
				}
			}
		}
	}

	// Apply command-line overrides AFTER file loading
	if viper.GetString("persona1.persona") != "" {
		persona1Override := viper.GetString("persona1.persona")

		if library.IsDirectPath(persona1Override) {
			// Direct path - load it
			slog.Debug("loading persona1 override from path", slog.String("path", persona1Override))
			if err := config.Persona1.LoadFromFile(persona1Override); err != nil {
				slog.Warn("error loading persona1 override from path", slog.Any("error", err))
				return nil, err
			}
		} else {
			// Check if it's a library reference - add .yaml if not present
			personaName := persona1Override
			if !strings.HasSuffix(personaName, ".yaml") {
				personaName = personaName + ".yaml"
			}
			persona1FilePath := filepath.Join(personaLibrary, personaName)
			if library.FileExists(personaLibrary, personaName) {
				// Library file - load it
				slog.Debug("loading persona1 override from library", slog.String("persona", persona1Override))
				if err := config.Persona1.LoadFromFile(persona1FilePath); err != nil {
					slog.Warn("error loading persona1 override", slog.Any("error", err))
					return nil, err
				}
			} else {
				// Not a file, treat as persona description string
				config.Persona1.Persona = persona1Override
			}
		}
	}

	if viper.GetString("persona2.persona") != "" {
		persona2Override := viper.GetString("persona2.persona")

		if library.IsDirectPath(persona2Override) {
			// Direct path - load it
			slog.Debug("loading persona2 override from path", slog.String("path", persona2Override))
			if err := config.Persona2.LoadFromFile(persona2Override); err != nil {
				slog.Warn("error loading persona2 override from path", slog.Any("error", err))
				return nil, err
			}
		} else {
			// Check if it's a library reference - add .yaml if not present
			personaName := persona2Override
			if !strings.HasSuffix(personaName, ".yaml") {
				personaName = personaName + ".yaml"
			}
			persona2FilePath := filepath.Join(personaLibrary, personaName)
			if library.FileExists(personaLibrary, personaName) {
				// Library file - load it
				slog.Debug("loading persona2 override from library", slog.String("persona", persona2Override))
				if err := config.Persona2.LoadFromFile(persona2FilePath); err != nil {
					slog.Warn("error loading persona2 override", slog.Any("error", err))
					return nil, err
				}
			} else {
				// Not a file, treat as persona description string
				config.Persona2.Persona = persona2Override
			}
		}
	}

	// Apply other command-line overrides
	if viper.GetString("opening") != "" {
		config.OpeningPrompt = viper.GetString("opening")
	}

	if len(viper.GetStringSlice("prompts")) > 0 {
		config.ExtraPrompts = append(config.ExtraPrompts, viper.GetStringSlice("prompts")...)
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

	// Apply global model override LAST
	if viper.GetString("model") != "" {
		globalModel := viper.GetString("model")
		slog.Debug("applying global model override", slog.String("model", globalModel))
		config.Persona1.Model = globalModel
		config.Persona2.Model = globalModel
	}

	return &config, nil
}
