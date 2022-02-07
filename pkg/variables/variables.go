package variables

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/edribeirojunior/terraform-cloud-tool/pkg/client"
	"github.com/hashicorp/go-tfe"
)

func Create(ws client.Workspaces) {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		if *ws.Variables != "" {
			variable, err := ws.Cl.Variables.Create(ctx, ws.List[i].ID, tfe.VariableCreateOptions{
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
	}
}

func Read(ws client.Workspaces) {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		fmt.Printf("Listing all Variables in %s\n", ws.List[i].Name)
		variable, _ := ws.Cl.Variables.List(ctx, ws.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})

		for j := 0; j < len(variable.Items); j++ {
			if !variable.Items[j].Sensitive {
				fmt.Printf("Name: %s, Value: %s\n", variable.Items[j].Key, variable.Items[j].Value)
				// for _, j := range wsList {
				// 	fmt.Printf("Name: %s\n", j.Key)
				// }
			} else {
				fmt.Printf("Name: %s\n", variable.Items[j].Key)
			}
		}
		fmt.Println()
	}
}

func Delete(ws client.Workspaces) {
	ctx := context.Background()

	for i := 0; i < len(ws.List); i++ {

		variable, err := ws.Cl.Variables.List(ctx, ws.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})
		if err != nil {
			fmt.Println("Failed to list the variables")
			log.Fatal(err)
		}

		for j := 0; j < len(variable.Items); j++ {

			if variable.Items[j].Key == *ws.Variables {
				err = ws.Cl.Variables.Delete(ctx, ws.List[i].ID, variable.Items[j].ID)
				if err != nil {
					fmt.Println("Failed to delete variable")
					log.Fatal(err)
				}

				fmt.Printf("Variable deleted: %s, in Workspace %s\n", variable.Items[j].Key, ws.List[i].Name)
			}
		}
	}
}

func ApproveChanges(ws client.Workspaces, action string) string {

	var wsList []string

	for i := 0; i < len(ws.List); i++ {
		wsList = append(wsList, ws.List[i].Name)
	}
	if action == "create" {
		fmt.Printf("The Variable %s will be created in these workspaces:\n", *ws.Variables)
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
	} else if action == "delete" {
		fmt.Printf("The Variable %s will be deleted in these workspaces:\n", *ws.Variables)
		for _, i := range wsList {
			fmt.Printf("%s\n", i)
		}
	} else {
		fmt.Println("operation not permitted, exiting...")
		os.Exit(1)

	}

	fmt.Printf("\n Approve changes? (y/n)\n")
	var approved string
	fmt.Scanln(&approved)

	return approved

}
