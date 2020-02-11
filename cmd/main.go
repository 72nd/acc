package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.com/72th/acc/pkg/bimpf"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
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
				Name:  "bimpf",
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
							&cli.StringFlag{
								Name:  "nc-path",
								Usage: "path to nextcloud project folder which is used by Bimpf",
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
				Aliases: []string{"doc", "document"},
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
			{
				Name:    "new",
				Aliases: []string{"n", "create"},
				Usage:   "generates a new acc project with all needed files, use sub-commands to create only a subset",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "overwrite",
						Aliases:     []string{"o"},
						Usage:       "enable overwrite of existing files",
					},
				},
				Action: func(c *cli.Context) error {
					schema.NewProject(getFolderPath(c))
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

// getPathOrExit checks if the given context contain a string in the args, if not output an error and exit.
func getPathOrExit(c *cli.Context) string {
	if c.Args().Len() < 1 || c.Args().First() == "" {
		logrus.Fatal("path as a argument is needed")
	}
	return c.Args().First()
}

// getFolderPath reads the first argument, validates and formats the folder path.
// It will also be checked, if the file can be overwritten.
func getFolderPath(c *cli.Context) string {
	var pth string
	if c.Args().Len() == 0 || c.Args().First() == "" {
		logrus.Info("will use working dir, as no folder path is given as an argument")
		var err error
		pth, err = os.Getwd()
		if err != nil {
			logrus.Error(err)
		}
	}
	pth = c.Args().First()

	fi, err := os.Stat(pth)
	if pth != "" {
		if err != nil {
			logrus.Fatal(err)
		}
		if !fi.IsDir() {
			logrus.Fatalf("given path (%s) is not a folder")
		}
	}
	if !c.Bool("overwrite") && filesExist(pth) {
		logrus.Fatal("use -o to overwrite files")
	}
	return pth
}

// filesExist checks if there are no default project files existent.
// If this is the case, the application will be terminated.
func filesExist(folderPath string) bool {

	for i := range schema.DefaultProjectFiles {
		pth := path.Join(folderPath, schema.DefaultProjectFiles[i])
		if _, err := os.Stat(pth); err == nil {
			logrus.Infof("file %s exists", pth)
			return true
		}
	}
	return false
}
