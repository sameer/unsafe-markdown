package markdown

import (
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	allEOLChars = "\\r\\n\\v\\f\\p{Zl}\\p{Zp}"
)

var (
	// All constructs are restricted to operating on a single line

	// I tried to explain some of it but I had to make too many changes, so the explanations might not match exactly.

	// EOLs
//	brExp = regexp.MustCompile("(?m)^(?:\\r\\n|\\n\\r|\\n|\\r|\\p{Zl}|\\p{Zp}|\\v|\\f)$")
    brExp = regexp.MustCompile("(?m)^$")
	// BOL, 1 to 6 #s, a space, some text, EOL.
	headerExp = regexp.MustCompile("(?m)^([#]{1,6})[ ]([^" + allEOLChars + "]+?)$")
	// BOL, optional text, 1 star, optional char that's not a star, required text, 1 star, optional char that's not a star, optional text, EOL. Requires >= 2 characters inside.
	italicsExp = regexp.MustCompile(`(?m)(?:^[^\r\n\v\f\p{Zl}\p{Zp}*]|^[^\r\n\v\f\p{Zl}\p{Zp}]*?[^*])[*]{1}([^\r\n\v\f\p{Zl}\p{Zp}*]+?)[*]{1}(?:$|[^\r\n\v\f\p{Zl}\p{Zp}*]$|[^*][^\r\n\v\f\p{Zl}\p{Zp}]*?$)`)
	// Duplicate of italics except with backticks.
	codeExp = regexp.MustCompile("(?m)(?:^[^\\r\\n\\v\\f\\p{Zl}\\p{Zp}`]|^[^\\r\\n\\v\\f\\p{Zl}\\p{Zp}]*?[^`])[`]{1}([^\\r\\n\\v\\f\\p{Zl}\\p{Zp}`]+?)[`]{1}(?:$|[^\\r\\n\\v\\f\\p{Zl}\\p{Zp}`]$|[^`][^\\r\\n\\v\\f\\p{Zl}\\p{Zp}]*?$)")
	// Duplicate of italics except with two stars.
	boldExp = regexp.MustCompile(`(?m)(?:^[^\r\n\v\f\p{Zl}\p{Zp}*]|^[^\r\n\v\f\p{Zl}\p{Zp}]*?[^*])[*]{2}([^\r\n\v\f\p{Zl}\p{Zp}*]+?)[*]{2}(?:$|[^\r\n\v\f\p{Zl}\p{Zp}*]$|[^*][^\r\n\v\f\p{Zl}\p{Zp}]*?$)`)
	// Duplicate of italics except with two tildes.
	strikethroughExp = regexp.MustCompile(`(?m)(?:^[^\r\n\v\f\p{Zl}\p{Zp}~]|^[^\r\n\v\f\p{Zl}\p{Zp}]*?[^~])[~]{2}([^\r\n\v\f\p{Zl}\p{Zp}~]+?)[~]{2}(?:$|[^\r\n\v\f\p{Zl}\p{Zp}~]$|[^~][^\r\n\v\f\p{Zl}\p{Zp}]*?$)`)
	// BOL, >, a space, optional text, EOL.
	blockquoteExp = regexp.MustCompile("(?m)(?:^>)[ ]([^" + allEOLChars + "]*?)$")
	// Plug this into an online regexp explainer and you'll see why, too complex
	linkExp = regexp.MustCompile(`(?m)(?:^|^[^` + allEOLChars + `]*?[^!])\[([^` + allEOLChars + `]+?)\]\(([^` + allEOLChars + `]*?)\)[^` + allEOLChars + `]*?$`)
	// Plug this into an online regexp explainer and you'll see why, too complex
	imgExp = regexp.MustCompile(`(?m)^[^` + allEOLChars + `]*?!\[([^` + allEOLChars + `]+?)\]\(([^` + allEOLChars + `]*?)\)[^` + allEOLChars + `]*?$`)
)

// Convert a string of pre-escaped Markdown text into an HTML string.
func MarkdownToHtmlString(md string) string {
	html := md

	for {
		actionCount := 0

		html = replaceGeneric(html, headerExp, actionHeader, &actionCount)

		html = replaceGeneric(actionGeneric(html, italicsExp, &actionCount))
		html = replaceGeneric(actionGeneric(html, boldExp, &actionCount))
		html = replaceGeneric(actionGeneric(html, blockquoteExp, &actionCount))
		html = replaceGeneric(actionGeneric(html, codeExp, &actionCount))
		html = replaceGeneric(actionGeneric(html, strikethroughExp, &actionCount))

		html = replaceGeneric(actionLinklike(html, linkExp, &actionCount))
		html = replaceGeneric(actionLinklike(html, imgExp, &actionCount))

		if actionCount == 0 {
			break
		}
	}
	// This can only be done once because we maintain the newlines for readability,
	// if it was in the loop it would infinitely replace.
	html = replaceGeneric(html, brExp, actionBr, nil)
	return html
}

