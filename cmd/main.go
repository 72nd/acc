package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.com/72th/acc/pkg/bimpf"
	"gitlab.com/72th/acc/pkg/camt"
	"gitlab.com/72th/acc/pkg/document/invoices"
	"gitlab.com/72th/acc/pkg/document/records"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
	"path"
)

func main() {
	addFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "asset",
			Aliases: []string{"a"},
			Usage:   "path to the asset file",
		},
		&cli.BoolFlag{
			Name:    "default",
			Aliases: []string{"d", "default-values"},
			Usage:   "use default values and do not use interactive input",
		},
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "acc project file",
		},
	}

	app := &cli.App{
		Name:                 "acc",
		Usage:                "collection of tools for accounting with hledger",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add an element (expense, invoice etc.) to a project",
				Action: func(c *cli.Context) error {
					_ = cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "customer",
						Aliases: []string{"cst"},
						Usage:   "add a customer",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							acc := schema.OpenProject(inputPath)
							if c.Bool("default") {
								acc.Parties.Customers = append(acc.Parties.Customers, schema.NewPartyWithUuid())
							} else {
								acc.Parties.Customers = append(acc.Parties.Customers, schema.InteractiveNewCustomer(acc))
							}
							acc.SaveProject()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "employee",
						Aliases: []string{"epy"},
						Usage:   "add a employee",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							acc := schema.OpenProject(inputPath)
							if c.Bool("default") {
								acc.Parties.Employees = append(acc.Parties.Employees, schema.NewPartyWithUuid())
							} else {
								acc.Parties.Employees = append(acc.Parties.Employees, schema.InteractiveNewEmployee(acc))
							}
							acc.SaveProject()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "expense",
						Aliases: []string{"exp"},
						Usage:   "add a expense",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							acc := schema.OpenProject(inputPath)
							if c.Bool("default") {
								acc.Expenses = append(acc.Expenses, schema.NewExpenseWithUuid())
							} else {
								acc.Expenses = append(acc.Expenses, schema.InteractiveNewExpense(acc, c.String("asset")))
							}
							acc.SaveProject()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "invoice",
						Aliases: []string{"inv"},
						Usage:   "add a invoice",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							acc := schema.OpenProject(inputPath)
							if c.Bool("default") {
								acc.Invoices = append(acc.Invoices, schema.NewInvoiceWithUuid())
							} else {
								acc.Invoices = append(acc.Invoices, schema.InteractiveNewInvoice(acc, c.String("asset")))
							}
							acc.SaveProject()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "transaction",
						Aliases: []string{"trn"},
						Usage:   "add a transaction",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							acc := schema.OpenProject(inputPath)
							if c.Bool("default") {
								acc.BankStatement.Transactions = append(acc.BankStatement.Transactions, schema.NewTransactionWithUuid())
							} else {
								acc.BankStatement.Transactions = append(acc.BankStatement.Transactions, schema.InteractiveNewTransaction(acc.BankStatement))
							}
							acc.SaveProject()
							return nil
						},
						Flags: addFlags,
					},
				},
			},
			{
				Name:  "bank",
				Usage: "import bank-to-customer statement (camt.053.001.04)",
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					btcStatement := camt.NewBankToCustomerStatement(getReadPathOrExit(c, "input", "camt xml file"))
					acc := schema.OpenProject(inputPath)
					acc.BankStatement.AddTransaction(btcStatement.Transactions())
					acc.SaveProject()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "statement",
						Aliases: []string{"s"},
						Usage:   "path to camt xml",
					},
				},
			},
			{
				Name:  "bimpf",
				Usage: "bimpf related functions",
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
							project.SaveProjectToFolder(outPath)
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
							&cli.StringFlag{
								Name:    "output-folder",
								Aliases: []string{"output", "o"},
								Usage:   "path to the folder where the Acc project files should be written",
							},
						},
					},
					{
						Name:  "validate",
						Usage: "validate bimpf json dumps and searches for missing information in the Bimpf utils",
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
				Name:  "invoices",
				Usage: "generate simple invoices based on a project",
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					if err := os.MkdirAll(c.String("output-folder"), os.ModePerm); err != nil {
						logrus.Fatal("creation of document output folder failed: ", err)
					}
					acc := schema.OpenProject(inputPath)
					if c.Bool("all") {
						invoices.GenerateAllInvoices(
							acc,
							c.String("output-folder"),
							c.String("place"),
							c.Bool("do-overwrite"),
						)
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "export all existing invoices",
					},
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "output-folder",
						Aliases: []string{"output", "o"},
						Value:   "invoices",
						Usage:   "path to the folder where the exported documents should be stored",
					},
					&cli.BoolFlag{
						Name:    "do-overwrite",
						Aliases: []string{"overwrite"},
						Value:   false,
						Usage:   "force overwrite existing documents",
					},
					&cli.StringFlag{
						Name:  "place",
						Value: "PLACE-UNSET",
						Usage: "place where the invoice originates from",
					},
				},
			},
			{
				Name:    "new",
				Aliases: []string{"n", "nw", "create"},
				Usage:   "generates a new acc project with all needed files, use sub-commands to create only a subset",
				Action: func(c *cli.Context) error {
					outputPath := getFolderPath(c, "output-folder", c.Bool("force"), true)
					if !c.Bool("default") {
						fmt.Println("assistant for new acc project, use --default for non interactive use")
					}
					schema.NewProject(outputPath, c.String("logo"), true, !c.Bool("default"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "default",
						Aliases: []string{"d", "default-values"},
						Usage:   "use default values and do not use interactive input",
					},
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "force overwrite of existing files",
					},
					&cli.StringFlag{
						Name:  "logo",
						Usage: "path to the logo file",
					},
					&cli.StringFlag{
						Name:    "output-folder",
						Aliases: []string{"output", "o"},
						Usage:   "path to the folder where the Acc project files should be written",
					},
				},
			},
			{
				Name:    "records",
				Aliases: []string{"rec"},
				Usage:   "aggregate all business records associated with a type of business case",
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					if err := os.MkdirAll(c.String("output-folder"), os.ModePerm); err != nil {
						logrus.Error("creation of document output folder failed: ", err)
					}
					acc := schema.OpenProject(inputPath)
					records.GenerateExpensesRec(acc, c.String("output-folder"), c.Bool("do-overwrite"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "output-folder",
						Aliases: []string{"output", "o"},
						Value:   "records",
						Usage:   "path to the folder where the exported documents should be stored",
					},
					&cli.BoolFlag{
						Name:    "do-overwrite",
						Aliases: []string{"overwrite", "r"},
						Value:   false,
						Usage:   "force overwrite existing documents",
					},
				},
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "validates the current project",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "force overwrite of existing report",
					},
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "report",
						Aliases: []string{"r", "o"},
						Usage:   "path for the report",
					},
				},
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					outputPath := getPathOrExit(c, c.Bool("force"), "acc-report.txt", "report", "the validation report")
					acc := schema.OpenProject(inputPath)
					acc.ValidateAndReportProject(outputPath)
					logrus.Info("report saved as ", outputPath)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

/*
// validValueOrExit checks the existence of a given string flag and it's content against a list of allowed variables.
// If both tests are positive the function returns the content of the flag otherwise it quits the application.
func getValidValueOrExit(c *cli.Context, flag string, allowed []string) string {
	content := c.String(flag)
	if content == "" {
		logrus.Fatalf("flag %s has to be provided", flag)
	}
	for i := range allowed {
		if allowed[i] == content {
			return content
		}
	}
	logrus.Fatalf("flag %s was provided with an illegal value (%s). Allowed: %+v", flag, content, allowed)
	return ""
}
*/

// getPathOrExit reads the content of the first argument or a given string flag and validates it.
// The application will exit, when the argument is not provided or the file already exists but doOverwrite is false.
// When a fallback path is provided it will be used when the user input is empty.
func getPathOrExit(c *cli.Context, doOverwrite bool, fallback, flag, to string) string {
	var pth string
	if flag == "" {
		if c.Args().Len() < 1 || c.Args().First() == "" {
			if fallback == "" {
				_ = cli.ShowCommandHelp(c, c.Command.Name)
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
				_ = cli.ShowCommandHelp(c, c.Command.Name)
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
