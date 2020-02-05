package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.com/72th/acc/pkg/bimpf"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "acc",
		Usage: "collection of tools for accounting with hledger",
		Action: func(c *cli.Context) error {
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "bimpf",
				Usage: "Bimpf related functions",
				Action: func(c *cli.Context) error {
					_ = cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "import",
						Usage: "import Bimpf json dumps and converts them to acc json compatible files",
						Action: func(c *cli.Context) error {
							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "input",
								Aliases: []string{"i"},
								Usage:   "path to bimpf dump json",
							},
						},
					},
					{
						Name:  "validate",
						Usage: "validate bimpf json dumps and searches for missing information in the Bimpf data",
						Action: func(c *cli.Context) error {
							if c.String("input") == "" {
								logrus.Fatal("no input file specified (-i)")
							}
							inputPath := c.String("input")
							outputPath := c.Args().First()
							if outputPath == "" {
								outputPath = "bimpf-dump-report.txt"
							}
							dump := bimpf.OpenDump(inputPath)
							dump.ValidateAndReport(outputPath)
							logrus.Info("report saved as ", outputPath)
							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "input",
								Aliases: []string{"i"},
								Usage:   "path to bimpf dump json",
							},
						},
					},
				},
			},
			{
				Name:    "documents",
				Aliases: []string{"doc"},
				Usage:   "aggregate all documents associated with a type of business case",
				Action: func(c *cli.Context) error {
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "expenses",
						Aliases: []string{"e", "expense"},
						Usage:   "path to expenses json file",
					},
					&cli.StringFlag{
						Name:    "invoice",
						Aliases: []string{"i", "invoices"},
						Usage:   "path to invoices json file",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
