package cmd

import (
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/variables"
	"github.com/spf13/cobra"
)

func init() {
	varsCmd.PersistentFlags().StringVar(&wtags, "wtg", "", "The tags to filter the workspaces")
	varsCmd.PersistentFlags().StringVar(&wtype, "wt", "", "Filter the Workspace Name (REGEX)")
	varsCmd.PersistentFlags().StringVar(&varName, "vn", "", "Variable Name")
	varsCmd.PersistentFlags().StringVar(&varValue, "vv", "", "Variable Value")
	varsCmd.PersistentFlags().BoolVar(&varSensitive, "vs", false, "Variable Value is Sensitive")

	varsCmd.AddCommand(varsListCmd)
	varsCmd.AddCommand(varsApplyCmd)
	varsCmd.AddCommand(varsDeleteCmd)
}

var varsCmd = &cobra.Command{
	Use:   "variable",
	Short: "Variable function to Terraform Cloud",
	Long:  `Create/Delete/Edit Variables from Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var varsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables in a Workspace",
	Long:  "List variables in a Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		variables.Read(nCl)
	},
}

var varsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Variable function to Terraform Cloud",
	Long:  `Create Variables from Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		approved := variables.ApproveChanges(nCl, "create")

		if approved == "y" || approved == "yes" {
			variables.Create(nCl)
		}
	},
}

var varsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Variable function to Terraform Cloud",
	Long:  `Delete Variables from Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		approved := variables.ApproveChanges(nCl, "delete")

		if approved == "y" || approved == "yes" {
			variables.Delete(nCl)
		}

	},
}
