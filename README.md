# docfx-azure

## Prerequisites

* local:
  * git.exe
  * Azure CLI: Use to manage Azure resource. Ensure Azure DevOps extension installed:

    ```cli
    az extension add --name azure-devops
    ```

* service:
  * Azure subsciption / resource group
  * Azure DevOps / organization / project
    * Service connection to subscription established
    * Azure Repos / Azure Pipeline enabled

## What this tool do

* Create Azure Blob
* Create Azure Repos
* Create Azure Pipeline
