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
							dumpPath := getReadPathOrExit(c, "input", "the Bimpf dump JSON")
							ncPath := getReadPathOrExit(c, "nextcloud-folder", "Bimpf Nextcloud project folder")
							if !fileExist(ncPath) {
								if c.Bool("ignore") {
									logrus.Warnf("Nextcloud folder at given path (%s) not found. As --ignore is used, execution will continue.")
								} else {
									logrus.Fatalf("Nextcloud folder at given path (%s) not found, this will lead to incorrect project files. Use --ignore to continue with the execution.", ncPath)
								}
							}
							outPath := getFolderPath(c, "output-folder", c.Bool("force"), true)

							dump := bimpf.OpenDump(dumpPath)
							project := dump.Convert(outPath, ncPath)
							project.SaveProject(outPath, !c.Bool("no-indentation"))
							return nil
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Usage:   "force overwrite of existing Acc project",
							},
							&cli.BoolFlag{
								Name:  "ignore",
								Usage: "ignore warnings",
							},
							&cli.StringFlag{
								Name:    "input",
								Aliases: []string{"i"},
								Usage:   "path to Bimpf dump JSON",
							},
							&cli.StringFlag{
								Name:    "nextcloud-folder",
								Aliases: []string{"n", "nc-folder"},
								Usage:   "path to nextcloud project folder which is used by Bimpf",
							},
							&cli.BoolFlag{
								Name:  "no-indentation",
								Usage: "suppress indentation of the output JSON files",
							},
							&cli.StringFlag{
								Name:    "output-folder",
								Aliases: []string{"output", "o"},
								Usage:   "path to the folder where the Acc project files should be written",
							},
						},
					},
					{
						Name:  "validate",
						Usage: "validate bimpf json dumps and searches for missing information in the Bimpf data",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "the Bimpf input file")
							outputPath := getPathOrExit(c, c.Bool("force"), "bimpf-dump-report.txt", "report", "the validation report")
							dump := bimpf.OpenDump(inputPath)
							dump.ValidateAndReport(outputPath)
							logrus.Info("report saved as ", outputPath)
							return nil
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Usage:   "force overwrite of existing report",
							},
							&cli.StringFlag{
								Name:    "input",
								Aliases: []string{"i"},
								Usage:   "path to Bimpf dump JSON",
							},
							&cli.StringFlag{
								Name:    "report",
								Aliases: []string{"r", "o"},
								Usage:   "path for the report",
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
				Action: func(c *cli.Context) error {
					outputPath := getFolderPath(c, "output-folder", c.Bool("force"), true)
					schema.NewProject(outputPath, false, true)
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "force overwrite of existing files",
					},
					&cli.StringFlag{
						Name:    "output-folder",
						Aliases: []string{"output", "o"},
						Usage:   "path to the folder where the Acc project files should be written",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

// getPathOrExit reads the content of the first argument or a given string flag and validates it.
// The application will exit, when the argument is not provided or the file already exists but doOverwrite is false.
// When a fallback path is provided it will be used when the user input is empty.
func getPathOrExit(c *cli.Context, doOverwrite bool, fallback, flag, to string) string {
	var pth string
	if flag == "" {
		if c.Args().Len() < 1 || c.Args().First() == "" {
			if fallback == "" {
				logrus.Fatalf("path to %s as first argument is needed", to)
			}
			logrus.Infof("as no path to %s is provided as first argument the default value (%s) will be used", to, fallback)
			pth = fallback
		} else {
			pth = c.Args().First()
		}
	} else {
		if c.String(flag) == "" {
			if fallback == "" {
				logrus.Fatalf("the flag -%s is needed as the path to %s", flag, to)
			}
			logrus.Infof("as no path to %s is provided with -%s the default value (%s) will be used", to, flag, fallback)
			pth = fallback
		} else {
			pth = c.String(flag)
		}
	}
	if !doOverwrite && fileExist(pth) {
		logrus.Fatalf("file (%s) already exists, use -f to overwrite files", pth)
	}
	return pth
}

// getReadPathOrExit does the same as getPathOrExit but does not check whether a file already exits.
// This function doesn't provide a fallback name, as an input file should be always explicit stated by the user.
func getReadPathOrExit(c *cli.Context, flag, to string) string {
	return getPathOrExit(c, true, "", flag, to)
}

// getFolderPath reads the content of first argument or a given string flag, validates and formats the folder path.
// It will also be checked, if the file can be overwritten. If no flag is provided (empty string) the first argument will be used.
// If mkDir is true, the folder tree will be created according to the user input.
func getFolderPath(c *cli.Context, flag string, doOverwrite, mkDir bool) string {
	var pth string
	useWd := false
	if flag == "" {
		if c.Args().Len() == 0 || c.Args().First() == "" {
			logrus.Info("will use working dir, as no folder path is given as an argument")
			useWd = true
		} else {
			pth = c.Args().First()
		}
	} else {
		if c.String(flag) == "" {
			logrus.Infof("will use working dir, aso no folder path is given with the -%s", flag)
			useWd = true
		} else {
			pth = c.String(flag)
		}
	}
	if useWd {
		var err error
		pth, err = os.Getwd()
		if err != nil {
			logrus.Error(err)
		}
	}
	if mkDir && !fileExist(pth) {
		logrus.Infof("will create «%s» as the folder doesn't exist", pth)
		if err := os.MkdirAll(pth, os.ModePerm); err != nil {
			logrus.Fatal(err)
		}
	}

	fi, err := os.Stat(pth)
	if err != nil {
		logrus.Fatal(err)
	}
	if !fi.IsDir() {
		logrus.Fatalf("given path (%s) is not a folder", pth)
	}
	if !doOverwrite && projectFilesExist(pth) {
		logrus.Fatalf("folder (%s) already exists, use -f to overwrite files", pth)
	}
	return pth
}

// projectFilesExist checks if there are no default project files existent.
// If this is the case, the application will be terminated.
func projectFilesExist(folderPath string) bool {
	for i := range schema.DefaultProjectFiles {
		pth := path.Join(folderPath, schema.DefaultProjectFiles[i])
		if fileExist(pth) {
			logrus.Infof("file %s exists", pth)
			return true
		}
	}
	return false
}

// fileExist checks if the file with the given path already exists.
func fileExist(pth string) bool {
	_, err := os.Stat(pth)
	return !os.IsNotExist(err)
}
