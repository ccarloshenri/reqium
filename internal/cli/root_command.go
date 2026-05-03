package cli

import (
	"github.com/spf13/cobra"

	"reqium/internal/infrastructure/filesystem"
	httpinfra "reqium/internal/infrastructure/http"
)

func NewRootCommand() *cobra.Command {
	reader := filesystem.NewOSFileReader()
	client := httpinfra.NewNetHTTPClient()

	root := &cobra.Command{
		Use:          "reqium",
		Short:        "A fast, minimal terminal API client",
		SilenceUsage: true,
	}

	root.AddCommand(newMethodCommand("get", "GET", client, reader))
	root.AddCommand(newMethodCommand("post", "POST", client, reader))
	root.AddCommand(newMethodCommand("put", "PUT", client, reader))
	root.AddCommand(newMethodCommand("patch", "PATCH", client, reader))
	root.AddCommand(newMethodCommand("delete", "DELETE", client, reader))

	return root
}
