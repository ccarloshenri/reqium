package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"reqium/internal/app"
	respformatter "reqium/internal/implementations/formatter"
)

func newMethodCommand(methodName string, method string, deps dependencies) *cobra.Command {
	opts := requestOptions{
		timeoutSec: 30,
		pretty:     true,
	}

	cmd := &cobra.Command{
		Use:   methodName + " URL",
		Short: "Send a " + method + " request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := buildRequest(method, args[0], opts, deps.reader)
			if err != nil {
				return err
			}

			variables, err := deps.envService.Variables(cmd.Context(), opts.env)
			if err != nil {
				return err
			}
			req, err = deps.resolver.ResolveRequest(req, variables)
			if err != nil {
				return err
			}

			service := app.NewRequestServiceWithHistory(deps.client, respformatter.NewResponseFormatter(opts.pretty), deps.historyRepo)
			output, err := service.SendFormatted(cmd.Context(), req)
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
	cmd.Flags().StringVar(&opts.env, "env", "", "environment to resolve {{variables}}")

	return cmd
}
