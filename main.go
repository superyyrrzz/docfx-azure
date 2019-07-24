package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

	storage := fmt.Sprintf("%sdocfxsite", name)
	deployStorage(storage, group, sub)
}

func deployStorage(name string, group string, subscription string) {
	fmt.Println("Deploying storage...")
	execute(fmt.Sprintf("az storage account create -n %s -g %s --subscription %s --kind StorageV2", name, group, subscription))
	execute(fmt.Sprintf("az storage blob service-properties update --account-name %s --static-website --index-document index.html", name))
	fmt.Println("Finish deploying storage.")
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

/*
func execute(command string) string {
	fmt.Println(command)
	out, err := exec.Command("cmd", "/c", command).Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: %s", err, out))
	}
	return fmt.Sprintf("%s", out)
}
*/
