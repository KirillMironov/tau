package cobrax

import "github.com/spf13/cobra"

type Flag interface {
	Apply(*cobra.Command)
}

type StringFlag struct {
	Name     string
	Alias    string
	Usage    string
	Required bool
}

func (f *StringFlag) Apply(command *cobra.Command) {
	command.Flags().StringP(f.Name, f.Alias, "", f.Usage)

	if f.Required {
		_ = command.MarkFlagRequired(f.Name)
	}
}
