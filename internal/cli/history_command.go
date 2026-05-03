package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	respformatter "reqium/internal/implementations/formatter"
	"reqium/internal/models"
)

func newHistoryCommand(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Inspect and replay request history",
	}

	var limit int
	list := &cobra.Command{
		Use:   "list",
		Short: "List recent requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := deps.historyService.List(cmd.Context(), limit)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s  %s  %d  %s\n", entry.ID, entry.Method, entry.URL, entry.StatusCode, entry.Duration)
			}
			return nil
		},
	}
	list.Flags().IntVarP(&limit, "limit", "n", 20, "number of history entries")

	show := &cobra.Command{
		Use:   "show ID",
		Short: "Show one history entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := deps.historyService.Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), formatHistoryEntry(entry))
			return nil
		},
	}

	replay := &cobra.Command{
		Use:   "replay ID",
		Short: "Replay one request from history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := deps.historyService.Replay(cmd.Context(), args[0], 30*time.Second)
			if err != nil {
				return err
			}
			output, err := respformatter.NewResponseFormatter(true).Format(response)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		},
	}

	cmd.AddCommand(list, show, replay)
	return cmd
}

func formatHistoryEntry(entry models.HistoryEntry) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "ID: %s\n", entry.ID)
	fmt.Fprintf(&builder, "Method: %s\n", entry.Method)
	fmt.Fprintf(&builder, "URL: %s\n", entry.URL)
	fmt.Fprintf(&builder, "Status: %d\n", entry.StatusCode)
	fmt.Fprintf(&builder, "Duration: %s\n", entry.Duration)
	fmt.Fprintf(&builder, "Executed At: %s\n", entry.ExecutedAt.Format(time.RFC3339))
	if entry.Error != "" {
		fmt.Fprintf(&builder, "Error: %s\n", entry.Error)
	}
	builder.WriteString("Headers:\n")
	for key, value := range entry.Headers {
		fmt.Fprintf(&builder, "  %s: %s\n", key, value)
	}
	if len(entry.Body) > 0 {
		fmt.Fprintf(&builder, "Body:\n%s\n", entry.Body)
	}
	if len(entry.Response) > 0 {
		fmt.Fprintf(&builder, "Response:\n%s\n", entry.Response)
	}
	return builder.String()
}
