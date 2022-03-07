package workspace

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

func NewClient(client *tfe.Client, ws []tfe.Workspace, setTags, setTerraformVersion string) Workspace {

	return Workspace{
		Cl:               client,
		List:             ws,
		Tags:             &setTags,
		TerraformVersion: &setTerraformVersion,
	}
}

func (ws *Workspace) Read(org string) {

	for i := 0; i < len(ws.List); i++ {
		fmt.Printf("Workspaces Name %s \n", ws.List[i].Name)

	}
}

func (ws *Workspace) Create() {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		if *ws.Tags != "" {
			tag := tfe.Tag{Name: *ws.Tags}
			tags := []*tfe.Tag{&tag}
			err := ws.Cl.Workspaces.AddTags(ctx, ws.List[i].ID, tfe.WorkspaceAddTagsOptions{
				Tags: tags,
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Tags created: %s, in Workspace %s\n", *ws.Tags, ws.List[i].Name)
		}

		if *ws.TerraformVersion != "" {

			tfVersion := *ws.TerraformVersion
			_, err := ws.Cl.Workspaces.UpdateByID(ctx, ws.List[i].ID, tfe.WorkspaceUpdateOptions{
				TerraformVersion: tfe.String(tfVersion),
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Terraform Version updated to %s, in Workspace %s\n", *ws.TerraformVersion, ws.List[i].Name)
		}
	}
}

// func Update(ws client.Workspaces, org string) {
// 	ctx := context.Background()

// 	for i := 0; i < len(ws.List); i++ {
// 		if *ws.TerraformVersion != "" {

// 			tfVersion := *ws.TerraformVersion
// 			_, err := ws.Cl.Workspaces.Update(ctx, org, ws.List[i].Name, tfe.WorkspaceUpdateOptions{
// 				TerraformVersion: tfe.String(tfVersion),
// 			})

// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			fmt.Printf("Terraform Version updated to %s, in Workspace %s\n", *ws.TerraformVersion, ws.List[i].Name)
// 		}
// 	}
// }

func (ws *Workspace) ApproveChanges(action string) string {

	var wsList []string

	for i := 0; i < len(ws.List); i++ {
		wsList = append(wsList, ws.List[i].Name)
	}
	if action == "create" {
		fmt.Printf("The Tags/Terraform Version will be created in these workspaces:\n")
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
		// } else if action == "delete" {
		// 	fmt.Printf("The Variable %s will be deleted in these workspaces:\n", *ws.Variables)
		// 	for _, i := range wsList {
		// 		fmt.Printf("%s\n", i)
		// 	}
	} else {
		fmt.Println("operation not permitted, exiting...")
		os.Exit(1)

	}

	fmt.Printf("\nApprove changes? (y/n)\n")
	var approved string
	fmt.Scanln(&approved)

	return approved

}
