package main

import (
	"eea-cli/registry"
	"github.com/urfave/cli"
	"log"
	"os"
)

var app = cli.NewApp()

func info() {
	app.Name = "eea"
	app.Usage = "A cli tools"
	app.Author = "nxest.com"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "registry",
			Aliases: []string{"r"},
			Usage:   "Get docker registry's latest login token.",
			Action: func(c *cli.Context) {
				registry.QueryRegistryToken()
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
