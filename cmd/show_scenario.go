package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/isometry/yaketty/internal/library"
)

var showScenarioCmd = &cobra.Command{
	Use:   "show-scenario [name]",
	Short: "Display the contents of an embedded scenario",
	Long:  `Show the YAML content of an embedded scenario by name (without .yaml extension).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := library.GetEmbeddedFile("scenarios", args[0])
		if err != nil {
			slog.Error("failed to get scenario", "name", args[0], "error", err)
			os.Exit(1)
		}
		fmt.Print(string(content))
	},
}

func init() {
	rootCmd.AddCommand(showScenarioCmd)
}
