package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newEnvironmentCommand(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "env",
		Aliases: []string{"environment"},
		Short:   "Manage environments and variables",
	}

	create := &cobra.Command{
		Use:   "create NAME",
		Short: "Create an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deps.envService.Create(cmd.Context(), args[0])
		},
	}

	set := &cobra.Command{
		Use:   "set ENV KEY VALUE",
		Short: "Set an environment variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deps.envService.Set(cmd.Context(), args[0], args[1], args[2])
		},
	}

	use := &cobra.Command{
		Use:   "use NAME",
		Short: "Set the active environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deps.envService.Use(cmd.Context(), args[0])
		},
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List environments",
		RunE: func(cmd *cobra.Command, args []string) error {
			envs, err := deps.envService.List(cmd.Context())
			if err != nil {
				return err
			}
			for _, env := range envs {
				active := " "
				if env.Active {
					active = "*"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", active, env.Name)
				keys := make([]string, 0, len(env.Variables))
				for key := range env.Variables {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				for _, key := range keys {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", key, env.Variables[key])
				}
			}
			return nil
		},
	}

	cmd.AddCommand(create, set, use, list)
	return cmd
}
