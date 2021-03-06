# acc

 <p align="center">
  <img width="128" height="128" src="misc/icon-text.png">
</p>

First planned as a simple tool chain for easing the work with the plain text accounting software [hledger](https://hledger.org/), acc evolved into some sort of plain-text ERP system. It's capable of keeping track of your customers, employees, expenses and invoices as well as importing bank statements (via [ISO 20022](https://en.wikipedia.org/wiki/ISO_20022) camt). Acc can generate most of your hledger account records based on this data. Also it's possible to export all business records (expenses etc.) with their respective receipt as a PDF for further archiving. There is also a quit powerful query functionality for easy finding of records. For adding new data, a interactive prompt can be used (you also can just edit the [YAML](https://en.wikipedia.org/wiki/YAML) files by yourself).

_Current status:_ The project is under active use and development at the [Genossenschaft Solutionsbüro](https://buero.io) and Acc evolves along our needs. Almost all core features are implemented and we use this software in your everyday work. The biggest throwback in the moment are the missing tests and a bad code documentation (I'm working on it). Also a manual (or at least a detailed and long README) is missing at this point.

![Screenshots](misc/screenshots.png)


## Table of contents

* [Todo](#todo)
* [Installation](#installation)
* [Basic Concepts](#basic-concepts)
    + [Types of data/records](#types-of-data-records)
    + [Modes](#modes)
* [Usage and Functions](#usage-and-functions)
    + [add](#add)
    + [bimpf](#bimpf)
	* [camt](#camt)
	+ [complete](#complete)
	+ [distributed](#distributed)
	+ [filter](#filter)
	+ [invoices](#invoices)
	+ [ledger](#ledger)
	+ [new](#new)
	+ [query](#query)
	+ [records](#records)
	+ [validate](#validate)
* [Workflows](#workflows)
	+ [Workflow with Bimpf](#workflow-with-bimpf)
	+ [Simple theater project](#simple-theater-project)


## Todo

- [x] First Run with 2019 Taxes
- [ ] Documentation in Readme done
- [x] Filter function
- [x] Project Mode
- [x] Shorts for commands 
- [ ] Code Documentation
- [ ] Misc Documents in complete Transactions
- [ ] Debug flag
- [ ] Multiple Bank Statements
- [x] Amounts with Amount type
- [ ] pain.001 for payment generation
	- [ ] Add IBAN to employee (data field, interactive add, assisted completion, validation)
	- [x] Rewrite `sba-pay` as library for go
	- [ ] Add functionality to iso20022 package



## Installation

### Perquisites

First install perquisites (optional step): 

```shell script
# For experimental features:
sudo apt install wmctrl xdg-mime xprop

# Export embedded PDF's in records
sudo apt install poppler-utils ghostscript
```

### With `go install` 

```shell script
go install github.com/72nd/acc
```


### From binary file

There are some pre-compiled binaries in the [release section](https://github.com/72nd/acc/releases/latest). Download the binary for your system, make it executable and make sure that the binary is in the `PATH`.


### From source

You can also download the source code and build Acc yourself. Make sure [go is installed](https://golang.org/doc/install).

```shell script
wget https://github.com/72nd/acc/archive/v.0.4.7.tar.gz
tar -zxvf v.0.4.7.tar.gz
cd acc-v.0.4.7
go build -o acc main.go
```


### Tab Completion

Acc builds on [urfave's cli framework](https://github.com/urfave/cli/) thus to activate command completion do the following:

#### ZSH

```shell script
cd path/to/store/autocomplete-script/
wget https://raw.githubusercontent.com/urfave/cli/master/autocomplete/zsh_autocomplete
chmod +x zsh_autocomplete
```

Open your ZSH config file (e.g. `.zshrc`) and add the following lines:

```shell script
PROG=acc
_CLI_ZSH_AUTOCOMPLETE_HACK=1
source path/to/store/autocomplete-script/zsh_autocomplete
```


#### Bash

```shell script
cd path/to/store/autocomplete-script/
wget https://raw.githubusercontent.com/urfave/cli/master/autocomplete/bash_autocomplete
chmod +x bash_autocomplete
```

Open your bash config file and add this line:

```shell script
PROG=acc source path/to/store/autocomplete-script/bash_autocomplete
``` 


### Build for multiple platforms

This project uses the [Task](https://github.com/go-task/task) task runner to build this application for multiple platforms. With `task` present run `task build`in the repository root. You'll find all the binaries in the `build` folder.

To build the deb packages:

```shell script
// Dependencies
go install github.com/mh-cbon/go-bin-deb
sudo apt install build-essential lintian

// Build
task deb
```

## Basic Concepts

To start with Acc, it's important to understand some fundamental concepts and ideas which shape the handling of the software.

The idea of Acc is to collect data about your business and then generating a number of different outputs from these. This diagram should give you an idea about some possibilities of Acc:

![Example flow of data](misc/flow-of-data.svg)


### Types of data/records

**config** Commonly named `acc.yaml` contains all the basic data about a Acc project (learn morn about this below in the «Modes» section) as well as all the definitions for the automatic account records generation. This file is also the entry point for the application. While using Acc you always have to state this file with the `-i` flag (exception: generation of a new Acc project with `acc new`).

**expense** Expenses represent an event, where the company has to pay some money. This can be the receiving of a bill (eg. tax bill) or paying a purchase directly with the companies debit card. But also the advancing employee scenario can be handled. Sometimes an employee has to pay something with his/her own money. Acc provides functionality to keep track of such advanced expenses and also generating payment order files (ISO 20022 pain.001) for easy transferring your debts. 

**invoice** Represents an invoice you've sent to a customer. Acc also contains an experimental feature to render simple invoice letters as PDFs.

**misc-record** Sometimes there are other documents which have to be archived or are the cause for some transaction on the bank account (example: the final account of the health insurance which states a refund). As this documents don't fit into the other categories, there is this misc category.

**party** Either a customer or an employee containing the usual information (name, street etc.). Most of the other documents or records are somehow associated with one or multiple parties (ex: projects belong to a customer, an invoice was sent to a customer, a expense was advanced by a employee). Learn more about this interconnections in the diagram below.

**project** If you use Acc for a more complex scenario it makes sense to group expenses and invoices per customer project. Expenses and invoices can be linked to a project. Each project has a associated customer. By using _distributed mode_ you can also group your files in project folders (learn more about below). The use of projects is optional, for simple cases you don't have to use them.

**statement** A bank statement contains bank account transactions for a certain period of time. In the future a Acc project should be able to have multiple statements separated by a period (month, year).

**transaction** A transaction describes the receiving or payment of a amount on your bank account. Some data types (like invoices and expenses) can be associated with one or multiple transactions. This way you can keep track of the payment of your invoices and transfer outstanding advanced employee balances. Acc can import this transactions directly from your bank account via ISO 20022 pain.001. You can learn more on how to link records to (imported) transactions in the _complete_ sub-command section.

Mainly for reference: This diagram shows all possible interconnections between the different types. If this confuses you, just ignore it for now.

![Data Model](misc/data-model.svg)


### Modes

Acc has two different ways of saving all the data: _flat_ and _distributed_ mode. Both modes using YAML files but differ in the way of arranging the data in the file system.

**flat mode** In flat mode there is one YAML file per record type. Typically all files are in the same folder and the `acc.yaml` file contains the paths to all other files. This is the default mode and fits small projects (like a single theater production) well. A flat mode file structure will look something like this:

```
.
├── acc.yaml
├── bank-statement.yaml
├── expenses.yaml
├── invoices.yaml
├── misc.yaml
├── parties.yaml
└── projects.yaml

```

**distributed mode** The distributed mode on the other hand organizes the records according to their customer and project using not only files but also a folder structure.  The aim is to use Acc for multiple projects, customers and over the curse of multiple years. A distributed Acc project consists of a base folder, containing the `acc.yaml` and `employees.yaml` files and the folders `internal` (containing all records which are internal and thus not associated with any customer and/or project) and `projects`. The `projects` folder contains a directory for each customer. Each customer folder on the other hand contains a `customer.yaml` (with the data for this customer) and a folder for each project associated with this customer. Invoices and expenses linked to a project are stored in a `project.yaml` file present in each and every project folder. Please note: This mode is much more opinionated on how to store and arrange the data.

```
.
├── acc.yaml
├── employees.yaml
├── internal
│   └── expenses-2019.yaml
└── projects
    └── max-mustermann
        ├── building-a-space-rocket
        │   └── project.yaml
        └── customer.yaml

```



## Usage and Functions

You can always add the `--help` flag to all commands to learn more about a certain (sub) command.

### add

Add new elements (customer, employee, expense, expense-category, invoice or transaction) to your acc project. If you don't want to use the interactive prompt, use the `--default` flag. Some of the elements contain paths to files by using the `--asset` flag you can specify this paths in advance and thus use the tab-completion of your shell.

```shell script
acc add customer -i acc.yaml
acc add invoice -i acc.yaml --asset /path/to/sent-invoice.pdf
```


### bimpf 

This command allows to import JSON dumps of our old administration software [Bimpf](https://gitlab.com/solutionsbuero/bimpf), you'll probably never need this :)

Hint: Dump with the following commands on the server: 

```shell script
git clone https://gitlab.com/solutionsbuero/bimpf.git
cd bimpf
python3 -m venv init .
source venv/bin/activate
pip3 install -r requirements.txt
python3 setup.py install
python3 bimpf/cli/__main__.py dump.json
```

Import Bimpf JSON dumps and converts them into to the acc-formatted JSON files:

```shell script
acc bimpf import -i bimpf.json PATH/TO/FOLDER/
```

Validate a Bimpf JSON dump and saves the result to an output file:

```shell script
acc bimpf validate -i bimpf.json report.txt
```


### camt

Import bank statements from a ISO 20022 camt.001 xml into your acc project.

```shell script
acc camt -i acc.yaml -statement /path/to/camt.xml
```


### complete

Interactive complete expenses, invoices and transactions. Runs trough all elements and prompts an interface to complete certain elements. It supports searching for linked elements (like obliged customer for an invoice). Flags:

- `--auto-save` Save each element record immediately after each completion. This can be helpful when working with big data sets which can't be all completed in one go.
- `force` Normally only not valid or completed elements will be prompted for completion. Using this flag all elements will be prompted for completion. Caution: You can loose data.
- `--input` as usual the path to the `acc.yaml` file.
- `--open-attachement` _Experimentally feature!_ Open (and most of the time) close the associated file in the default application. Will only work on Linux with `xdg-mime`, `xprop` installed and the application desktop files located under `/usr/share/applications/`. 
- `--retain-focus` _Hacky feature!_ Tries to regain focus of the terminal acc runs in. Will only work on certain Linux installations if `wmctrl` is installed.

The `acc complete repopulate` on the other hand can be used to link expenses and invoices to transactions which already are linked to the expense/invoice.



### distributed

Place for some utility commands for working in distributed mode.


### filter

Filter expenses and invoices with a data range and saves the subset in new files. Feature not completed.


### invoices

_Experimentally feature!_ Create some very basic invoices for customers.


### ledger

Create a [hleder](https://hledger.org) journal based on the acc project.


### new

Creates default JSON files in a given folder to get start working on a new project.

```shell script
acc create PATH/TO/FOLDER/
```


### query

Search for certain elements.


### records

Exports expenses and invoices as an annotated business records for taxes and activation.


### validate

Check your data.


## Workflows

### Workflow with Bimpf

1. Import from Bimpf
	- Dump Bimpf data as JSON.
	- Validate dump with `acc bimpf validate -i dump.json` and resolve the problems in Bimpf, then re-export the data.
	- Import the date with `acc bimpf import -i dump.json -n /path/to/nc-folder`.
	- Validate the import with `acc validate -i acc.yaml`.
2. Open `acc.yaml` and complete the configuration. Do not forget to...
	- ...enter your company details (name, street, etc.).
	- ...change the journal accounts to your needs.
	- ...enter all account aliases.
	- ...define known expense categories and their appropriate journal account (you also can add categories while complete expenses).
3. Complete...
	- ...expenses with `acc complete expenses`.
	- ...invoices with `acc complete invoices`.
4. Validate the result with `acc validate -i acc.yaml`.
5. Bank statement:
	- Import the bank statement with `acc camt -i acc.yaml -s camt.xml`.
	- Complete with `acc complete expenses -i acc.yaml`.
	- Validate with `acc validate -i acc.yaml`.
6. Repopulate expenses and invoices with `acc complete repopulate -i acc.yaml`.
7. Complete transactions with `acc complete transactions -i acc.yaml`
9. Export the journal with `acc ledger -i acc.yaml -o ledger.journal -y 20XX` while replacing «20XX» with the year you want to filter for. Take a close look at the warnings and change the files if needed. 
10. Manually correct the journal:
	- Manually book all opening entries.
	- Check **all** generated bookings.
	- Amortize the positions you're legally obliged/allowed to.
	- Book all necessary reserve assets.
	- Finalize your journal.


### Simple theater project

1. Prepare:
	- Collect all receipts of expenses in single folder (e.g. `raw/`).
2. Create new project `acc new`.
3. Add all people which payed some expense in advance for the project as employee `acc add employee -i acc.yaml`.
4. Add all expenses with `acc add expense -i acc.yaml -a raw/receipt.pdf`.
5. Import or create all bank transactions for the project. Pay attention you only add transactions which have some connection to the theater production.
6. Complete the transactions with `acc complete transactions -i acc.yaml -d` using the document-only mode (only associated documents will be completed).
7. Repopulate the expenses with `acc complete repopulate -i acc.yaml`.
