# acc


## Todo

- [ ] First Run with 2019 Taxes
- [ ] Documentation in Readme done
- [ ] Filter function
- [ ] Project Mode
- [ ] Shorts for commands 
- [ ] Code Documentation

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


## Workflow

1. Dump Bimpf data
2. Validate dump with `acc bimpf validate -i dump.json` and resolve the problems in Bimpf, then re-export the data.
3. Import the date with `acc bimpf import -i dump.json -n /path/to/nc-folder`.
4. Validate the import with `acc validate -i acc.yaml`.
5. Complete...
	- ...expenses with `acc complete expenses`.
	- ...invoices with `acc complete invoices`.
6. Validate the result with `acc validate -i acc.yaml`.
7. Bank statement:
	- Import the bank statement with `acc camt -i acc.yaml -s camt.xml`.
	- Complete with `acc complete expenses`.
	- Validate with `acc validate -i acc.yaml`.
8. Repopulate expenses and invoices with `acc complete repopulate -i acc.yaml`.
