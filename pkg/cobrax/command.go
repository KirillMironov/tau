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
}

func (c *Command) Build() *cobra.Command {
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
		cmd.AddCommand(sub.Build())
	}

	return cmd
}
