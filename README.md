# acc

## Usage and Functions

### create

Creates default JSON files in a given folder to get start working on a new project.

```shell script
acc create PATH/TO/FOLDER/
```

### bimpf

Hint: You can dump the data of [Bimpf](https://gitlab.com/solutionsbuero/bimpf) with the following command on the server: 

```shell 
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

### documents

### hledger

## Workflow

1. Dump Bimpf data
2. Validata dump: `acc bimpf validate -i dump.json`