package cobrax

import "github.com/spf13/cobra"

func disableFlagsSorting(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	for _, sub := range cmd.Commands() {
		disableFlagsSorting(sub)
	}
}

func hideHelpCommand(cmd *cobra.Command) {
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func hideHelpFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("help", false, "")
	_ = cmd.PersistentFlags().MarkHidden("help")
}

func shortFlag(flag string) string {
	for _, v := range flag {
		return string(v)
	}

	return ""
}
