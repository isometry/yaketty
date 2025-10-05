package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/isometry/yaketty/internal/library"
)

func listPersonasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-personas",
		Short: "List all available embedded personas",
		Long:  `List all persona files that are embedded in the yaketty binary.`,
		Run: func(cmd *cobra.Command, args []string) {
			personas, err := library.ListPersonas()
			if err != nil {
				slog.Error("failed to list personas", "error", err)
				return
			}

			for _, persona := range personas {
				fmt.Println(persona)
			}
		},
	}
}
