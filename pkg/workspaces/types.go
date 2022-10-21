package workspace

import "github.com/hashicorp/go-tfe"

type Workspace struct {
	Cl               *tfe.Client
	List             []tfe.Workspace
	ListLocked       []tfe.Workspace
	Tags             *string `json:",omitempty"`
	TerraformVersion *string `json:",omitempty"`
	RunsList         []tfe.RunList
}
