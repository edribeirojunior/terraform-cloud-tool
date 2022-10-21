package cmd

import (
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/client"
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/variables"
	workspace "github.com/edribeirojunior/terraform-cloud-tool/pkg/workspaces"
)

var (
	token, org, wtags, wtype, varName, varValue, setTags, setTerraformVersion string
	varSensitive, setLock                                                     bool
)

func NewVariableClient(token, org, wtags, wtype, varName, varValue string, varSensitive bool) variables.Variable {

	cl := client.NewTfClient(token)

	o := client.GetOrg(cl, org)

	wsList, _ := client.GetWorkspacesList(cl, o, wtags, wtype, setLock)

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
	wsList, wsListLocked := client.GetWorkspacesList(cl, o, wtags, wtype, setLock)

	//wsRunsList := client.GetWorkspacesRunsList(cl, o, wsList)
	ws := workspace.Workspace{
		Cl:               cl,
		List:             wsList,
		ListLocked:       wsListLocked,
		Tags:             &wtags,
		TerraformVersion: &setTerraformVersion,
		//RunsList:         wsRunsList,
	}

	return ws
}
