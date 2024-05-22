package main

import (
	"testing"
)

func TestGetDocNameUrl(t *testing.T) {
	url := "https://www.vatican.va/content/leo-xiii/fr/letters/documents/hf_l-xiii_let_19010629_en-tout-temps.html"
	want := "en-tout-temps"
	msg := getDocName(url)
	if want != msg {
		t.Fatalf(`getDocName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/content/francesco/fr/speeches/2015/december/documents/papa-francesco_20151218_donatori-presepio-albero.html"
	want = "donatori-presepio-albero"
	msg = getDocName(url)
	if want != msg {
		t.Fatalf(`getDocName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/content/francesco/fr/cotidie/2013/documents/papa-francesco-cotidie_20131220.html"
	want = ""
	msg = getDocName(url)
	if want != msg {
		t.Fatalf(`getDocName("%q") returned %q, want match for %#q`, url, msg, want)
	}
}

func TestGetDocDateWrongFormat(t *testing.T) {
	url := "https://www.vatican.va/content/pius-x/fr/encyclicals/documents/hf_p-x_enc_15041905_acerbo-nimis.html"
	want := "1905-04-15"
	msg := getDocDate(url)
	if want != msg {
		t.Fatalf(`getDocDate("%q") returned %q, want match for %#q`, url, msg, want)
	}
}

func TestGetDocDateWithoutDateInUrl(t *testing.T) {
	url := "https://www.vatican.va/content/francesco/fr/travels/2013/documents/papa-francesco-programma-gmg-rio-de-janeiro-2013.html"
	want := ""
	msg := getDocDate(url)
	if want != msg {
		t.Fatalf(`getDocDate("%q") returned %q, want match for %#q`, url, msg, want)
	}
}

func TestGetFileNameUrl(t *testing.T) {
	url := "https://www.vatican.va/content/leo-xiii/fr/letters/documents/hf_l-xiii_let_19010629_en-tout-temps.html"
	want := "1901-06-29-en-tout-temps.md"
	msg := getFileName(url)
	if want != msg {
		t.Fatalf(`getFileName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/content/francesco/fr/speeches/2015/december/documents/papa-francesco_20151218_donatori-presepio-albero.html"
	want = "2015-12-18-donatori-presepio-albero.md"
	msg = getFileName(url)
	if want != msg {
		t.Fatalf(`getFileName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/content/john-xxiii/la/apost_constitutions/1962/documents/hf_j-xxiii_apc_19620222_veterum-sapientia.html"
	want = "1962-02-22-veterum-sapientia.latin.md"
	msg = getFileName(url)
	if want != msg {
		t.Fatalf(`getFileName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/content/francesco/fr/cotidie/2013/documents/papa-francesco-cotidie_20131220.html"
	want = "2013-12-20.md"
	msg = getFileName(url)
	if want != msg {
		t.Fatalf(`getFileName("%q") returned %q, want match for %#q`, url, msg, want)
	}

	url = "https://www.vatican.va/roman_curia/congregations/ccdds/documents/rc_con_ccdds_doc_20190125_decreto-celebrazione-paolovi_fr.html"
	want = "2019-01-25-decreto-celebrazione-paolovi_fr.md"
	msg = getFileName(url)
	if want != msg {
		t.Fatalf(`getFileName("%q") returned %q, want match for %#q`, url, msg, want)
	}
}

func TestGetFileNameWithoutDateInUrl(t *testing.T) {
	url := "https://www.vatican.va/content/francesco/fr/travels/2013/documents/papa-francesco-programma-gmg-rio-de-janeiro-2013.html"
	want := "papa-francesco-programma-gmg-rio-de-janeiro-2013.md"
	msg := getFileName(url)
	if want != msg {
		t.Fatalf(`getDocName("%q") returned %q, want match for %#q`, url, msg, want)
	}
}

func TestGetFilePath(t *testing.T) {
	url := "https://www.vatican.va/content/francesco/fr/travels/2013/documents/papa-francesco-programma-gmg-rio-de-janeiro-2013.html"
	mainPath := "/download/vatican/"
	want := "/download/vatican/francesco/travels/2013"
	msg := getFilePath(mainPath, url)
	if want != msg {
		t.Fatalf(`getFilePath(%q, %q) returned %q, want match for %#q`, mainPath, url, msg, want)
	}

	url = "https://www.vatican.va/content/benedict-xv/fr/letters/1921/documents/hf_ben-xv_let_19210124_la-singolare.html"
	want = "/download/vatican/benedict-xv/letters/1921"
	msg = getFilePath(mainPath, url)
	if want != msg {
		t.Fatalf(`getFilePath(%q, %q) returned %q, want match for %#q`, mainPath, url, msg, want)
	}

	url = "https://www.vatican.va/content/john-xxiii/la/apost_constitutions/1962/documents/hf_j-xxiii_apc_19620222_veterum-sapientia.html"
	want = "/download/vatican/john-xxiii/apost_constitutions/1962"
	msg = getFilePath(mainPath, url)
	if want != msg {
		t.Fatalf(`getFilePath(%q, %q) returned %q, want match for %#q`, mainPath, url, msg, want)
	}

	url = "https://www.vatican.va/roman_curia/congregations/cfaith/cti_documents/rc_con_cfaith_doc_20040723_communion-stewardship_fr.html"
	want = "cti/"
	msg = getFilePath(mainPath, url)
	if want != msg {
		t.Fatalf(`getFilePath(%q, %q) returned %q, want match for %#q`, mainPath, url, msg, want)
	}
}
