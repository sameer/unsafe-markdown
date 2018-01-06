package markdown

import (
	"regexp"

)

var (
	// BOL, 1 to 6 #s, a space, some text, EOL.
	header  = regexp.MustCompile("(?m)^([#]{1,6})[ ]([\\s\\S]+?)$")
	// 1 star, char that's not a star, optional text, star, char that's not a star or is EOL. Multiline.
	italics = regexp.MustCompile("[*]{1}([^*][\\s\\S]*?[*]{1}[^*$]{1})")
	// 2 stars, char that's not a star, optional text, 2 stars, char that's not a star. Multiline.
	bold = regexp.MustCompile("[*]{2}([^*][\\s\\S]*?[*]{1}[^*]{2})")
	// >, a space, and required text.
	blockquote = regexp.MustCompile("(?m)>[ ]([\\s\\S]+?)$")
	// non ! char, [, required text, ], (, optional text, ).
	link = regexp.MustCompile("[^!]\\[([\\s\\S])+?\\]\\(([\\s\\S]*?)\\)")
	// !, [, required text, ], (, optional text, ).
	img = regexp.MustCompile("!\\[([\\s\\S])+?\\]\\(([\\s\\S]*?)\\)")
)

func MarkdownToHtml(md string) string {
	return ""
}
