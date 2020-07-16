# acc

First planned as a simple tool chain for easing the work with the plain text accounting software [hledger](https://hledger.org/), acc evolved into some sort of plain-text ERP system. It's capable of keeping track of your customers, employees, expenses and invoices as well as importing bank statements (via ISO 20022 camt). Acc can generate most of your hledger account records based on this data. Also it's possible to export all business records (expenses etc.) with their respective receipt as a PDF for further archiving. There is also a quit powerful query functionality for easy finding of records. For adding new data, a interactive prompt can be used (you also can just edit the YAML files by yourself).

_Current status:_ The project is under active use and development in the [Genossenschaft Solutionsbüro](https://buero.io) and Acc evolves along our needs. Almost all core features are implemented and we use this software in your everyday work. The biggest throwback in the moment are the missing tests and a bad code documentation (I'm working on it). Also a manual (or at least a detailed and long README) is missing at this point.

## Todo

- [ ] First Run with 2019 Taxes
- [ ] Documentation in Readme done
- [x] Filter function
- [ ] Project Mode _(almost done)_
- [x] Shorts for commands 
- [ ] Code Documentation
- [ ] Misc Documents in complete Transactions
- [ ] Amounts with Amount type
- [ ] pain.001 for payment generation
	- [ ] Add IBAN to employee (data field, interactive add, assisted completion, validation)
	- [x] Rewrite `sba-pay` as library for go
	- [ ] Add functionality to iso20022 package



## Installation

```shell script
sudo apt install wmctrl
```


## Basic Concepts

To start with Acc, it's important to understand some fundamental concepts and ideas which shape the handling of the software.


### General Idea

The basic idea of Acc is to collect data about your business and then generating a number of different outputs from these.

![Example flow of data](misc/flow-of-data.svg)


### Types of data/records


### Modes


## Usage and Functions

### add

Add new elements (customer, employee, expense, expense-category, invoice or transaction) to your acc project. If you don't want to use the interactive prompt, use the `--default` flag. Some of the elements contain paths to files by using the `--asset` flag you can specify this paths in advance and thus use the tab-completion of your shell.

```shell script
acc add customer -i acc.yaml
acc add invoice -i acc.yaml --asset /path/to/sent-invoice.pdf
```


### bimpf

Hint: You can dump the data of [Bimpf](https://gitlab.com/solutionsbuero/bimpf) with the following command on the server: 

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


## camt

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
- `--retain-focus` _Experimentally and hacky feature!_ Tries to regain focus of the terminal acc runs in. Will only work on certain Linux installations if `wmctrl` is installed.

The `acc complete repopulate` on the other hand can be used to link expenses and invoices to transactions which already are linked to the expense/invoice.


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
	- 


### Simple theater project

1. Prepare:
	- Collect all receipts of expenses in single folder (e.g. `raw/`).
2. Create new project `acc new`.
3. Add all people which payed some expense in advance for the project as employee `acc add employee -i acc.yaml`.
4. Add all expenses with `acc add expense -i acc.yaml -a raw/receipt.pdf`.
5. Import or create all bank transactions for the project. Pay attention you only add transactions which have some connection to the theater production.
6. Complete the transactions with `acc complete transactions -i acc.yaml -d` using the document-only mode (only associated documents will be completed).
7. Repopulate the expenses with `acc complete repopulate -i acc.yaml`.
