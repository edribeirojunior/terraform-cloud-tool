package workspace

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

func NewClient(client *tfe.Client, ws []tfe.Workspace, wsLocked []tfe.Workspace, setTags, setTerraformVersion string) Workspace {

	return Workspace{
		Cl:               client,
		List:             ws,
		ListLocked:       wsLocked,
		Tags:             &setTags,
		TerraformVersion: &setTerraformVersion,
	}
}

func (ws *Workspace) Read(org string) {

	for i := 0; i < len(ws.List); i++ {
		fmt.Printf("Workspaces Name %s \n", ws.List[i].Name)

	}
}

func (ws *Workspace) Delete() {
	ctx := context.Background()

	if len(ws.RunsList) > 0 {

		for j := 0; j < len(ws.RunsList); j++ {

			for h := 0; j < len(ws.RunsList[j].Items); h++ {
				discardRun := "Discarded by CLI"
				err := ws.Cl.Runs.Cancel(ctx, ws.RunsList[j].Items[h].ID, tfe.RunCancelOptions{
					Comment: &discardRun,
				})

				if err != nil {
					log.Fatal(err)
				}
			}
		}

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

func (ws *Workspace) Unlock() {
	ctx := context.Background()

	for _, w := range ws.ListLocked {

		fmt.Println("Unlocking Workspace", w.Name, w.Tags)

		_, err := ws.Cl.Workspaces.ForceUnlock(ctx, w.ID)

		if err != nil {
			fmt.Println(err)
		}

	}

}

func (ws *Workspace) ApproveChanges(action string) string {

	if action == "create" {
		fmt.Printf("The Tags/Terraform Version will be created in these workspaces:\n")
		for _, i := range ws.List {
			fmt.Printf("%s\n", i.Name)
		}

	} else if action == "cancel" {
		fmt.Printf("The runs for the following workspaces, will be cancelled:\n")
		for _, i := range ws.List {
			fmt.Printf("%s\n", i.Name)
		}

	} else if action == "unlock" {
		fmt.Printf("Unlock workspaces:\n")
		for _, i := range ws.ListLocked {
			fmt.Printf("%s\n", i.Name)
		}

	} else {
		fmt.Println("operation not permitted, exiting...")
		os.Exit(1)

	}

	fmt.Printf("\nApprove changes? (y/n)\n")
	var approved string
	fmt.Scanln(&approved)

	return approved

}
