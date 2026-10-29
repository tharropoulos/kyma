package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/museslabs/kyma/internal/tui"
)

var gotoCmd = &cobra.Command{
	Use:   "goto <slide>",
	Short: "Show a slide, by its number, in the running presentation",
	Long: "Show a slide, by its number counting from 1, in the presentation running on this machine.\n" +
		"Bind it to a remote, such as a KDE Connect command, to jump around without the keyboard.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		slideNum, err := strconv.Atoi(args[0])
		if err != nil || slideNum < 1 {
			return fmt.Errorf("expected a slide number from 1 up, got: %v", args[0])
		}

		client, err := tui.NewSyncClient()
		if err != nil {
			return fmt.Errorf("no presentation is running: %w", err)
		}
		defer client.Close()

		return client.GoTo(slideNum)
	},
}
