# Terraform Cloud Tool

This is a CLI to help changing and doing stuff in Terraform Cloud.

## Terraform CLI Functions

```  bash
$ terraform-cloud-tool
Terraform Cloud Tool is a tool to manage Terraform Cloud.

Usage:
  terraform-cloud-tool [flags]
  terraform-cloud-tool [command]

Available Commands:
  help        Help about any command
  variable    Variable function to Terraform Cloud
  workspace   Workspace function to Terraform Cloud

Flags:
  -h, --help       help for terraform-cloud-tool
      --o string   The organization to use to authenticate in TFCloud
      --t string   The token to use to authenticate in TFCloud

Use "terraform-cloud-tool [command] --help" for more information about a command.
```

This tool was made using [Cobra](https://github.com/spf13/cobra), so it's based with differnet commands, we'll list below:

### Workspaces

```  bash
$ terraform-cloud-tool workspace 
Create/Delete/Edit Workspaces in Terraform Cloud

Usage:
  terraform-cloud-tool workspace [flags]
  terraform-cloud-tool workspace [command]

Available Commands:
  apply       apply to workspaces
  delete      delete to workspaces

Flags:
  -h, --help        help for workspace
      --ts string   The tags to set in the workspace

Global Flags:
      --o string   The organization to use to authenticate in TFCloud
      --t string   The token to use to authenticate in TFCloud

Use "terraform-cloud-tool workspace [command] --help" for more information about a command.
```

Currently we have just the `apply` and `delete` Tags in Workspaces.

### Variable

``` bash
$ terraform-cloud-tool variable
Create/Delete/Edit Variables from Terraform Cloud

Usage:
  terraform-cloud-tool variable [flags]
  terraform-cloud-tool variable [command]

Available Commands:
  apply       Apply Variable function to Terraform Cloud
  delete      Delete Variable function to Terraform Cloud
  list        List variables in a Workspace
  read        Read variable in a Workspace

Flags:
  -h, --help         help for variable
      --vn string    Variable Name
      --vs           Variable Value is Sensitive
      --vv string    Variable Value
      --wt string    Filter the Workspace Name (REGEX)
      --wtg string   The tags to filter the workspaces

Global Flags:
      --o string   The organization to use to authenticate in TFCloud
      --t string   The token to use to authenticate in TFCloud

Use "terraform-cloud-tool variable [command] --help" for more information about a command.
```

Flags:

| Flag      | Description | Scope |
| --------- | ----------- | ----- |
| --o       | Organization Name | Global |
| --t       | The token responsible to authe in TFCloud | Global |
| --vn      | Variable Name | Variable |
| --vs      | if value is Sensitive | Variable |
| --vv      | Variable  Value| Variable |
| --wt      | Filter the workspace using Regex | Variable |
| --wtg     | Tags to filter workspaces | Variable |

All this flags will be used to variable command.

### Variable - List

```  bash
$ terraform-cloud-tool variable list --o "organization-stamps" --wt "testing-.*-test"
Listing all Variables in testing-1-2-3-test
Name: var1, Value: number1
Name: var2, Value: number2
Name: varSensitive1 
```

### Variable Read

```  bash
$ terraform-cloud-tool variable read --o "organization-stamps" --wt "testing-.*-test" --vn "var1"
Read Variable var1 for testing-1-2-3-test
Name: var1, Value: number1, Sensitive: false
```

### Variable Delete

```  bash
terraform-cloud-tool variable delete --o "organization-stamps" --wt "testing-.*-test" --vn "var1"
```

### Variable Apply

```  bash
terraform-cloud-tool variable apply --o "organization-stamps" --wt "testing-.*-test" --vn "var1" --vv "number1"
```
