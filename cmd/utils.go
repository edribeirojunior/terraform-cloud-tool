package cmd

import (
	"github.com/edribeirojunior/terraform-cloud-tool/pkg/client"
)

var (
	token, org, wtags, wtype, varName, varValue, setTags, setVersion string
	varSensitive                                                     bool
)

func NewClient(token, org, wtags, wtype, varName, varvalue, setTags string, varSensitive bool) client.Workspaces {

	cl := client.NewTfClient(token)

	o := client.GetOrg(cl, org)

	wsList := client.GetWorkspacesList(cl, o, wtags)

	ws := client.NewWorkspace(cl, wsList, &wtype, &setTags, &setVersion, &varName, &varValue, &varSensitive)
	return ws

}
