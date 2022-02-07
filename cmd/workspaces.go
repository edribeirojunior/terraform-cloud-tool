package cmd

import (
	"fmt"

	workspace "github.com/edribeirojunior/terraform-cloud-tool/pkg/workspaces"
	"github.com/spf13/cobra"
)

func init() {
	worksCmd.PersistentFlags().StringVar(&setTags, "ts", "", "The tags to set in the workspace")

	worksCmd.AddCommand(worksApplyCmd)
	worksCmd.AddCommand(worksDeleteCmd)
}

var worksCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace function to Terraform Cloud",
	Long:  `Create/Delete/Edit Workspaces in Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var worksApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply to workspaces",
	Long:  `Apply Workspaces in Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		approved := workspace.ApproveChanges(nCl, "create")

		if approved == "y" || approved == "yes" {
			workspace.Create(nCl)
		}
	},
}

var worksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete to workspaces",
	Long:  `Delete Workspaces in Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete")
		cmd.Help()
	},
}
