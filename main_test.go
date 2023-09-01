package main

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	sampleFilePath := "testdata/sample.go"

	if err := execute(sampleFilePath); err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Open the generated CSV file
	file, err := os.Open("out.csv")
	if err != nil {
		t.Fatalf("Could not open output CSV file: %v", err)
	}
	defer file.Close()

	// Parse CSV and validate records
	csvReader := csv.NewReader(file)
	records := []*FuncDetails{}
	if err := gocsv.UnmarshalCSV(csvReader, &records); err != nil {
		t.Fatalf("Could not unmarshal CSV: %v", err)
	}

	// Validate results
	assert.Equal(t, 2, len(records), "Incorrect number of records")
	assert.Equal(t, "HelloWorld", records[0].FunctionName)
	assert.Equal(t, 1, records[0].LoC)
	assert.Equal(t, "Add", records[1].FunctionName)
	assert.Equal(t, 1, records[1].LoC)
}
