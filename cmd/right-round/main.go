package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	rightround "github.com/cboone/right-round"
	"github.com/cboone/right-round/internal/data"
	"github.com/cboone/right-round/internal/tui"
)

var version = "dev"

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var (
		typeFlag  string
		groupFlag string
	)

	cmd := &cobra.Command{
		Use:     "right-round",
		Short:   "Browse and preview terminal progress indicators",
		Version: version,
		Long: `right-round is a TUI for browsing 433 terminal progress indicators
(333 spinners, 100 progress bars) from 26 open-source collections.

Navigate, preview live animations, and copy entries as JSON.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate --type flag
			if typeFlag != "" && typeFlag != "spinner" && typeFlag != "progress_bar" {
				return fmt.Errorf("invalid --type %q: must be \"spinner\" or \"progress_bar\"", typeFlag)
			}

			jsonBytes := rightround.EmbeddedCatalogJSON()
			grouped, err := data.LoadCatalog(jsonBytes)
			if err != nil {
				return fmt.Errorf("loading catalog: %w", err)
			}

			model := tui.New(grouped, typeFlag, groupFlag)
			p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("running TUI: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&typeFlag, "type", "", "lock to \"spinner\" or \"progress_bar\"")
	cmd.Flags().StringVar(&groupFlag, "group", "", "start with a specific group selected")

	return cmd
}
