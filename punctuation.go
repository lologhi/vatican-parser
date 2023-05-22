package main

import (
    "strings"
    "regexp"
    md "github.com/JohannesKaufmann/html-to-markdown"
    "github.com/PuerkitoBio/goquery"
)

func getPunctuationRule() md.Rule {
	var footnote_4 = regexp.MustCompile(` \( \[(\d+)\]\(\#\d+\)\)`)
	var footnote_1 = regexp.MustCompile(` \[\[(\d+)\]\]\(\#\_ftn\d+.*?\)`)
	var footnote_1_note = regexp.MustCompile(`\[\[(\d+)\]\]\(\#\_ftnref\d+.*?\)`)
	var footnote_1_1 = regexp.MustCompile(` \[ \[(\d+)\]\(\#\_ftn\d+.*?\)\]`)
	var footnote_1_1_note = regexp.MustCompile(`\[ \[(\d+)\]\(\#\_ftnref\d+.*?\)\]`)
	var footnote_2 = regexp.MustCompile(` \[(\d+)\]`)
	var footnote_2_note = regexp.MustCompile(`\[(\d+)\] `)
	var footnote_3 = regexp.MustCompile(` \((\d{1,3})\)`)
	var footnote_3_note = regexp.MustCompile(`\((\d{1,3})\) `)

	var opening_parenthesis = regexp.MustCompile(`\( `)
	var closing_parenthesis = regexp.MustCompile(` \)`)
	var ponctuation_double = regexp.MustCompile(`(\w)([\:\;\?\!])`)
	var ponctuation_double_simple_space = regexp.MustCompile(` ([\:\;\?\!])`)
	var opening_guillemet = regexp.MustCompile(`«(\w)`)
	var closing_guillemet = regexp.MustCompile(`(\w)»`)
	var short_simple_quotes = regexp.MustCompile(`"(\w+)"`)
	var simple_quotes = regexp.MustCompile(`"(.*?)"`)

	var spaces = regexp.MustCompile(`\s+`)

	changeText := md.Rule{
		Filter: []string{"#text"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			text := selec.Text()
			if trimmed := strings.TrimSpace(text); trimmed == "" {
				return md.String("")
			}

			text = footnote_4.ReplaceAllString(text, "[^$1]")
			text = footnote_1.ReplaceAllString(text, "[^$1]")
			text = footnote_1_note.ReplaceAllString(text, "[^$1]:")
			text = footnote_1_1.ReplaceAllString(text, "[^$1]")
			text = footnote_1_1_note.ReplaceAllString(text, "[^$1]:")
			text = footnote_2.ReplaceAllString(text, "[^$1]")
			text = footnote_2_note.ReplaceAllString(text, "[^$1]:")
			text = footnote_3.ReplaceAllString(text, "[^$1]")
			text = footnote_3_note.ReplaceAllString(text, "[^$1]:")

			text = opening_parenthesis.ReplaceAllString(text, "(")
			text = closing_parenthesis.ReplaceAllString(text, ")")
			text = ponctuation_double.ReplaceAllString(text, "$1 $2")
			text = ponctuation_double_simple_space.ReplaceAllString(text, " $1")
			text = opening_guillemet.ReplaceAllString(text, "« $1")
			text = closing_guillemet.ReplaceAllString(text, "$1 »")
			text = short_simple_quotes.ReplaceAllString(text, "« $1 »")
			text = simple_quotes.ReplaceAllString(text, "« *$1* »")
			
			text = spaces.ReplaceAllString(text, " ")

			// NOTE: See the #text rule for commonmark for all the
			// other logic that should happen here...

			return md.String(text)
		},
	}

	return changeText
}
