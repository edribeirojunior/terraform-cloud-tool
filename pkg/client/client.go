package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	vr "github.com/edribeirojunior/terraform-cloud-tool/pkg/variables"
	ws "github.com/edribeirojunior/terraform-cloud-tool/pkg/workspaces"
	"github.com/hashicorp/go-tfe"
)

type Execute interface {
	Create()
	Delete()
	Read()
	ApproveChanges() string
}

type Workspaces struct {
	Workspace ws.Workspace `json:",omitempty"`
	Variable  vr.Variable  `json:",omitempty"`
}

func NewTfClient(token string) *tfe.Client {

	config := &tfe.Config{}
	if token != "" {
		config.Token = token
	} else {

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		tftokenpath := home + "/.terraform.d/credentials.tfrc.json"
		config.Token = readTfToken(tftokenpath)
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return client

}

func readTfToken(tPath string) string {

	jsonFile, err := os.Open(tPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var token map[string]interface{}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &token)

	cred := token["credentials"].(map[string]interface{})
	tftoken := cred["app.terraform.io"].(map[string]interface{})

	return tftoken["token"].(string)

}

func GetOrg(client *tfe.Client, organization string) *tfe.Organization {
	ctx := context.Background()

	org, err := client.Organizations.Read(ctx, organization)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return org

}

func GetWorkspacesList(client *tfe.Client, org *tfe.Organization, wtags, wsfilter string) []tfe.Workspace {
	ctx := context.Background()

	workspace, err := client.Workspaces.List(ctx, org.Name, tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
		Tags:        &wtags,
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	wsList := GetWorkspace(workspace, wsfilter)

	return wsList

}

func GetWorkspace(workspaceName *tfe.WorkspaceList, wsfilter string) []tfe.Workspace {
	var wsList []tfe.Workspace
	for i := 0; i < len(workspaceName.Items); i++ {

		found, _ := regexp.MatchString(wsfilter, workspaceName.Items[i].Name)
		if found {
			wsList = append(wsList, *workspaceName.Items[i])
		}

	}

	return wsList
}
