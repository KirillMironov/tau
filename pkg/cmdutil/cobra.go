package cmdutil

import "github.com/spf13/cobra"

func DisableFlagsSorting(command *cobra.Command) {
	command.Flags().SortFlags = false
	command.PersistentFlags().SortFlags = false

	for _, cmd := range command.Commands() {
		DisableFlagsSorting(cmd)
	}
}

func HideHelpCommand(command *cobra.Command) {
	command.SetHelpCommand(&cobra.Command{Hidden: true})
}

func HideHelpFlags(command *cobra.Command) {
	command.PersistentFlags().Bool("help", false, "")
	_ = command.PersistentFlags().MarkHidden("help")
}

func ShortFlag(flag string) string {
	for _, v := range flag {
		return string(v)
	}

	return ""
}
