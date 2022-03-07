package cmd

import (
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/client"
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/variables"
	workspace "github.com/edribeirojunior/terraform-cloud-tool/pkg/workspaces"
)

var (
	token, org, wtags, wtype, varName, varValue, setTags, setTerraformVersion string
	varSensitive                                                              bool
)

func NewVariableClient(token, org, wtags, wtype, varName, varValue string, varSensitive bool) variables.Variable {

	cl := client.NewTfClient(token)

	o := client.GetOrg(cl, org)

	wsList := client.GetWorkspacesList(cl, o, wtags, wtype)

	vS := variables.Variable{
		Cl:                 cl,
		List:               wsList,
		Variables:          &varName,
		VariablesValue:     &varValue,
		VariablesSensitive: &varSensitive,
	}

	return vS
}

func NewWorkspaceClient(token, org, wtags, wtype, setTags, setTerraformVersion string) workspace.Workspace {
	cl := client.NewTfClient(token)
	o := client.GetOrg(cl, org)
	wsList := client.GetWorkspacesList(cl, o, wtags, wtype)
	ws := workspace.Workspace{
		Cl:               cl,
		List:             wsList,
		Tags:             &wtags,
		TerraformVersion: &setTerraformVersion,
	}

	return ws
}
