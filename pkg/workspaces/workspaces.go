package workspace

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/edribeirojunior/terraform-cloud-tool/pkg/client"
	"github.com/hashicorp/go-tfe"
)

func Read(ws client.Workspaces) {
	//	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

	}
}

func Create(ws client.Workspaces) {
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

		if *ws.Version != "" {
			version := tfe.String(*ws.Version)
			_, err := ws.Cl.Workspaces.UpdateByID(ctx, ws.List[i].ID, tfe.WorkspaceUpdateOptions{
				TerraformVersion: version,
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Version Changed: %s, in Workspace %s\n", *ws.Version, ws.List[i].Name)
		}
	}
}

func ApproveChanges(ws client.Workspaces, action string) string {

	var wsList []string

	for i := 0; i < len(ws.List); i++ {
		wsList = append(wsList, ws.List[i].Name)
	}
	if action == "create" {
		if *ws.Tags != "" {
			fmt.Printf("The Tags will be created in these workspaces:\n")
			for _, i := range wsList {
				fmt.Printf("%s\n", i)
			}
		}
		if *ws.Version != "" {
			fmt.Printf("The Terraform version of environment will change in these workspaces:\n")
			for _, i := range wsList {
				fmt.Printf("%s\n", i)
			}
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
