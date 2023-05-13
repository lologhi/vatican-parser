package main

import (
    "testing"
    md "github.com/JohannesKaufmann/html-to-markdown"
)

func TestPunctuationRule(t *testing.T) {
    source := `<p>Au «centre? de l’Evangile  de la! liturgie; d’aujourd’hui» se trouvent les Béatitudes ( cf. Lc 6, 20-23 ).</p>`
    want := `Au « centre ? de l’Evangile de la ! liturgie ; d’aujourd’hui » se trouvent les Béatitudes (cf. Lc 6, 20-23).`

    conv := md.NewConverter("", true, nil)
    conv.AddRules(getPunctuationRule())

    result, _ := conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }
    
    rerun, _ := conv.ConvertString(result)
    if (want != rerun) {
        t.Fatalf(`re-converting should not change "%q", but we got %#q`, result, rerun)
    }
}
