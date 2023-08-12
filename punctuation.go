package main

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func getPunctuationRule() md.Rule {
	var footnote4 = regexp.MustCompile(` \( \[(\d+)]\(#\d+\)\)`)
	var footnote1 = regexp.MustCompile(` \[\[(\d+)]]\(#_ftn\d+.*?\)`)
	var footnote1Note = regexp.MustCompile(`\[\[(\d+)]]\(#_ftnref\d+.*?\)`)
	var footnote11 = regexp.MustCompile(` \[ \[(\d+)]\(#_ftn\d+.*?\)]`)
	var footnote11Note = regexp.MustCompile(`\[ \[(\d+)]\(#_ftnref\d+.*?\)]`)
	var footnote2 = regexp.MustCompile(` \[(\d+)]`)
	var footnote2Note = regexp.MustCompile(`\[(\d+)] `)
	var footnote3 = regexp.MustCompile(` \((\d{1,3})\)`)
	var footnote3Note = regexp.MustCompile(`\((\d{1,3})\) `)

	var openingParenthesis = regexp.MustCompile(`\( `)
	var closingParenthesis = regexp.MustCompile(` \)`)
	var ponctuationDouble = regexp.MustCompile(`(\w)([:;?!])`)
	var ponctuationDoubleSimpleSpace = regexp.MustCompile(` ([:;?!])`)
	var openingGuillemet = regexp.MustCompile(`«(\w)`)
	var closingGuillemet = regexp.MustCompile(`(\w)»`)
	var shortSimpleQuotes = regexp.MustCompile(`"(\w+)"`)
	var simpleQuotes = regexp.MustCompile(`"(.*?)"`)

	var spaces = regexp.MustCompile(`\s+`)

	changeText := md.Rule{
		Filter: []string{"#text"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			text := selec.Text()
			if trimmed := strings.TrimSpace(text); trimmed == "" {
				return md.String("")
			}

			text = footnote4.ReplaceAllString(text, "[^$1]")
			text = footnote1.ReplaceAllString(text, "[^$1]")
			text = footnote1Note.ReplaceAllString(text, "[^$1]:")
			text = footnote11.ReplaceAllString(text, "[^$1]")
			text = footnote11Note.ReplaceAllString(text, "[^$1]:")
			text = footnote2.ReplaceAllString(text, "[^$1]")
			text = footnote2Note.ReplaceAllString(text, "[^$1]:")
			text = footnote3.ReplaceAllString(text, "[^$1]")
			text = footnote3Note.ReplaceAllString(text, "[^$1]:")

			text = openingParenthesis.ReplaceAllString(text, "(")
			text = closingParenthesis.ReplaceAllString(text, ")")
			text = ponctuationDouble.ReplaceAllString(text, "$1 $2")
			text = ponctuationDoubleSimpleSpace.ReplaceAllString(text, " $1")
			text = openingGuillemet.ReplaceAllString(text, "« $1")
			text = closingGuillemet.ReplaceAllString(text, "$1 »")
			text = shortSimpleQuotes.ReplaceAllString(text, "« $1 »")
			text = simpleQuotes.ReplaceAllString(text, "« *$1* »")

			text = spaces.ReplaceAllString(text, " ")

			// NOTE: See the #text rule for commonmark for all the
			// other logic that should happen here...

			return md.String(text)
		},
	}

	return changeText
}
