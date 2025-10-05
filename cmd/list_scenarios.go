package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/isometry/yaketty/internal/library"
)

var listScenariosCmd = &cobra.Command{
	Use:   "list-scenarios",
	Short: "List all available embedded scenarios",
	Long:  `List all scenario files that are embedded in the yaketty binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		scenarios, err := library.ListScenarios()
		if err != nil {
			slog.Error("failed to list scenarios", "error", err)
			return
		}

		for _, scenario := range scenarios {
			fmt.Println(scenario)
		}
	},
}

func init() {
	rootCmd.AddCommand(listScenariosCmd)
}
