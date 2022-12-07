package cobrax

import "github.com/spf13/cobra"

type Command struct {
	Usage       string
	Description string
	Example     string
	Args        cobra.PositionalArgs
	Flags       []Flag
	Action      func(cmd *cobra.Command, args []string) error
	Subcommands []*Command
	Options     CommandOptions
}

type CommandOptions struct {
	SilenceErrors       bool
	SilenceUsage        bool
	DisableDefaultCmd   bool
	DisableFlagsSorting bool
	HideHelpCommand     bool
	HideHelpFlags       bool
}

func (c *Command) Execute() error {
	cmd := c.toCobra()

	cmd.SilenceErrors = c.Options.SilenceErrors
	cmd.SilenceUsage = c.Options.SilenceUsage
	cmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: c.Options.DisableDefaultCmd,
	}

	if c.Options.DisableFlagsSorting {
		disableFlagsSorting(cmd)
	}
	if c.Options.HideHelpCommand {
		hideHelpCommand(cmd)
	}
	if c.Options.HideHelpFlags {
		hideHelpFlags(cmd)
	}

	return cmd.Execute()
}

func (c *Command) toCobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:     c.Usage,
		Long:    c.Description,
		Example: c.Example,
		Args:    c.Args,
		RunE:    c.Action,
	}

	for _, flag := range c.Flags {
		flag.Apply(cmd)
	}

	for _, sub := range c.Subcommands {
		cmd.AddCommand(sub.toCobra())
	}

	return cmd
}
