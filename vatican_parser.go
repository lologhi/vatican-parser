package main

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/gocolly/colly"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
    // parsePopes()
    parseCommissions()
}

func parseCommissions() {
    c := getCollector()
    c.AllowedDomains = []string{"vatican.va", "www.vatican.va"}

    // Parcourir tous les dicastères
    c.OnHTML("#accordionmenu a" , func(e *colly.HTMLElement) {
        // fmt.Println("Dicastère:", strings.TrimSpace(e.Text))
        e.Request.Visit(e.Attr("href"))
    })

    // Page listant les "Documents"
    c.OnHTML(".documenti h2 > a", func(e *colly.HTMLElement) {
        e.Request.Visit(e.Attr("href"))
    })

    // Documents en Français
    c.OnHTML(".documento > .testo > .text ul > li b > a", func(e *colly.HTMLElement) {
        // docUrl          := e.Request.URL.String()
        fmt.Println("Document:", e.Text)
        e.Request.Visit(e.Attr("href"))
    })

    // Documents en Latin
    c.OnHTML(".documento > .testo > .text > ul > li > a", func(e *colly.HTMLElement) {
        if ("Latin" == e.Text) {
            fmt.Println("Document latin:", e.Text)
            e.Request.Visit(e.Attr("href"))
        }
    })

    converter := getConverter()
    c.OnHTML(".documento .testo", func(e *colly.HTMLElement) {
        docUrl          := e.Request.URL.String()
        docContent      := converter.Convert(e.DOM)

        fileName := getFileName(docUrl) // 1901-06-29-en-tout-temps.md
        filePath := getFilePath(docUrl) // leo-xiii/fr/letters/documents
        os.MkdirAll(filePath, 0750)
        // fmt.Println("fileName: ", fileName)
        if err := os.WriteFile(filepath.Join(filePath, fileName), []byte(docContent), 0666); err != nil {
            fmt.Println(err)
        }
    })

    c.Visit(curie)
}

func parsePopes() {
    c := getCollector()
    c.AllowedDomains = []string{"vatican.va", "www.vatican.va"}

    // Parcourir tous les papes
    c.OnHTML("#corpo > table > tbody > tr:nth-child(2) > td > table > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody", func(e *colly.HTMLElement) {
        e.ForEach("tr td a", func(_ int, el *colly.HTMLElement) {
            fmt.Println("Pape:", el.Text)
            el.Request.Visit(el.Attr("href"))
        })
    })

    // Parcourir tous les types de documents
    c.OnHTML("#accordionmenu > ul > li", func(e *colly.HTMLElement) {
        //fmt.Println("Doc type:", e.ChildText("a[0]"))
        e.Request.Visit(e.ChildAttr("a", "href"))

        // Parcourir les sous sections (années et mois) d'un type de document
        e.ForEach("ul li a", func(_ int, el *colly.HTMLElement) {
            //fmt.Println("subDoc type:", strings.TrimSpace(e.Text))
            el.Request.Visit(el.Attr("href"))
        })
    })

    // Parcourir les documents d'un type donné
    c.OnHTML(".vaticanindex h1 a", func(e *colly.HTMLElement) {
        e.Request.Visit(e.Attr("href"))
    })

    c.OnHTML(".vaticanindex h2 a", func(e *colly.HTMLElement) {
        if ("Latin" == e.Text) {
            e.Request.Visit(e.Attr("href"))
        }
    })

    converter := getConverter()
    c.OnHTML(".documento .testo", func(e *colly.HTMLElement) {
        docUrl          := e.Request.URL.String()
        docContent      := converter.Convert(e.DOM)

        fileName := getFileName(docUrl) // 1901-06-29-en-tout-temps.md
        filePath := getFilePath(docUrl) // leo-xiii/fr/letters/documents
        os.MkdirAll(filePath, 0750)
        // fmt.Println("fileName: ", fileName)
        if err := os.WriteFile(filepath.Join(filePath, fileName), []byte(docContent), 0666); err != nil {
            fmt.Println(err)
        }
    })

    // Parse all eleven previous popes :
    c.Visit(popes)
}


var (
    reName = regexp.MustCompile(`([0-9]{8})(.*)`)
    // 1 : pope name and doc type (skipped as already in the path)
    // 2 : doc date
    // 3 : _
    // 4 : doc name
    reDate = regexp.MustCompile(`([0-9]{8})`)
    popes = "https://www.vatican.va/holy_father/index_fr.htm"
    curie = "https://www.vatican.va/content/romancuria/fr/segreteria-di-stato/segreteria-di-stato.index.html"
)

func getCollector() *colly.Collector {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
    })

    return c
}

func getConverter() *md.Converter {
    opt := &md.Options{
        CodeBlockStyle: "fenced", // default: indented
        EmDelimiter: "*", // default: _
    }
    converter := md.NewConverter("", true, opt)
    converter.AddRules(getPunctuationRule())
    converter.Use(plugin.GitHubFlavored())

    return converter
}

func getOriginalName(url string) string {
    // url example : https://www.vatican.va/content/leo-xiii/fr/letters/documents/hf_l-xiii_let_19010629_en-tout-temps.html
    baseUrl   := path.Base(url) // hf_l-xiii_let_19010629_en-tout-temps.html
    
    return strings.TrimSuffix(baseUrl, filepath.Ext(baseUrl)) // hf_l-xiii_let_19010629_en-tout-temps
}

func getDocName(url string) string {
    matchedName := reName.FindStringSubmatch(getOriginalName(url))
    if (2 > len(matchedName)) {
        return ""
    } else {
        return strings.Trim(matchedName[2], "-_")
    }
}

func getDocDate(url string) string {
    dateStringInURL := reDate.FindString(url)
    if ("" != dateStringInURL) {
        docDate, _ := time.Parse("20060102", dateStringInURL)
        if ("0001-01-01" == docDate.Format("2006-01-02")) {
            docDate, _ = time.Parse("02012006", dateStringInURL)
        }
        return docDate.Format("2006-01-02")
    } else {
        return ""
    }
}

func getFileName(url string) string {
    ext := ".md"
    if (strings.Contains(url, "/la/")) {
        ext = ".latin.md"
    }

    docName := getDocName(url)
    docDate := getDocDate(url)

    if ("" == docName && "" == docDate) {
        return getOriginalName(url) + ext
    }
    if ("" == docName) {
        return docDate + ext
    }
    if ("" == docDate) {
        return docName + ext
    }

    return docDate + "-" + docName + ext
}

func getFilePath(docUrl string) string {
    cleanedPath := strings.TrimSuffix(path.Dir(docUrl), "documents")
    cleanedPath = strings.Replace(cleanedPath, "/fr/", "/", 1)
    cleanedPath = strings.Replace(cleanedPath, "/la/", "/", 1)
    return strings.TrimPrefix(cleanedPath, "https:/www.vatican.va/content/")
}
