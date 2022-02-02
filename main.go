package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	tfe "github.com/hashicorp/go-tfe"
)

const (
	terraformTokenPath = "/Users/ero/.terraform.d/credentials.tfrc.json"
)

type Execute interface {
	Create()
	Delete()
	Read()
	ApproveChanges() string
}

type Workspaces struct {
	cl                 *tfe.Client
	List               []tfe.Workspace
	Tags               *string `json:",omitempty"`
	Variables          *string `json:",omitempty"`
	VariablesValue     *string `json:",omitempty"`
	VariablesSensitive *bool   `json:",omitempty"`
	variableDelete     *bool   `json:",omitempty"`
}

func main() {

	tokenPtr := flag.String("t", "", "The token to use to authenticate in TFCloud")
	organizationPtr := flag.String("o", "organization-name", "The organization to use to authenticate in TFCloud")
	workspaceTagsPtr := flag.String("wtags", "", "The tags to filter the workspaces")
	workspaceTypePtr := flag.String("wtype", "", "Filter the workspace (Regex)")
	setTags := flag.String("ts", "", "The tags to set in the workspace")
	variableNamePtr := flag.String("vn", "", "The variable name to set")
	variableValuePtr := flag.String("vv", "", "The variable value to set")
	variableSenstivePtr := flag.Bool("vs", false, "The variable is sensitive")
	variableDelete := flag.Bool("vd", false, "The variable should be deleted")
	variableList := flag.Bool("vl", false, "Reading Variables")

	flag.Parse()

	config := &tfe.Config{}
	if *tokenPtr != "" {
		config.Token = *tokenPtr
	} else {
		config.Token = readTfToken()
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get the organization.
	org, err := client.Organizations.Read(ctx, *organizationPtr)
	if err != nil {
		log.Fatal(err)
	}

	// Get the workspace.
	tags := *workspaceTagsPtr
	workspace, err := client.Workspaces.List(ctx, org.Name, tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
		Tags:        &tags})
	if err != nil {
		log.Fatal(err)
	}

	ws := NewWorkspace(client, workspace, workspaceTypePtr, setTags, variableNamePtr, variableValuePtr, variableSenstivePtr, variableDelete)

	if *variableList {
		Read(ws)
	} else {

		approved := ApproveChanges(ws)

		if approved != "y" {
			fmt.Println("Changes not Approved.")
			os.Exit(1)
		}

		Middleware(ws, *variableDelete, *variableNamePtr, *variableValuePtr, *setTags)

		fmt.Println("Finished.")
	}
}

func readTfToken() string {
	jsonFile, err := os.Open(terraformTokenPath)
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

func NewWorkspace(client *tfe.Client, wsList *tfe.WorkspaceList, workspaceTypePtr, setTags, variableNamePtr, variableValuePtr *string, variableSenstivePtr, variableDelete *bool) Workspaces {
	wList := getWorkspace(wsList, *workspaceTypePtr)

	return Workspaces{
		cl:                 client,
		List:               wList,
		Tags:               setTags,
		Variables:          variableNamePtr,
		VariablesValue:     variableValuePtr,
		VariablesSensitive: variableSenstivePtr,
		variableDelete:     variableDelete,
	}

}

func getWorkspace(workspaceName *tfe.WorkspaceList, wsfilter string) []tfe.Workspace {
	var wsList []tfe.Workspace
	for i := 0; i < len(workspaceName.Items); i++ {

		//fmt.Println(workspaceName.Items[i].Name)
		found, _ := regexp.MatchString(wsfilter, workspaceName.Items[i].Name)
		if found {
			wsList = append(wsList, *workspaceName.Items[i])
		}

	}

	return wsList
}

func Create(ws Execute) {
	ws.Create()
}

func Delete(ws Execute) {
	ws.Delete()
}

func Read(ws Execute) {
	ws.Read()
}

func ApproveChanges(ws Execute) string {
	return ws.ApproveChanges()
}

func Middleware(ws Execute, wsDelete bool, wsVar, wsVal, wsTags string) {

	if !wsDelete && wsVar != "" && wsVal != "" || wsTags != "" {
		ws.Create()
	} else if wsDelete && wsVar != "" {
		ws.Delete()
	} else {
		fmt.Println("Variable value or name is null")
		os.Exit(1)
	}

}

func (ws Workspaces) Delete() {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		variable, err := ws.cl.Variables.List(ctx, ws.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})
		if err != nil {
			fmt.Println("Failed to list the variables")
			log.Fatal(err)
		}

		for j := 0; j < len(variable.Items); j++ {

			if variable.Items[j].Key == *ws.Variables {
				err = ws.cl.Variables.Delete(ctx, ws.List[i].ID, variable.Items[j].ID)
				if err != nil {
					fmt.Println("Failed to delete variable")
					log.Fatal(err)
				}

				fmt.Printf("Variable deleted: %s, in Workspace %s\n", variable.Items[j].Key, ws.List[i].Name)
			}
		}
	}
}

func (ws Workspaces) Create() {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		if *ws.Variables != "" && !*ws.variableDelete {
			variable, err := ws.cl.Variables.Create(ctx, ws.List[i].ID, tfe.VariableCreateOptions{
				Key:       ws.Variables,
				Value:     ws.VariablesValue,
				Sensitive: ws.VariablesSensitive,
				Category:  tfe.Category("terraform"),
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Variable created: %s, in Workspace %s\n", variable.Key, ws.List[i].Name)

		}

		if *ws.Tags != "" {
			tag := tfe.Tag{Name: *ws.Tags}
			tags := []*tfe.Tag{&tag}
			err := ws.cl.Workspaces.AddTags(ctx, ws.List[i].ID, tfe.WorkspaceAddTagsOptions{
				Tags: tags,
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Tags created: %s, in Workspace %s\n", *ws.Tags, ws.List[i].Name)
		}
	}
}

func (ws Workspaces) Read() {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		fmt.Printf("Listing all Variables in %s\n", ws.List[i].Name)
		variable, _ := ws.cl.Variables.List(ctx, ws.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})

		for j := 0; j < len(variable.Items); j++ {
			if !variable.Items[j].Sensitive {
				fmt.Printf("Name: %s, Value: %s\n", variable.Items[j].Key, variable.Items[j].Value)
				// for _, j := range wsList {
				// 	fmt.Printf("Name: %s\n", j.Key)
				// }
			}
		}
	}
}

func (ws Workspaces) ApproveChanges() string {

	var wsList []string

	for i := 0; i < len(ws.List); i++ {
		wsList = append(wsList, ws.List[i].Name)
	}

	if *ws.Variables != "" && !*ws.variableDelete {
		fmt.Printf("The Variable %s will be created in these workspaces:\n", *ws.Variables)
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
	}

	if *ws.Variables != "" && *ws.variableDelete {
		fmt.Printf("The Variable %s will be deleted in these workspaces:\n", *ws.Variables)
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
	}

	if *ws.Tags != "" {
		fmt.Printf("The Variable %s will be created in these workspaces:\n", *ws.Variables)
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
	}

	if *ws.Variables == "" && *ws.Tags == "" {

		return "Nothing to do here..."
	}

	fmt.Printf("\nApprove changes? (y/n)\n")
	var approved string
	// Taking input from user
	fmt.Scanln(&approved)

	return approved
}
