package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"reqium/internal/app"
	"reqium/internal/infrastructure/formatter"
	"reqium/internal/infrastructure/http"
	"reqium/internal/interfaces"
)

func newMethodCommand(methodName string, method string, client *http.NetHTTPClient, reader interfaces.FileReader) *cobra.Command {
	opts := requestOptions{
		timeoutSec: 30,
		pretty:     true,
	}

	cmd := &cobra.Command{
		Use:   methodName + " URL",
		Short: "Send a " + method + " request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := buildRequest(method, args[0], opts, reader)
			if err != nil {
				return err
			}

			service := app.NewRequestService(client, formatter.NewResponseFormatter(opts.pretty))
			output, err := service.Send(cmd.Context(), req)
			if err != nil {
				return err
			}

			fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&opts.headers, "header", "H", nil, "custom header in 'Key: Value' format")
	cmd.Flags().StringVarP(&opts.body, "body", "b", "", "raw request body")
	cmd.Flags().StringVarP(&opts.bodyFile, "body-file", "f", "", "load request body from file")
	cmd.Flags().IntVarP(&opts.timeoutSec, "timeout", "t", 30, "request timeout in seconds")
	cmd.Flags().BoolVar(&opts.pretty, "pretty", true, "pretty-print JSON responses")

	return cmd
}
