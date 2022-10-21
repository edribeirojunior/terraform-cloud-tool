package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/hashicorp/go-tfe"
)

type Execute interface {
	Create()
	Delete()
	Read()
	ApproveChanges() string
}

type Workspaces struct {
	Cl                 *tfe.Client
	List               []tfe.Workspace
	ListLocked         []tfe.Workspace
	Tags               *string `json:",omitempty"`
	Version            *string `json:",omitempty"`
	Variables          *string `json:",omitempty"`
	VariablesValue     *string `json:",omitempty"`
	VariablesSensitive *bool   `json:",omitempty"`

	//Workspace ws.Workspace `json:",omitempty"`
	//Variable  vr.Variable  `json:",omitempty"`
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

func NewWorkspace(client *tfe.Client, wsList []*tfe.Workspace, workspaceTypePtr, setTags, setVersion, variableNamePtr, variableValuePtr *string, variableSenstivePtr *bool, lock bool) Workspaces {

	wList, wsListLocked := GetWorkspace(wsList, *workspaceTypePtr)

	return Workspaces{
		Cl:                 client,
		List:               wList,
		ListLocked:         wsListLocked,
		Tags:               setTags,
		Version:            setVersion,
		Variables:          variableNamePtr,
		VariablesValue:     variableValuePtr,
		VariablesSensitive: variableSenstivePtr,
	}

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

func GetWorkspacesList(client *tfe.Client, org *tfe.Organization, wtags, wsfilter string, lock bool) ([]tfe.Workspace, []tfe.Workspace) {
	ctx := context.Background()

	var workspaces []*tfe.Workspace
	workspace, err := client.Workspaces.List(ctx, org.Name, tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
		Tags:        &wtags,
	})

	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	for i := 0; i < workspace.TotalPages; i++ {
		workspaces = append(workspaces, workspace.Items...)
		if workspace.CurrentPage == workspace.TotalPages {
			break
		}
	}

	wsList, wsListLocked := GetWorkspace(workspaces, wsfilter)

	return wsList, wsListLocked

}

func GetWorkspace(workspaces []*tfe.Workspace, wsfilter string) ([]tfe.Workspace, []tfe.Workspace) {
	var wsList []tfe.Workspace
	var wsListLocked []tfe.Workspace

	fmt.Printf("Getting %d workspaces\n", len(workspaces))

	for _, w := range workspaces {

		found, err := regexp.MatchString(wsfilter, w.Name)

		if err != nil {
			fmt.Println(err)
		}

		if found {

			if !w.Locked {
				wsList = append(wsList, *w)
			} else {
				fmt.Printf("Post: Workspace Name %s, State %t\n", w.Name, w.Locked)
				wsListLocked = append(wsListLocked, *w)
			}

		}
	}

	return wsList, wsListLocked
}

func GetWorkspacesRunsList(client *tfe.Client, org *tfe.Organization, wsList []tfe.Workspace) []tfe.RunList {
	ctx := context.Background()

	var wsRunsList []tfe.RunList

	for i := 0; i < len(wsList); i++ {
		runList, err := client.Runs.List(ctx, wsList[i].ID, tfe.RunListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})

		if err != nil {
			log.Fatal(err)
			return nil
		}

		wsRunsList = append(wsRunsList, *runList)

	}

	return wsRunsList

}
