package markdown

import (
	"regexp"
)

const (
	allEOLChars = "\\r\\n\\v\\f\\p{Zl}\\p{Zp}"
)

var (
	// All constructs are restricted to operating on a single line

	// EOLs
	br = regexp.MustCompile("(?:\\r\\n|\\n\\r|\\n|\\r|\\p{Zl}|\\p{Zp}|\\v|\\f)")
	// BOL, 1 to 6 #s, a space, some text, EOL.
	header = regexp.MustCompile("(?m)^([#]{1,6})[ ]([^" + allEOLChars + "]+?)$")
	// BOL, optional text, 1 star, optional char that's not a star, required text, 1 star, optional char that's not a star, optional text, EOL. Requires >= 2 characters inside.
	italics = regexp.MustCompile("(?m)^[^" + allEOLChars + "*]*?[*]([^*" + allEOLChars + "][^" + allEOLChars + "]*[^*" + allEOLChars + "])[*][^" + allEOLChars + "*]*?$")
	// Duplicate of italics except with two stars.
	bold = regexp.MustCompile("(?m)^[^" + allEOLChars + "*]*?[*]{2}([^*" + allEOLChars + "][^" + allEOLChars + "]*[^*" + allEOLChars + "])[*]{2}[^" + allEOLChars + "*]*?$")
	// BOL, >, a space, optional text, EOL.
	blockquote = regexp.MustCompile("(?m)^>[ ]([^" + allEOLChars + "]*?)$")
	// Plug this into an online regexp explainer and you'll see why
	link = regexp.MustCompile(`(?m)(?:^|^[^` + allEOLChars + `]*?[^!])\[([^` + allEOLChars + `]+?)\]\(([^` + allEOLChars + `]*?)\)[^` + allEOLChars + `]*?$`)
	// Plug this into an online regexp explainer and you'll see why
	img = regexp.MustCompile(`(?m)^[^` + allEOLChars + `]*?!\[([^` + allEOLChars + `]+?)\]\(([^` + allEOLChars + `]*?)\)[^` + allEOLChars + `]*?$`)
)

func MarkdownToHtml(md string) string {
	html := md
	matched := 0
	for {

		if matched == 0 {
			break
		}
	}
	html, _ = replaceGeneric(md, br, actionBr)
	return html
}

func replaceGeneric(md string, exp *regexp.Regexp, action func(string, []int) string) (string, bool) {
	actionHappened := false
	matches := exp.FindAllStringSubmatchIndex(md, -1)
	reverse(matches)
	for _, match := range matches {
		md = action(md, match)
		actionHappened = true
	}
	return md, actionHappened
}

func actionBr(md string, match []int) string {
	temp := "<br>" + md[match[0]:match[1]]
	if len(md) > match[1] {
		md = md[:match[0]] + temp + md[match[1]:]
	} else {
		md = md[:match[0]] + temp
	}
	return md
}

func reverse(a [][]int) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}
