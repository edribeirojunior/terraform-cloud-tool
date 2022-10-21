package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	worksCmd.PersistentFlags().StringVar(&setTags, "ts", "", "The tags to set in the workspace")
	worksCmd.PersistentFlags().StringVar(&setTerraformVersion, "tfv", "1.0.4", "The terraform version to set in the workspace")
	worksCmd.PersistentFlags().StringVar(&wtags, "wtg", "", "The tags to filter the workspaces")
	worksCmd.PersistentFlags().BoolVar(&setLock, "tl", true, "Lock Workspace")
	worksCmd.AddCommand(worksApplyCmd)
	worksCmd.AddCommand(worksDeleteCmd)
	worksCmd.AddCommand(worksReadCmd)
	worksCmd.AddCommand(worksRunsDelete)
	worksCmd.AddCommand(worksUnlock)
}

var worksCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace function to Terraform Cloud",
	Long:  `Create/Delete/Edit Workspaces in Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var worksRunsDelete = &cobra.Command{
	Use:   "runs-cancel",
	Short: "runs cancel all",
	Long:  `Cancel all current Runs`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewWorkspaceClient(token, org, wtags, wtype, setTags, setTerraformVersion)

		approved := nCl.ApproveChanges("cancel")

		if approved == "y" || approved == "yes" {
			nCl.Delete()
		}
	},
}

var worksUnlock = &cobra.Command{
	Use:   "unlock",
	Short: "unlock",
	Long:  `Unlock Workspace`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewWorkspaceClient(token, org, wtags, wtype, setTags, setTerraformVersion)

		approved := nCl.ApproveChanges("unlock")

		if approved == "y" || approved == "yes" {
			nCl.Unlock()
		}
	},
}

var worksApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply to workspaces",
	Long:  `Apply Workspaces in Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewWorkspaceClient(token, org, wtags, wtype, setTags, setTerraformVersion)

		approved := nCl.ApproveChanges("create")

		if approved == "y" || approved == "yes" {
			nCl.Create()
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

var worksReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read Workspace",
	Long:  "Read Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewWorkspaceClient(token, org, wtags, wtype, setTags, setTerraformVersion)

		nCl.Read(org)
	},
}
