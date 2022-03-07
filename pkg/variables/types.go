package variables

import "github.com/hashicorp/go-tfe"

type Variable struct {
	Cl                 *tfe.Client
	List               []tfe.Workspace
	Variables          *string `json:",omitempty"`
	VariablesValue     *string `json:",omitempty"`
	VariablesSensitive *bool   `json:",omitempty"`
}
