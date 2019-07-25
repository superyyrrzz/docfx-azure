# docfx-azure

## Usage

```
docfx-azure deploy --resource-group {RES} --subscription-id {SUB} --organization-uri {ORG} --project {PROJ} --service-connection {CONN} --name {NAME}
```

## Prerequisites

* local tools:
  * git.exe
  * Azure CLI: Use to manage Azure resource. Ensure Azure DevOps extension installed, and already signed in:

    ```cli
    az login
    az extension add --name azure-devops
    ```

* existing service:
  * Azure subsciption
  * Azure DevOps

## What this tool do

* Create Azure Blob
* Create Azure Repos
* Create Azure Pipeline
