package markdown

import (
	"testing"
	"fmt"
)

type testString struct {
	string
	expected interface{}
}

func TestHeader(t *testing.T) {
	testStrs := [...]testString{
		{"# h1", true},
		{"## h2", true},
		{"### h3", true},
		{"#### h4", true},
		{"##### h5", true},
		{"###### h6", true},
		{"## h1 ##", true},
		{" ## h2", false},
		{"", false},
		{"\n# h1\n", true},
		{"\r\n# h1\n\r", true},
	}
	for _, testStr := range testStrs {
		if val := header.MatchString(testStr.string); val != testStr.expected {
			t.Error("For", testStr.string, "expected", testStr.expected, "got", val)
			break
		}
	}

	testStrs = [...]testString{
		{"# h1", [][]string{{"# h1", "#", "h1"}}},
		{"## h2", [][]string{{"## h2", "##", "h2"}}},
		{"### h3", [][]string{{"### h3", "###", "h3"}}},
		{"#### h4", [][]string{{"#### h4", "####", "h4"}}},
		{"##### h5", [][]string{{"##### h5", "#####", "h5"}}},
		{"###### h6", [][]string{{"###### h6", "######", "h6"}}},
		{"## h1 ##", [][]string{{"## h1 ##", "##", "h1 ##"}}},
		{" ## h2", [][]string{}},
		{"", [][]string{}},
		{"\n# h1\n", [][]string{{"# h1", "#", "h1"}}},
		{"\r\n# h1\n\r", [][]string{{"# h1", "#", "h1"}}},
	}
	for _, testStr := range testStrs {
		if val := header.FindAllStringSubmatch(testStr.string, -1); !strSliceEqual(val, testStr.expected.([][]string)) {
			t.Error("For", testStr.string, "expected", fmt.Sprintf("%q", testStr.expected), "got", fmt.Sprintf("%q", val))
			break
		}
	}
}

func strSliceEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j, _ := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
