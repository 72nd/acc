# acc

## Usage and Functions

### create

Creates default JSON files in a given folder to get start working on a new project.

```shell script
acc create PATH/TO/FOLDER/
```

### bimpf
[Bimpf](https://gitlab.com/solutionsbuero/bimpf)

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