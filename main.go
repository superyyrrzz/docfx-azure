package main

import (
	"fmt"
	"log"
	"os"

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
					Name:  "subscription-id",
					Usage: "Azure subscription id",
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
	fmt.Println("deploying ...")
	fmt.Println("subscription id:", c.String("subscription-id"))
	fmt.Println("Complete!")
}
