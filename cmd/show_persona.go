package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/isometry/yaketty/internal/library"
)

var showPersonaCmd = &cobra.Command{
	Use:   "show-persona [name]",
	Short: "Display the contents of an embedded persona",
	Long:  `Show the YAML content of an embedded persona by name (without .yaml extension).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := library.GetEmbeddedFile("personas", args[0])
		if err != nil {
			slog.Error("failed to get persona", "name", args[0], "error", err)
			os.Exit(1)
		}
		fmt.Print(string(content))
	},
}

func init() {
	rootCmd.AddCommand(showPersonaCmd)
}
