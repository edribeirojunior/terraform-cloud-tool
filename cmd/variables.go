package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	varsCmd.PersistentFlags().StringVar(&wtags, "wtg", "", "The tags to filter the workspaces")
	varsCmd.PersistentFlags().StringVar(&varName, "vn", "", "Variable Name")
	varsCmd.PersistentFlags().StringVar(&varValue, "vv", "", "Variable Value")
	varsCmd.PersistentFlags().BoolVar(&varSensitive, "vs", false, "Variable Value is Sensitive")

	varsCmd.AddCommand(varsReadCmd)
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

var varsReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read variable in a Workspace",
	Long:  "Read variable in a Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		varCl := NewVariableClient(token, org, wtags, wtype, varName, varValue, varSensitive)

		varCl.Read()
	},
}

var varsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables in a Workspace",
	Long:  "List variables in a Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		varCl := NewVariableClient(token, org, wtags, wtype, varName, varValue, varSensitive)

		varCl.ListVars()
	},
}

var varsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Variable function to Terraform Cloud",
	Long:  `Create Variables from Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		varCl := NewVariableClient(token, org, wtags, wtype, varName, varValue, varSensitive)

		approved := varCl.ApproveChanges("create")

		if approved == "y" || approved == "yes" {
			varCl.Apply()
		}
	},
}

var varsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Variable function to Terraform Cloud",
	Long:  `Delete Variables from Terraform Cloud`,
	Run: func(cmd *cobra.Command, args []string) {

		varCl := NewVariableClient(token, org, wtags, wtype, varName, varValue, varSensitive)

		approved := varCl.ApproveChanges("delete")

		if approved == "y" || approved == "yes" {
			varCl.Delete()
		}

	},
}
