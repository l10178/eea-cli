package main

import (
	"github.com/l10178/eea-cli/cmd/eea/gitlab"
	"github.com/urfave/cli/v2"
	"os"
)

var app = cli.NewApp()

func info() {
	app.Name = "eea"
	app.Usage = "A cli tools"
}

func commands() {

	app.Commands = []*cli.Command{
		{
			Name:    "add-tag",
			Aliases: []string{"at"},
			Usage:   "Add new tag for git repository.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "project",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "The project id.",
				},
				&cli.StringFlag{
					Name:     "tag",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "The new tag name.",
				},
				&cli.StringFlag{
					Name:     "ref",
					Aliases:  []string{"r"},
					Required: true,
					Usage:    "Create tag using commit SHA, another tag name, or branch name.",
				},
				&cli.StringFlag{
					Name:    "message",
					Aliases: []string{"m"},
					Usage:   "The message.",
				},
			},
			Action: func(c *cli.Context) error {
				req := &gitlab.AddTagReq{
					ProjectId: c.String("project"),
					TagName:   c.String("tag"),
					Ref:       c.String("ref"),
					Message:   c.String("message"),
				}
				_, e := gitlab.AddTag(req)
				return e
			},
		},
		{
			Name:      "batch-tag",
			Aliases:   []string{"bt"},
			Usage:     "Batch tags by file.",
			UsageText: "eea batch-tag --file projects.txt --tag 1.2.0-release --ref master",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "file",
					Aliases:  []string{"f"},
					Required: true,
					Usage:    "The file contains all repositories.",
				},
				&cli.StringFlag{
					Name:     "tag",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "The new tag name.",
				},
				&cli.StringFlag{
					Name:    "ref",
					Aliases: []string{"r"},
					Usage:   "Create tag using another tag name or branch name.",
				},
			},
			Action: func(c *cli.Context) error {
				return gitlab.BatchTag(c.String("file"), c.String("tag"), c.String("ref"))
			},
		},
		{
			Name:      "batch-commit",
			Aliases:   []string{"bc"},
			Usage:     "Batch get tag's commit id by the file.",
			UsageText: "eea batch-commit --file projects.txt --tag 1.2.0-release",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "file",
					Aliases:  []string{"f"},
					Required: true,
					Usage:    "The file contains all repositories.",
				},
				&cli.StringFlag{
					Name:     "tag",
					Aliases:  []string{"t"},
					Required: true,
					Usage:    "The tag name.",
				},
			},
			Action: func(c *cli.Context) error {
				return gitlab.BatchCommit(c.String("file"), c.String("tag"))
			},
		},
	}
}

func main() {
	info()
	commands()
	_ = app.Run(os.Args)
}
