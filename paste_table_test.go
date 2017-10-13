package main

import (
	"strings"
	"testing"
)

func equal(t *testing.T, in string, want string, index int) {
	if strings.TrimSpace(in) != want {
		t.Errorf("Test #%d: %q != %q", index, in, want)
	}
}

func isFalse(t *testing.T, b bool) {
	if b == true {
		t.Errorf("Expected false")
	}
}

func isTrue(t *testing.T, b bool) {
	if b == false {
		t.Errorf("Expected true")
	}
}

// describe('paste table', function() {
//   describe('should handle mysql cli output', function() {
//     it('with a single column, single row', function() {
//       var data = '+---------+\n| 117700 |\n+---------+\n';
//       var program = mockProgram();
//       should(pasteTable.commandLine(data.split('\n'), program)).equal('117700');
//     });
func TestParseTableWithDefault(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		options PasteTableOptions
	}{
		{
			"+---------+\n| 117700 |\n+---------+\n",
			"117700",
			PasteTableOptions{},
		},
		// by default we get only the first column
		{
			"+---------+\n| 117700 | 22 |\n+---------+\n",
			"117700",
			PasteTableOptions{},
		},
		// switch the column
		{
			"+---------+\n| 117700 | 22 |\n+---------+\n",
			"22",
			PasteTableOptions{column: 1},
		},
		// all columns
		{
			"+---------+\n| 117700 | 22 |\n+---------+\n",
			"117700,22",
			PasteTableOptions{column: -1},
		},
		// multiple rows, single column
		{
			"+---------+\n| 117700 |\n| 22 |\n+---------+\n",
			"117700,22",
			PasteTableOptions{},
		},
		// multiple rows, single column w newLinePerColumn option
		{
			"+---------+\n| 117700 |\n| 22 |\n+---------+\n",
			"117700\n22",
			PasteTableOptions{newLinePerColumn: true},
		},
	}

	for index, testCase := range cases {
		equal(t, PasteTable(testCase.in, testCase.options), testCase.want, index)
	}
}

func TestParseTableWithCsv(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		options PasteTableOptions
	}{
		{
			"117700",
			"117700",
			PasteTableOptions{delimiter: ","},
		},
		// output one column per record
		{
			"117700,22",
			"117700",
			PasteTableOptions{delimiter: ","},
		},
		// output one column per record, one record per line
		{
			"117700,22\n128811,33",
			"117700,128811",
			PasteTableOptions{delimiter: ","},
		},
		// output all columns, one record per line
		{
			"117700,22",
			"117700,22",
			PasteTableOptions{delimiter: ",", column: -1},
		},
		// output all columns, one record per line
		{
			"117700,22\n128811,33",
			"117700,22\n128811,33",
			PasteTableOptions{delimiter: ",", column: -1},
		},
		{
			"117700,22\n128811,33",
			"117700\n128811",
			PasteTableOptions{delimiter: ",", newLinePerColumn: true},
		},
		{
			"117700,22",
			"117700\n22",
			PasteTableOptions{delimiter: ",", column: -1, newLinePerColumn: true},
		},
		{
			"117700\t22",
			"117700",
			PasteTableOptions{delimiter: "\t"},
		},
		{
			"117700\t22",
			"22",
			PasteTableOptions{delimiter: "\t", column: 1},
		},
		{
			"117700\t22",
			"22",
			PasteTableOptions{delimiter: "\\t", column: 1},
		},
	}

	for index, testCase := range cases {
		equal(t, PasteTable(testCase.in, testCase.options), testCase.want, index)
	}
}

func TestNoOptions(t *testing.T) {
	opts := PasteTableOptions{}
	isTrue(t, noOptions(opts))

	opts = PasteTableOptions{column: 1}
	isFalse(t, noOptions(opts))
}
