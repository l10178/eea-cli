package main

import (
	"github.com/l10178/eea-cli/cmd/eear/registry"
	"github.com/urfave/cli/v2"
	"os"
)

var app = cli.NewApp()

func info() {
	app.Name = "eear"
	app.Usage = "Yeah, my docker registry tools."
}

func commands() {

	app.Commands = []*cli.Command{
		{
			Name:    "registry",
			Aliases: []string{"r"},
			Usage:   "Get docker registry's latest login token.",
			Action: func(c *cli.Context) error {
				registry.QueryRegistryToken()
				return nil
			},
		},
	}
}

func main() {
	info()
	commands()
	_ = app.Run(os.Args)
}
