package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"reqium/internal/app"
)

func newRunCommand(deps dependencies) *cobra.Command {
	var envName string
	var timeoutSec int

	cmd := &cobra.Command{
		Use:   "run COLLECTION",
		Short: "Run a saved collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runner := app.NewRunnerService(deps.collectionRepo, deps.requestService, deps.resolver, deps.envService)
			results, err := runner.Run(cmd.Context(), args[0], envName, time.Duration(timeoutSec)*time.Second)
			if err != nil {
				return err
			}
			for _, result := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s  %s  %d  %s  %s\n", result.Status, result.Method, result.URL, result.StatusCode, result.Duration, result.RequestName)
				if result.Error != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "  error: %s\n", result.Error)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&envName, "env", "", "environment to resolve {{variables}}")
	cmd.Flags().IntVarP(&timeoutSec, "timeout", "t", 30, "request timeout in seconds")
	return cmd
}
