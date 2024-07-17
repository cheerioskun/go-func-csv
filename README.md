# go-func-csv
## Generate a management friendly view of your Go codebase
Given a Go repository, this tool generates a CSV formatted text file where all the functions in the repository are listed with multiple other fields like
- package name
- cyclomatic complexity
- Location in codebase

These were the fields of interest selected for a specific use case of writing unit tests in retrospect for an existing codebase. Your needs might require you to customize
this part.

## Usage
```
go build . -o ftocsv
./ftocsv <path to package or repo>
# out.csv file is created
```

## Sample
```
# Running the tool on charmbracelet/gum
Package,Function Name,Complexity,Location
filter,(model).Update,24,../../charmbracelet/gum/filter/filter.go
filter,(Options).Run,22,../../charmbracelet/gum/filter/command.go
table,(Options).Run,17,../../charmbracelet/gum/table/command.go
choose,(Options).Run,16,../../charmbracelet/gum/choose/command.go
filter,(model).View,15,../../charmbracelet/gum/filter/filter.go
log,(Options).Run,13,../../charmbracelet/gum/log/command.go
```

