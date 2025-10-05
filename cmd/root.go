package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/isometry/yaketty/internal/config"
	"github.com/isometry/yaketty/internal/dialogue"
)

var rootCmd = &cobra.Command{
	Use:          "yaketty [config]",
	Args:         cobra.ExactArgs(1),
	Short:        "A CLI for driving conversational AI models",
	Long:         `Yaketty orchestrates conversational AI dialogues between two personas using the Ollama API. It enables scripted conversations where two AI models play different characters in various scenarios.`,
	PreRunE:      Load,
	RunE:         Run,
	SilenceUsage: true,
}

var (
	cfg       *config.Config
	path      string
	verbosity int
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flagSet := rootCmd.Flags()
	flagSet.StringP("scenario", "s", "", "Override the scenario for the dialogue")
	_ = viper.BindPFlag("scenario", flagSet.Lookup("scenario"))

	flagSet.StringP("scenarios", "S", "scenarios", "The path to library scenarios")
	_ = viper.BindPFlag("scenarios", flagSet.Lookup("scenarios"))

	flagSet.StringSliceP("prompts", "p", nil, "Additional system prompts for the dialogue")
	_ = viper.BindPFlag("prompts", flagSet.Lookup("prompts"))

	flagSet.StringP("persona1", "1", "", "Override the persona for the first bot")
	_ = viper.BindPFlag("persona1.persona", flagSet.Lookup("persona1"))

	flagSet.StringP("persona2", "2", "", "Override the persona for the second bot")
	_ = viper.BindPFlag("persona2.persona", flagSet.Lookup("persona2"))

	flagSet.StringP("personas", "P", "personas", "The path to library personas")
	_ = viper.BindPFlag("personas", flagSet.Lookup("personas"))

	flagSet.StringP("opening", "o", "", "Opening prompt for the first persona")
	_ = viper.BindPFlag("opening", flagSet.Lookup("opening"))

	flagSet.StringP("model", "m", "", "Override the model for both personas")
	_ = viper.BindPFlag("model", flagSet.Lookup("model"))

	flagSet.StringVarP(&path, "path", "c", ".", "The path to the configuration file")

	flagSet.Int("context", 8192, "size of the context window used to generate the next token")
	_ = viper.BindPFlag("options.num_ctx", flagSet.Lookup("context"))

	flagSet.Int("repeat-last-n", -1, "")
	_ = viper.BindPFlag("options.repeat_last_n", flagSet.Lookup("repeat-last-n"))

	flagSet.Float32("repeat-penalty", 1.1, "")
	_ = viper.BindPFlag("options.repeat_penalty", flagSet.Lookup("repeat-penalty"))

	flagSet.Float32("temperature", 0.8, "temperature of the model: increase to answer more creatively")
	_ = viper.BindPFlag("options.temperature", flagSet.Lookup("temperature"))

	flagSet.StringSlice("stop", []string{}, "stop tokens to end the conversation")
	_ = viper.BindPFlag("options.stop", flagSet.Lookup("stop"))

	flagSet.Int("top-k", 40, "a higher value (e.g. 100) will give more diverse answers, while a lower value (e.g. 10) will be more conservative")
	_ = viper.BindPFlag("options.top_k", flagSet.Lookup("top-k"))

	flagSet.Float32("top-p", 0.9, "a higher value (e.g., 0.95) will lead to more diverse text, while a lower value (e.g., 0.5) will generate more focused and conservative text")
	_ = viper.BindPFlag("options.top_p", flagSet.Lookup("top-p"))

	flagSet.CountVarP(&verbosity, "verbosity", "v", "Increase verbosity (can be used multiple times)")
	_ = viper.BindPFlag("verbosity", flagSet.Lookup("verbosity"))
}

func Load(cmd *cobra.Command, args []string) (err error) {
	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn - slog.Level(verbosity*4))

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stderr, handlerOpts)

	slog.SetDefault(slog.New(handler))

	cfg, err = config.Load(path, args[0])

	return err
}

func Run(cmd *cobra.Command, args []string) error {
	chat, err := dialogue.NewDialogue(cmd.Context(), cfg)
	if err != nil {
		return err
	}

	return chat.Start()
}
