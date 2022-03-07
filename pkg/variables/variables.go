package variables

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

func (vr *Variable) Apply() {
	ctx := context.Background()

	for i := 0; i < len(vr.List); i++ {

		varID := ReadVarID(vr, vr.List[i])

		if varID != "" {
			variable, err := vr.Cl.Variables.Update(ctx, vr.List[i].ID, varID, tfe.VariableUpdateOptions{
				Key:       vr.Variables,
				Value:     vr.VariablesValue,
				Sensitive: vr.VariablesSensitive,
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Variable updated: %s, in Workspace: %s\n", variable.Key, vr.List[i].Name)
		}

		if *vr.Variables != "" && varID == "" {
			variable, err := vr.Cl.Variables.Create(ctx, vr.List[i].ID, tfe.VariableCreateOptions{
				Key:       vr.Variables,
				Value:     vr.VariablesValue,
				Sensitive: vr.VariablesSensitive,
				Category:  tfe.Category("terraform"),
			})

			if err != nil {
				fmt.Println("Error from here")
				log.Fatal(err)
			}

			fmt.Printf("Variable created: %s, in Workspace %s\n", variable.Key, vr.List[i].Name)

		}
	}
}

func (vr Variable) Read() string {
	ctx := context.Background()

	for i := 0; i < len(vr.List); i++ {
		fmt.Printf("Read Variable %s for %s\n", *vr.Variables, vr.List[i].Name)
		variable, _ := vr.Cl.Variables.List(ctx, vr.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})

		for j := 0; j < len(variable.Items); j++ {
			if variable.Items[j].Key == *vr.Variables {
				fmt.Printf("Name: %s, Value: %s,  Sensitive: %t\n", variable.Items[j].Key, variable.Items[j].Value, variable.Items[j].Sensitive)
				return variable.Items[j].ID
			}
		}

		fmt.Printf("Variable not found for Workspace %s\n", vr.List[i].Name)

	}

	return ""
}

func ReadVarID(vr *Variable, ws tfe.Workspace) string {
	ctx := context.Background()

	variable, _ := vr.Cl.Variables.List(ctx, ws.ID, tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	})

	for j := 0; j < len(variable.Items); j++ {
		if variable.Items[j].Key == *vr.Variables {
			return variable.Items[j].ID
		}
	}

	return ""

}

func (vr Variable) ListVars() {
	ctx := context.Background()

	for i := 0; i < len(vr.List); i++ {

		fmt.Printf("Listing all Variables in %s\n", vr.List[i].Name)
		variable, _ := vr.Cl.Variables.List(ctx, vr.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})

		for j := 0; j < len(variable.Items); j++ {
			if !variable.Items[j].Sensitive {
				fmt.Printf("Name: %s, Value: %s\n", variable.Items[j].Key, variable.Items[j].Value)
				// for _, j := range vrList {
				// 	fmt.Printf("Name: %s\n", j.Key)
				// }
			} else {
				fmt.Printf("Name: %s\n", variable.Items[j].Key)
			}
		}
		fmt.Println()
	}
}

func (vr Variable) Delete() {
	ctx := context.Background()

	for i := 0; i < len(vr.List); i++ {

		variable, err := vr.Cl.Variables.List(ctx, vr.List[i].ID, tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{PageSize: 100},
		})
		if err != nil {
			fmt.Println("Failed to list the variables")
			log.Fatal(err)
		}

		for j := 0; j < len(variable.Items); j++ {

			if variable.Items[j].Key == *vr.Variables {
				err = vr.Cl.Variables.Delete(ctx, vr.List[i].ID, variable.Items[j].ID)
				if err != nil {
					fmt.Println("Failed to delete variable")
					log.Fatal(err)
				}

				fmt.Printf("Variable deleted: %s, in Workspace %s\n", variable.Items[j].Key, vr.List[i].Name)
			}
		}
	}
}

func (vr Variable) ApproveChanges(action string) string {

	var vrList []string

	for i := 0; i < len(vr.List); i++ {
		vrList = append(vrList, vr.List[i].Name)
	}
	if action == "create" {
		fmt.Printf("The Variable %s will be created/updated in these workspaces:\n", *vr.Variables)
		for _, i := range vrList {
			fmt.Printf("%s\n", i)
		}
	} else if action == "delete" {
		fmt.Printf("The Variable %s will be deleted in these workspaces:\n", *vr.Variables)
		for _, i := range vrList {
			fmt.Printf("%s\n", i)
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
