# docfx-azure

[![Build Status](https://dev.azure.com/renzeyu/docfx-azure/_apis/build/status/docfx-azure.CI?branchName=master)](https://dev.azure.com/renzeyu/docfx-azure/_build/latest?definitionId=7&branchName=master)

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
