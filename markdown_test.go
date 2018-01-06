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
		if val := headerExp.MatchString(testStr.string); val != testStr.expected {
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
		if val := headerExp.FindAllStringSubmatch(testStr.string, -1); !strSliceEqual(val, testStr.expected.([][]string)) {
			t.Error("For", testStr.string, "expected", fmt.Sprintf("%q", testStr.expected), "got", fmt.Sprintf("%q", val))
			break
		}
	}
}

func TestItalics(t *testing.T) {
	testStrs := []testString{
		{"*italics*", true},
		{"bblahblahblahblah*italics      yo!*blah", true},
		{"** italics **", false},
		{"\n* **italics** *\n", true},
		{"\r\n*italics*\n\r", true},
		{"\r\n*italics\n*\r", false},
		{"\r\n** *italics* **\n\r", false},
	}
	for _, testStr := range testStrs {
		if val := italicsExp.MatchString(testStr.string); val != testStr.expected {
			t.Error("For", testStr.string, "expected", fmt.Sprintf("%q", testStr.expected), "got", fmt.Sprintf("%q", val))
			break
		}
	}
	//
	//testStrs = []testString{
	//	{"# h1", [][]string{{"# h1", "#", "h1"}}},
	//	{"## h2", [][]string{{"## h2", "##", "h2"}}},
	//	{"### h3", [][]string{{"### h3", "###", "h3"}}},
	//	{"#### h4", [][]string{{"#### h4", "####", "h4"}}},
	//	{"##### h5", [][]string{{"##### h5", "#####", "h5"}}},
	//	{"###### h6", [][]string{{"###### h6", "######", "h6"}}},
	//	{"## h1 ##", [][]string{{"## h1 ##", "##", "h1 ##"}}},
	//	{" ## h2", [][]string{}},
	//	{"", [][]string{}},
	//	{"\n# h1\n", [][]string{{"# h1", "#", "h1"}}},
	//	{"\r\n# h1\n\r", [][]string{{"# h1", "#", "h1"}}},
	//}
	//for _, testStr := range testStrs {
	//	if val := italics.FindAllStringSubmatch(testStr.string, -1); !strSliceEqual(val, testStr.expected.([][]string)) {
	//		t.Error("For", testStr.string, "expected", fmt.Sprintf("%q", testStr.expected), "got", fmt.Sprintf("%q", val))
	//		break
	//	}
	//}
}

func TestBlockquote(t *testing.T) {
	testStrs := []testString{
		{"> woot", true},
		{">", false},
		{"> ", true},
		{"blah blah > sdfsdf", false},
	}
	for _, testStr := range testStrs {
		if val := blockquoteExp.MatchString(testStr.string); val != testStr.expected {
			t.Error("For", testStr.string, "expected", fmt.Sprintf("%q", testStr.expected), "got", fmt.Sprintf("%q", val))
			break
		}
	}
}

func TestLink(t *testing.T) {
	testStrs := []testString{
		{"It's crazy! [This](link) gives you magical powers!", true},
		{"It's crazy! ![This](image) gives you magical powers!", false},
		{"![This](image) shows you one easy trick to make money fast!", false},
		{"![This](image) shows you one easy trick to make money fast! [This](link) does too!", true},
		{"![This](image) shows you one easy trick to make money fast! [Th!s](l!nk) does too!", true},
	}
	for _, testStr := range testStrs {
		if val := linkExp.MatchString(testStr.string); val != testStr.expected {
			t.Error("For", fmt.Sprintf("%q", testStr.string), "expected", fmt.Sprintf("%v", testStr.expected), "got", fmt.Sprintf("%v", val))
			break
		}
	}
}

func TestImage(t *testing.T) {
	testStrs := []testString{
		{"It's crazy! [This](link) gives you magical powers!", false},
		{"It's crazy! ![This](image) gives you magical powers!", true},
		{"![This](image) shows you one easy trick to make money fast!", true},
		{"[This](link) shows you one easy trick to make money fast! ![This](image) does too!", true},
		{"[This](link) shows you one easy trick to make money fast! ![Th!s](l!nk) does too!", true},
		{"[This](link) shows you one easy trick to make money fast!![Th!s](l!nk) does too!", true},
	}
	for _, testStr := range testStrs {
		if val := imgExp.MatchString(testStr.string); val != testStr.expected {
			t.Error("For", fmt.Sprintf("%q", testStr.string), "expected", fmt.Sprintf("%v", testStr.expected), "got", fmt.Sprintf("%v", val))
			break
		}
	}
}

func TestMarkdownToHtml(t *testing.T) {
	testStrs := [...]testString{
		{"# h1", "<h1>h1</h1>"},
		{"# h1\n# h1", "<h1>h1</h1><br>\n<h1>h1</h1>"},
		{"## h2", "<h2>h2</h2>"},
		{"### h3", "<h3>h3</h3>"},
		{"#### h4", "<h4>h4</h4>"},
		{"You must put an *emphasis* on your work.", "You must put an <i>emphasis</i> on your work."},
		{"You must be **bold** in the face of danger.", "You must be <b>bold</b> in the face of danger."},
		{"You must be **bold** in the face of danger, *yet* know when to run and hide.", "You must be <b>bold</b> in the face of danger, <i>yet</i> know when to run and hide."},
		{"[This](link) shows you one easy trick to make money fast!![Th!s](image!) does too!", "<a href='link'>This</a> shows you one easy trick to make money fast!<img src='image!' alt='Th!s'> does too!"},
	}
	for _, testStr := range testStrs {
		if val := MarkdownToHtml(testStr.string); val != testStr.expected {
			t.Error("For", testStr.string, "\nexpected", fmt.Sprintf("%q\n", testStr.expected), "got", fmt.Sprintf("%q\n", val))
			break
		}
	}
}

func strSliceEqual(a, b [][]string) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		if len(a) == len(b) { // nil slice has len 0
			return true
		}
		return false
	}
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
