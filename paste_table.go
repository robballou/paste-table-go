package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Column           int    `cli:"column" usage:"Column number to select" dft:"0"`
	Delimiter        string `cli:"delimiter" usage:"Delimiter used in the input for separating columns" dft:""`
	NewLinePerColumn bool   `cli:"newline" usage:"Put single columns on a new line" dft:"false"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		// ctx.String("%d, %s, %d, %s\n", argv.Basic, argv.Env, argv.Expr, argv.DevDir)

		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if fi.Size() == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: No data to read. Pipe in content to process.\n")
			ctx.WriteUsage()
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Error reading stdin")
		}
		in := string(b)

		opts := PasteTableOptions{column: argv.Column, newLinePerColumn: argv.NewLinePerColumn, delimiter: argv.Delimiter}
		ctx.String(PasteTable(in, opts))
		return nil
	})
}

// PasteTableOptions structure around CLI options
type PasteTableOptions struct {
	column           int
	delimiter        string
	newLinePerColumn bool
}

// PasteTable converts the input string into the desired output
//
// Figuring out what to do with the input can be done in a couple ways:
// - check the options
func PasteTable(input string, options PasteTableOptions) string {
	// default
	if noDelimiter(options) {
		return ParseTable(input, options)
	}

	return ParseDelimited(input, options)
}

// ParseTable parses the input from a MySQL-table string
func ParseTable(input string, options PasteTableOptions) string {
	var output []string

	// pattern to skip out on +----------+ lines
	tableLinesPattern := regexp.MustCompile(`^\+(-)+\+$`)

	// pattern to match the values in a table output, | value1 | value2 |
	tableCellsPattern := regexp.MustCompile(`\|?\s*([^|]+)\s*\|`)

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if tableLinesPattern.MatchString(line) {
			continue
		}
		lineData := tableCellsPattern.FindAllStringSubmatch(line, -1)
		if lineData == nil {
			continue
		}

		var thisLine string

		switch {

		// return a specific column
		case options.column >= 0 && lineData[options.column] != nil:
			thisLine = strings.TrimSpace(lineData[options.column][1])

		// return all columns as csv
		case options.column == -1:
			// build our column data...
			var thisLineColumns []string
			for _, column := range lineData {
				thisLineColumns = append(thisLineColumns, strings.TrimSpace(column[1]))
			}

			// use the csv library to write the CSV (in case values have commas)

			thisLine = sliceToCsv(thisLineColumns)
		}

		output = append(output, thisLine)
	}

	if options.column >= 0 && options.newLinePerColumn == false {
		return fmt.Sprintf("%s\n", sliceToCsv(output))
	}

	return fmt.Sprintf("%s\n", strings.Join(output, "\n"))
}

func ParseDelimited(input string, options PasteTableOptions) string {
	var output []string

	r := csv.NewReader(strings.NewReader(input))

	// handle cases where escaped characters are passed in
	switch options.delimiter {
	case "\\t":
		options.delimiter = "\t"
	case "\\n":
		options.delimiter = "\n"
	}

	delimiter := []rune(options.delimiter)
	r.Comma = delimiter[0]

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		lineLength := len(record)
		if lineLength == 0 {
			continue
		}

		switch {
		// requested a single column
		case options.column >= 0 && lineLength > options.column:
			output = append(output, strings.TrimSpace(record[options.column]))
		case options.column == -1 && options.newLinePerColumn == false:
			output = append(output, sliceToCsv(normalize(record)))
		case options.column == -1 && options.newLinePerColumn:
			output = append(output, strings.Join(normalize(record), "\n"))
		}
	}

	if options.column >= 0 && options.newLinePerColumn == false {
		return fmt.Sprintf("%s\n", sliceToCsv(output))
	}

	return fmt.Sprintf("%s\n", strings.Join(output, "\n"))
}

func normalize(input []string) []string {
	var output []string

	for _, in := range input {
		output = append(output, strings.TrimSpace(in))
	}

	return output
}

func sliceToCsv(s []string) string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Write(s)
	w.Flush()
	return strings.TrimSpace(buf.String())
}

func noDelimiter(options PasteTableOptions) bool {
	return options.delimiter == ""
}

// noOptions check if the options struct has no options set
func noOptions(options PasteTableOptions) bool {

	if options.column == 0 && noDelimiter(options) {
		return true
	}

	return false
}
