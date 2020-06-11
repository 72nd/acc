package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.com/72th/acc/pkg/bimpf"
	"gitlab.com/72th/acc/pkg/camt"
	"gitlab.com/72th/acc/pkg/config"
	"gitlab.com/72th/acc/pkg/document/invoices"
	"gitlab.com/72th/acc/pkg/document/records"
	"gitlab.com/72th/acc/pkg/ledger"
	"gitlab.com/72th/acc/pkg/query"
	"gitlab.com/72th/acc/pkg/schema"
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
	completeFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:  "ask-skip",
			Usage: "ask if non valid entries should be skipped if id and identifier is set",
		},
		&cli.BoolFlag{
			Name:    "auto-save",
			Aliases: []string{"a"},
			Usage:   "save each transaction immediately after completion, this is useful when working with large bank statemens",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "redo all transaction which ID's and Identifier's are already set",
		},
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "acc project file",
		},
		&cli.BoolFlag{
			Name:    "open-attachment",
			Aliases: []string{"o"},
			Usage:   "open attachment (experimental feature)",
		},
		&cli.BoolFlag{
			Name:  "retain-focus",
			Usage: "try to retain focus when open attachment",
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
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.Parties.Customers = append(s.Parties.Customers, schema.NewPartyWithUuid())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.Parties.Customers = append(s.Parties.Customers, schema.InteractiveNewCustomer(s))
							}
							s.Save()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "employee",
						Aliases: []string{"emp"},
						Usage:   "add a employee",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.Parties.Employees = append(s.Parties.Employees, schema.NewPartyWithUuid())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.Parties.Employees = append(s.Parties.Employees, schema.InteractiveNewEmployee(s))
							}
							s.Save()
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
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.Expenses = append(s.Expenses, schema.NewExpenseWithUuid())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.Expenses = append(s.Expenses, schema.InteractiveNewExpense(&s, c.String("asset")))
							}
							s.Save()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "expense-category",
						Aliases: []string{"cat"},
						Usage:   "add one or multiple expense category",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.JournalConfig.ExpenseCategories = append(s.JournalConfig.ExpenseCategories, schema.NewExpenseCategory())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.JournalConfig.ExpenseCategories = append(s.JournalConfig.ExpenseCategories, schema.InteractiveNewExpenseCategories(c.Bool("multiple"))...)
							}
							s.Save()
							return nil
						},
						Flags: append(addFlags, &cli.BoolFlag{
							Name:    "mulitple",
							Aliases: []string{"m"},
							Usage:   "add multiple expende categories in one go",
						}),
					},
					{
						Name:    "invoice",
						Aliases: []string{"inv"},
						Usage:   "add a invoice",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.Invoices = append(s.Invoices, schema.NewInvoiceWithUuid())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.Invoices = append(s.Invoices, schema.InteractiveNewInvoice(s, c.String("asset")))
							}
							s.Save()
							return nil
						},
						Flags: addFlags,
					},
					{
						Name:    "misc-record",
						Aliases: []string{"mrc"},
						Usage:   "add a misc business record",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.MiscRecords = append(s.MiscRecords, schema.NewMiscRecord())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.MiscRecords = append(s.MiscRecords, schema.InteractiveNewMiscRecord(s, c.String("asset")))
							}
							s.Save()
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
							s := config.OpenSchema(inputPath)
							if c.Bool("default") {
								s.Statement.Transactions = append(s.Statement.Transactions, schema.NewTransactionWithUuid())
							} else {
								fmt.Println(aurora.BrightMagenta("Use the --default flag to suppress interactive mode and use defaults."))
								s.Statement.Transactions = append(s.Statement.Transactions, schema.InteractiveNewTransaction(s.Statement))
							}
							s.Save()
							return nil
						},
						Flags: addFlags,
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
							s := dump.Convert(outPath, ncPath)
							s.Save()
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
				Name:  "camt",
				Usage: "import bank-to-customer statement (camt.053.001.04)",
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					btcStatement := camt.NewBankToCustomerStatement(getReadPathOrExit(c, "statement", "camt xml file"))
					s := config.OpenSchema(inputPath)
					s.Statement.AddTransaction(btcStatement.Transactions())
					s.Save()
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
				Name:  "complete",
				Usage: "complete incorrect validated entries",
				Action: func(c *cli.Context) error {
					_ = cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "expenses",
						Usage: "complete incorrect expenses",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							s.Expenses.AssistedCompletion(&s, c.Bool("force"), c.Bool("auto-save"), c.Bool("open-attachment"), c.Bool("retain-focus"))
							s.Save()
							return nil
						},
						Flags: completeFlags,
					},
					{
						Name:  "invoices",
						Usage: "complete incorrect invoices",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							s.Invoices.AssistedCompletion(s, c.Bool("force"), c.Bool("auto-save"), c.Bool("open-attachment"), c.Bool("retain-focus"))
							s.Save()
							return nil
						},
						Flags: completeFlags,
					},
					{
						Name:  "repopulate",
						Usage: "repopulate expenses and invoices with transaction id's",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							s.Expenses.Repopulate(s)
							s.Invoices.Repopulate(s)
							s.MiscRecords.Repopulate(s)
							s.Save()
							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "input",
								Aliases: []string{"i"},
								Usage:   "acc project file",
							},
						},
					},
					{
						Name:  "transactions",
						Usage: "complete incorrect transactions",
						Action: func(c *cli.Context) error {
							inputPath := getReadPathOrExit(c, "input", "acc project file")
							s := config.OpenSchema(inputPath)
							s.Statement.AssistedCompletion(s, c.Bool("force"), c.Bool("auto-save"), c.Bool("auto-mode"), c.Bool("ask-skip"))
							s.Save()
							return nil
						},
						Flags: append(completeFlags, &cli.BoolFlag{
							Name:  "auto-mode",
							Usage: "set all transactions to auto mode, so third party has to be reviewed",
						}),
					},
				},
			},
			{
				Name:  "filter",
				Usage: "filter elements by date",
				Action: func(c *cli.Context) error {
					types := []string{"customers", "employee", "expenses", "invoices"}
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					if c.String("types") != "" {
						types = getSlice("types", c.String("types"), types)
					}
					from := getDateOrExit(c, "from")
					to := getDateOrExit(c, "to")
					s := config.OpenSchema(inputPath)
					s.Filter(types, from, to, c.String("output"), c.Bool("force"), c.String("identifier"))
					s.Save()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from",
						Aliases: []string{"f"},
						Usage:   "older elements are ignored, format YYYY-MM-DD",
					},
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"r"},
						Usage:   "overwrite filtered output files",
					},
					&cli.StringFlag{
						Name:  "identifier",
						Usage: "filter identifiers by `REGEX`",
					},
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "suffix for filtered output",
						Value:   "filtered",
					},
					&cli.StringFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "newer elements are ignored, format YYYY-MM-DD",
					},
					&cli.StringFlag{
						Name:  "types",
						Usage: "types to be filtered seperated by comma (customers,employee,expenses,invoices)",
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
					s := config.OpenSchema(inputPath)
					if c.Bool("all") {
						invoices.GenerateAllInvoices(
							s,
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
				Name:    "ledger",
				Aliases: []string{"ldg"},
				Usage:   "generate hledger journal",
				Action: func(c *cli.Context) error {
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					outputPath := getPathOrExit(c, c.Bool("force"), "transactions.journal", "output", "the journal file")
					s := config.OpenSchema(inputPath)
					journal := ledger.JournalFromAcc(s, c.Int("year"))
					journal.SaveHLedgerFile(outputPath)
					logrus.Info("journal saved as ", outputPath)
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
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "path for the journal file",
					},
					&cli.BoolFlag{
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "update transaction id's in expenses and invoices",
					},
					&cli.IntFlag{
						Name:    "year",
						Aliases: []string{"y"},
						Usage:   "generate journal for specific year",
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
						fmt.Println(aurora.BrightMagenta("assistant for new acc project, use --default for non interactive use"))
					}
					config.NewSchema(outputPath, c.String("logo"), true, !c.Bool("default"))
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
					&cli.BoolFlag{
						Name:    "project-mode",
						Aliases: []string{"p"},
						Usage:   "enable project mode",
					},
				},
			},
			{
				Name:  "query",
				Usage: "find and display elements",
				Subcommands: []*cli.Command{
					{
						Name:  "types",
						Usage: "list all available element types",
						Action: func(c *cli.Context) error {
							query.AccQueryables.PPTypes()
							return nil
						},
					},
					{
						Name:  "keys",
						Usage: "list keys for given elemt type",
						Action: func(c *cli.Context) error {
							qry, err := query.AccQueryables.QueryablesFromUserInput(c.String("types"))
							if err != nil {
								logrus.Fatal(err)
							}
							mode := query.TableMode
							if c.Bool("yaml") {
								mode = query.YamlMode
							}
							qry.PPKeys(mode)
							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "types",
								Aliases: []string{"t"},
								Usage:   "types to be filtered seperated by comma, use 'acc query types' to get possibilities",
							},
							&cli.BoolFlag{
								Name:  "yaml",
								Usage: "output as YAML",
							},
						},
					},
				},
				Action: func(c *cli.Context) error {
					qry, err := query.AccQueryables.QueryablesFromUserInput(c.String("types"))
					if err != nil {
						logrus.Fatal(err)
					}
					inputPath := getReadPathOrExit(c, "input", "acc project file")
					s := config.OpenSchema(inputPath)
					mode := query.TableMode
					if c.Bool("yaml") {
						mode = query.YamlMode
					}
					qry.QueryAcc(s, c.String("match"), c.String("date"), c.String("select"), mode, !c.Bool("no-render"), c.Bool("strict"), c.Bool("open-attachment"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "filter keys by date ranges key:from:to as `REGEX:YYYY-MM-DD:YYYY-MM-DD` multiple can be seperated by comma",
					},
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "acc project file",
					},
					&cli.StringFlag{
						Name:    "match",
						Aliases: []string{"r"},
						Usage:   "match key:value compbinations with `REGEX:REGEX` multiple can be seperated by comma",
					},
					&cli.BoolFlag{
						Name:    "no-render",
						Aliases: []string{"n"},
						Usage:   "do not render the output values",
					},
					&cli.BoolFlag{
						Name:    "open-attachment",
						Aliases: []string{"o"},
						Usage:   "open attachment (experimental feature)",
					},
					&cli.StringFlag{
						Name:    "select",
						Aliases: []string{"s"},
						Usage:   "select displayed keys, multiple can be sperated by comma `KEY[,KEY]`",
					},
					&cli.BoolFlag{
						Name:  "strict",
						Usage: "case sensitive matching",
					},
					&cli.StringFlag{
						Name:    "types",
						Aliases: []string{"t"},
						Usage:   "types to be filtered seperated by comma, use 'acc query types' to get possibilities",
					},
					&cli.BoolFlag{
						Name:  "yaml",
						Usage: "output as YAML",
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
					s := config.OpenSchema(inputPath)
					s = s.FilterYear(c.Int("year"))
					records.GenerateExpensesRec(s, c.String("output-folder"), c.Bool("do-overwrite"), !c.Bool("skip-downconvert"))
					records.GenerateInvoicesRec(s, c.String("output-folder"), c.Bool("do-overwrite"), !c.Bool("skip-downconvert"))
					records.GenerateMiscsRec(s, c.String("output-folder"), c.Bool("do-overwrite"), !c.Bool("skip-downconvert"))
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
					&cli.BoolFlag{
						Name:  "skip-downconvert",
						Usage: "skip downconvert with pdftops and gs",
					},
					&cli.IntFlag{
						Name:    "year",
						Aliases: []string{"y"},
						Usage:   "generate journal for specific year",
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
					s := config.OpenSchema(inputPath)
					s.ValidateAndReportProject(outputPath)
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

func getDateOrExit(c *cli.Context, flag string) *time.Time {
	if c.String(flag) == "" {
		return nil
	}
	value, err := time.Parse("2006-01-02", c.String(flag))
	if err != nil {
		logrus.Fatalf("value «%s» from flag --%s could not be parsed with layout YYYY-MM-DD", flag, c.String(flag))
	}
	return &value
}

// projectFilesExist checks if there are no default project files existent.
// If this is the case, the application will be terminated.
func projectFilesExist(folderPath string) bool {
	for i := range config.DefaultProjectFiles {
		pth := path.Join(folderPath, config.DefaultProjectFiles[i])
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

func getSlice(name, input string, valid []string) []string {
	ele := strings.Split(input, ",")
	result := make([]string, len(ele))
	for i := range ele {
		e := strings.TrimPrefix(ele[i], " ")
		e = strings.TrimSuffix(e, " ")
		contained := false
		for j := range valid {
			if e == valid[j] {
				contained = true
				break
			}
		}
		if !contained {
			logrus.Fatalf("%s is not alloweed for --%s, use \"%s\"", e, name, valid)
		}
		result[i] = e
	}
	return result
}
