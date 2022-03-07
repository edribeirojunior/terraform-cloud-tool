package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terraform-cloud-tool",
	Short: "Terraform Cloud Tool is a tool to manage Terraform Cloud",
	Long:  `Terraform Cloud Tool is a tool to manage Terraform Cloud.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVar(&token, "t", "", "The token to use to authenticate in TFCloud")
	rootCmd.PersistentFlags().StringVar(&org, "o", "", "The organization to use to authenticate in TFCloud")
	rootCmd.PersistentFlags().StringVar(&wtype, "wt", "", "Filter the Workspace Name (REGEX)")
	rootCmd.AddCommand(worksCmd)
	rootCmd.AddCommand(varsCmd)

}
