package cobrax

import "github.com/spf13/cobra"

type Flag interface {
	Apply(*cobra.Command)
}

type StringFlag struct {
	Name     string
	Usage    string
	Required bool
}

func (f *StringFlag) Apply(cmd *cobra.Command) {
	cmd.Flags().StringP(f.Name, shortFlag(f.Name), "", f.Usage)

	if f.Required {
		_ = cmd.MarkFlagRequired(f.Name)
	}
}
