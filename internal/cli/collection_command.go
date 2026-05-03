package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCollectionCommand(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Manage request collections",
	}

	create := &cobra.Command{
		Use:   "create NAME",
		Short: "Create a collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deps.collectionService.Create(cmd.Context(), args[0])
		},
	}

	addOpts := requestOptions{timeoutSec: 30, pretty: true}
	add := &cobra.Command{
		Use:   "add COLLECTION REQUEST_NAME METHOD URL",
		Short: "Add a saved request to a collection",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			req, err := buildRequest(args[2], args[3], addOpts, deps.reader)
			if err != nil {
				return err
			}
			return deps.collectionService.AddRequest(cmd.Context(), args[0], args[1], req)
		},
	}
	add.Flags().StringArrayVarP(&addOpts.headers, "header", "H", nil, "custom header in 'Key: Value' format")
	add.Flags().StringVarP(&addOpts.body, "body", "b", "", "raw request body")
	add.Flags().StringVarP(&addOpts.bodyFile, "body-file", "f", "", "load request body from file")

	list := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			collections, err := deps.collectionService.List(cmd.Context())
			if err != nil {
				return err
			}
			for _, collection := range collections {
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d requests)\n", collection.Name, len(collection.Requests))
			}
			return nil
		},
	}

	show := &cobra.Command{
		Use:   "show NAME",
		Short: "Show collection requests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collection, err := deps.collectionService.Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", collection.Name)
			for _, req := range collection.Requests {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  %s\n", req.Name, req.Method, req.URL)
			}
			return nil
		},
	}

	cmd.AddCommand(create, add, list, show)
	return cmd
}
