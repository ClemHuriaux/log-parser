package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func fillMissingParts(parts []string, expectedLength int) []string {
	if len(parts) < expectedLength {
		for i := len(parts); i < expectedLength; i++ {
			parts = append(parts, "")
		}
	}
	return parts
}

func readLines(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func parseLogs(source string, destination string, columnsFile string, verbose bool, separator string) error {
	files, err := filepath.Glob(filepath.Join(source, "*.*"))
	if err != nil {
		return err
	}

	columns, err := readLines(columnsFile)
	if err != nil {
		return err
	}

	var data[][]string

	for i, file := range files {
		if verbose {
			fmt.Printf("Parsing file %d/%d\n", i+1, len(files))
		}

		lines, err := readLines(file)
		if err != nil {
			return err
		}

		for _, line := range lines {
			parts := strings.Split(line, separator)
			parts = fillMissingParts(parts, len(columns))
			data = append(data, parts)
			
		}

		if verbose {
			progressPercentage := (float64(i+1) / float64(len(files))) * 100
			fmt.Printf("Progress: %.2f%%\n", progressPercentage)
		}
	}
	file, err := os.Create(filepath.Join(destination, "results.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(columns)
	if err != nil {
		return err
	}

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

func main() {
	var verbose bool
	var separator string
	flag.BoolVar(&verbose, "v", false, "Show state of files parsing")
	flag.StringVar(&separator, "s", " ", "Specify the separator of the files")
	flag.Parse()

	columnsFile := flag.Arg(0)
	destination := flag.Arg(1)
	source := flag.Arg(2)

	startTime := time.Now()
	err := parseLogs(source, destination, columnsFile, verbose, separator)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	endTime := time.Now()

	fmt.Printf("Execution time: %v\n", endTime.Sub(startTime))
}
