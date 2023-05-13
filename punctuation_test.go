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

func TestWithRealExamples(t *testing.T) {
    conv := md.NewConverter("", true, nil)
    conv.AddRules(getPunctuationRule())

    source := `IV. Conditae sodalitatis constitutiones Ordinarius recognoscat : verum ne prius approbet, quam eas ad normam eorum, quae Sacrum Consilium in hac causa decrevit, exigendas curaverit [[3]](#_ftn3 "").`
    want := `IV. Conditae sodalitatis constitutiones Ordinarius recognoscat : verum ne prius approbet, quam eas ad normam eorum, quae Sacrum Consilium in hac causa decrevit, exigendas curaverit[^3].`
    result, _ := conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }

    source = `[[3]](#_ftnref3 "") Ordinarius nempe approbare nequit constitutiones Instituti, nisi hae conformes sint Normis, quibus Sancta Sedes uti solet in novis Institutis approbandis, quum consentaneum omnino sit ut Institutum iam a suo exordio iuxta Apostolicas regulas moderetur (N. R.).`
    want = `[^3]: Ordinarius nempe approbare nequit constitutiones Instituti, nisi hae conformes sint Normis, quibus Sancta Sedes uti solet in novis Institutis approbandis, quum consentaneum omnino sit ut Institutum iam a suo exordio iuxta Apostolicas regulas moderetur (N. R.).`
    result, _ = conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }

    source = `atque inde fit, quemadmodum S. Cyprianus monet [2], *ut Ecclesia super Episcopos constituatur, et omnis actus Ecclesiae per eosdem Praepositos gubernetur.*`
    want = `atque inde fit, quemadmodum S. Cyprianus monet[^2], *ut Ecclesia super Episcopos constituatur, et omnis actus Ecclesiae per eosdem Praepositos gubernetur.*`
    result, _ = conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }

    source = `[2] Epist. 29 *ad lapsos.*`
    want = `[^2]:Epist. 29 *ad lapsos.*`
    result, _ = conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }

    source = `Romani Imperii auctoritate saecula plurima sociavisset, is et proprius Apostolicae Sedis evaderet (3) et, posteritati servatus, christianos Europae populos alios cum aliis arto unitatis vinculo coniungeret.`
    want = `Romani Imperii auctoritate saecula plurima sociavisset, is et proprius Apostolicae Sedis evaderet[^3] et, posteritati servatus, christianos Europae populos alios cum aliis arto unitatis vinculo coniungeret.`
    result, _ = conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }

    source = `(3) Epist. S. Congr. Stud. *Vehementer sane,* ad Ep. universos, 1 Iul. 1908 : *Ench. Cler.,* N. 820. Cfr etiam Epist. Ap. Pii XI, *[Unigenitus Dei Filius](/content/pius-xi/la/apost_letters/documents/hf_p-xi_apl_19240319_unigenitus-dei.html)*, 19 Mar. 1924 : *A.A.S.* 16 (1924), 141.`
    want = `[^3]:Epist. S. Congr. Stud. *Vehementer sane,* ad Ep. universos, 1 Iul. 1908 : *Ench. Cler.,* N. 820. Cfr etiam Epist. Ap. Pii XI, *[Unigenitus Dei Filius](/content/pius-xi/la/apost_letters/documents/hf_p-xi_apl_19240319_unigenitus-dei.html)*, 19 Mar. 1924 : *A.A.S.* 16 (1924), 141.`
    result, _ = conv.ConvertString(source)
    if (want != result) {
        t.Fatalf(`converting "%q" returned\n %q, want match for\n %#q`, source, result, want)
    }
}
