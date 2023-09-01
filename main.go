package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
)

type FuncDetails struct {
	FunctionName string `csv:"Function Name"`
	LoC          int    `csv:"LoC"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run main.go <package-directory>")
		return
	}

	if err := execute(os.Args[1]); err != nil {
		log.Printf("execution failed: %v", err)
	}
}

// execute runs the logic recursively for all go files under path
func execute(basepath string) error {

	fset := token.NewFileSet()
	var records []*FuncDetails

	if fh, err := os.Open(basepath); err != nil {
		log.Printf("could not open the input path: %v", err)
	} else if info, _ := fh.Stat(); info.IsDir() {
		fsys := os.DirFS(basepath)
		records = processDirectory(fsys, fset, basepath)
	} else {
		records = processFile(basepath, fset)
	}

	file, err := os.OpenFile("out.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open output csv: %v", err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(records, file); err != nil {
		return fmt.Errorf("could not write to csv: %v", err)
	}
	return nil
}

// processDirectory walks through the directory and processes Go files
func processDirectory(fsys fs.FS, fset *token.FileSet, inputPath string) []*FuncDetails {
	records := []*FuncDetails{}
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			records = append(records, processFile(filepath.Join(inputPath, path), fset)...)
		}
		return nil
	})
	return records
}

// processFile parses and inspects the AST of a Go file
func processFile(path string, fset *token.FileSet) []*FuncDetails {
	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Printf("could not open file %s: %v", path, err)
		return nil
	}
	records := []*FuncDetails{}
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			loc := fset.Position(funcDecl.Body.Rbrace).Line - fset.Position(funcDecl.Body.Lbrace).Line - 1
			records = append(records, &FuncDetails{
				FunctionName: funcDecl.Name.Name,
				LoC:          loc,
			})
		}
		return true
	})

	return records
}
