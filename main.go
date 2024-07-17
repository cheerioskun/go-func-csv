package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/fzipp/gocyclo"
	"github.com/gocarina/gocsv"
)

type FuncDetails struct {
	PkgName      string `csv:"Package"`
	FunctionName string `csv:"Function Name"`
	// LoC          int    `csv:"LoC"`
	Complexity int    `csv:"Complexity"`
	Location   string `csv:"Location"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run main.go <package-directory>")
		return
	}

	if err := execute(os.Args[1:]); err != nil {
		log.Printf("execution failed: %v", err)
	}
}

// execute runs the logic recursively for all go files under path
func execute(paths []string) error {

	var records []*FuncDetails
	fmt.Printf("Analyzing files using gocyclo...")
	stats := gocyclo.Analyze(paths, regexp.MustCompile(".*test.go|asset.go")).SortAndFilter(-1, 1)
	fmt.Printf("Reformatting to CSV...")
	for _, stat := range stats {
		records = append(records, &FuncDetails{
			FunctionName: stat.FuncName,
			PkgName:      stat.PkgName,
			Location:     regexp.MustCompile(`.*/msp-controller/msp/`).ReplaceAllString(fmt.Sprintf("%s", stat.Pos.Filename), ""),
			Complexity:   stat.Complexity,
		})
	}

	file, err := os.OpenFile("out.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open output csv: %v", err)
	}
	defer file.Close()
	fmt.Printf("Writing to output file...")
	if err := gocsv.MarshalFile(records, file); err != nil {
		return fmt.Errorf("could not write to csv: %v", err)
	}
	return nil
}
