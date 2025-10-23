package cmd

import (
	"fmt"

	"github.com/fopina/proxyone/pkg/config"
	"github.com/fopina/proxyone/pkg/proxy"
	"github.com/spf13/cobra"
)

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxyone",
		Short: "A simple HTTP proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			p := proxy.NewProxy(cfg)
			return p.Start()
		},
	}

	cmd.AddCommand(newVersionCmd(version)) // version subcommand

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}
