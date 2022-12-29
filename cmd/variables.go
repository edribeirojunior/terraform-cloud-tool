package cmd

import (
	"fmt"

	"github.com/edribeirojunior/terraform-cloud-tool/pkg/variables"
	"github.com/spf13/cobra"
)

func init() {
	varsCmd.PersistentFlags().StringVar(&wtags, "wtg", "", "The tags to filter the workspaces")
	varsCmd.PersistentFlags().StringVar(&varName, "vn", "", "Variable Name")
	varsCmd.PersistentFlags().StringVar(&varValue, "vv", "", "Variable Value")
	varsCmd.PersistentFlags().BoolVar(&varSensitive, "vs", false, "Variable Value is Sensitive")

	varsCmd.AddCommand(varsReadCmd)
	varsReadCmd.Flags().Bool("show-all", false, "if set shows variables for all rings instead of only for the first ring")

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
	RunE: func(cmd *cobra.Command, args []string) error {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		showAll, err := cmd.Flags().GetBool("show-all")
		if err != nil {
			return fmt.Errorf("invalid flag: %s", err)
		}

		if showAll {
			variables.ReadAll(nCl)
		} else {
			variables.Read(nCl)
		}

		return nil
	},
}

var varsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables in a Workspace",
	Long:  "List variables in a Workspace",
	Run: func(cmd *cobra.Command, args []string) {

		nCl := NewClient(token, org, wtags, wtype, varName, varValue, setTags, varSensitive)

		variables.List(nCl)
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
			variables.Apply(nCl)
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
