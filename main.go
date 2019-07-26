package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "docfx-azure"
	app.Usage = "provision DocFX site on Azure"
	app.Author = "renzeyu"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:   "deploy",
			Usage:  "deploy a DocFX site on Azure",
			Action: deploy,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "site name, determine the name of resource created (required)",
				},
				cli.StringFlag{
					Name:  "subscription-id",
					Usage: "Azure subscription id (required)",
				},
				cli.StringFlag{
					Name:  "resource-group",
					Usage: "Azure resource group (required)",
				},
				cli.StringFlag{
					Name:  "organization-uri",
					Usage: "URI of Azure DevOps organization to host repository (required)",
				},
				cli.StringFlag{
					Name:  "project",
					Usage: "Azure DevOps project to host repository (required)",
				},
				cli.StringFlag{
					Name:  "service-connection",
					Usage: "service connection to Azure in Azure DevOps (required)",
				},
			},
		},
	}
}

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func deploy(c *cli.Context) {
	name := getFlag(c, "name")
	group := getFlag(c, "resource-group")
	sub := getFlag(c, "subscription-id")
	org := strings.TrimRight(getFlag(c, "organization-uri"), "/")
	proj := getFlag(c, "project")
	conn := getFlag(c, "service-connection")

	storage := fmt.Sprintf("%sdocfxsite", name)
	site := deployStorage(storage, group, sub)
	repo := deployRepoAndPipeline(name, org, proj, conn, storage)

	fmt.Println("")
	fmt.Println("All set!")
	fmt.Println("Wait a few minutes for your site's first build, then visit", site)
	fmt.Println("Your future changes in this created repository will be published automatically:", repo)
}

func deployStorage(name , group , subscription string) string {
	fmt.Println("Deploying storage...")
	execute(fmt.Sprintf("az storage account create -n %s -g %s --subscription %s --kind StorageV2", name, group, subscription))
	execute(fmt.Sprintf("az storage blob service-properties update --account-name %s --static-website --index-document index.html", name))
	fmt.Println("Finish deploying storage.")

	return strings.TrimSpace(execute(fmt.Sprintf("az storage account show -n %s -g %s --query primaryEndpoints.web --output tsv", name, group)))
}

func deployRepoAndPipeline(name , organization , project , connection , storage string) string {
	fmt.Println("Creating repository...")
	execute(fmt.Sprintf("az repos create --name %s --org %s -p %s", name, organization, project))
	dir := cloneTemplateRepo()
	updateTemplateRepo(dir, connection, storage)
	remote := fmt.Sprintf("%s/%s/_git/%s", organization, project, name)
	pushRepo(dir, remote)
	fmt.Println("Finish creating repository.")

	fmt.Println("Creating pipeline...")
	pipeline := fmt.Sprintf("%s-site-CI-CD", name)
	createPipeline(dir, pipeline, name)
	fmt.Println("Finish creating pipeline.")

	os.Remove(dir)
	return remote
}

func cloneTemplateRepo() string {
	dir, err := ioutil.TempDir(os.TempDir(), "docfx-azure-template")
	if err != nil {
		log.Fatal(err)
	}
	execute(fmt.Sprintf("git clone https://github.com/superyyrrzz/docfx-azure-template.git %s", dir))
	return dir
}

func updateTemplateRepo(dir , connection , storage string) {
	path := filepath.Join(dir, "azure-pipelines.yml")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	connUpdated := strings.Replace(string(content), "$SERVICE-CONNECTION", connection, 1)
	storageUpdated := strings.Replace(connUpdated, "$STORAGE", storage, 1)
	if err := ioutil.WriteFile(path, []byte(storageUpdated), 0666); err != nil {
		log.Fatal(err)
	}
}

func pushRepo(dir , remote string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	execute(fmt.Sprintf("git remote set-url origin %s", remote))
	execute("git add azure-pipelines.yml")
	execute(`git commit -m pipelines_updated`)
	execute("git push")
	err = os.Chdir(wd)
	if err != nil {
		log.Fatal(err)
	}
}

func createPipeline(dir , name , repoName string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	execute(fmt.Sprintf("az pipelines create --name %s --description docfx-azure-pipeline --repository %s --branch master --repository-type tfsgit --yml-path azure-pipelines.yml", name, repoName))
	err = os.Chdir(wd)
	if err != nil {
		log.Fatal(err)
	}
}

func getFlag(c *cli.Context, flag string) string {
	result := c.String(flag)
	if result == "" {
		log.Fatal(fmt.Sprintf("missing flag: %s", flag))
	}
	return result
}

func execute(str string) string {
	fmt.Println("[DEBUG] ", str)
	cmd := exec.Command("cmd", "/c", str)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprintf("%s", stdoutStderr))
	}
	return fmt.Sprintf("%s\n", stdoutStderr)
}