// Convert a byte array of pre-escaped Markdown text into an HTML string.
func MarkdownToHtmlByte(md []byte) []byte {
	// TODO: modify this so allocations can be reduced
	return []byte(MarkdownToHtmlString(string(md)))
}

// Convert Markdown text from an io.Reader and write it to an io.Writer.
func MarkdownToHtmlIO(r io.Reader, w io.Writer) error {
	html, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	w.Write(MarkdownToHtmlByte(html))
	return nil
}

func replaceGeneric(md string, exp *regexp.Regexp, action func(string, []int) string, actionCount *int) string {
	actionHappened := false
	matches := exp.FindAllStringSubmatchIndex(md, -1)
	reverse(matches)
	for _, match := range matches {
		md = action(md, match)
		actionHappened = true
	}
	if actionHappened && actionCount != nil {
		*actionCount++
	}
	return md
}

func actionHeader(md string, match []int) string {
	// Generic match 0 1, 1st group 2 3, 2nd group 4 5
	hType := strconv.Itoa(match[3] - match[2])
	openTag := "<h" + hType + ">"
	closeTag := "</h" + hType + ">"
	temp := openTag + md[match[4]:match[5]] + closeTag
	if len(md) > match[1] {
		md = md[:match[0]] + temp + md[match[1]:]
	} else {
		md = md[:match[0]] + temp
	}
	return md
}

func actionGeneric(html string, exp *regexp.Regexp, actionCount *int) (string, *regexp.Regexp, func(string, []int) string, *int) {
	var openTag, closeTag string
	if exp == italicsExp {
		openTag, closeTag = "<i>", "</i>"
	} else if exp == boldExp {
		openTag, closeTag = "<b>", "</b>"
	} else if exp == blockquoteExp {
		openTag, closeTag = "<blockquote>", "</blockquote>"
	} else if exp == codeExp {
		openTag, closeTag = "<code>", "</code>"
	} else if exp == strikethroughExp {
		openTag, closeTag = "<s>", "</s>"
	} else {
		panic("Unknown regex expression provided!")
	}
	return html, exp, func(md string, match []int) string {
		// Whole line match 0 1, 1st group 2 3
		temp := openTag + md[match[2]:match[3]] + closeTag
		leftOffset, rightOffset := 1, 1
		if exp == boldExp {
			// two stars
			leftOffset, rightOffset = 2, 2
		} else if exp == blockquoteExp {
			// blockquote shouldn't take on any functionality like this
			leftOffset, rightOffset = 2, 0
		}
		if len(md) > match[3]+rightOffset {
			temp += md[match[3]+rightOffset:]
		}
		if match[2]-leftOffset > 0 {
			temp = md[:match[2]-leftOffset] + temp
		}
		md = temp
		return md
	}, actionCount
}

func actionLinklike(html string, exp *regexp.Regexp, actionCount *int) (string, *regexp.Regexp, func(string, []int) string, *int) {
	var openTag, endLink, openDesc, endDesc, closeTag string
	if exp == linkExp {
		openTag, endLink, openDesc, endDesc, closeTag = "<a href='", "'", ">", "", "</a>"
	} else if exp == imgExp {
		openTag, endLink, openDesc, endDesc, closeTag = "<img src='", "'", " alt='", "'", ">"
	} else {
		panic("Unknown regex expression provided!")
	}
	return html, exp, func(md string, match []int) string {
		// Whole line match 0 1, 1st group 2 3 desc, 2nd group 4 5 link
		temp := openTag + md[match[4]:match[5]] + endLink + openDesc + md[match[2]:match[3]] + endDesc + closeTag
		if len(md) > match[5]+1 {
			temp += md[match[5]+1:]
		}
		leftOffset := 1
		if exp == imgExp { // Images have an exclamation mark
			leftOffset = 2
		}
		if match[2]-leftOffset > 0 {
			temp = md[:match[2]-leftOffset] + temp
		}
		md = temp
		return md
	}, actionCount
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
