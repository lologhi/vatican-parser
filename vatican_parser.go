package main

import (
    "os"
    "fmt"
    "github.com/gocolly/colly"
    "strings"
    "path"
    "path/filepath"
    "regexp"
    "time"
    md "github.com/JohannesKaufmann/html-to-markdown"
    "github.com/JohannesKaufmann/html-to-markdown/plugin"
)

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains("vatican.va", "www.vatican.va"),
    )
    opt := &md.Options{
        CodeBlockStyle: "fenced", // default: indented
        EmDelimiter: "*", // default: _
    }
    converter := md.NewConverter("", true, opt)
    converter.Use(plugin.GitHubFlavored())

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
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

    c.OnHTML(".documento .testo", func(e *colly.HTMLElement) {
        docUrl          := e.Request.URL.String()
        docBaseUrlExt   := path.Base(docUrl)
        docBaseUrl      := strings.TrimSuffix(docBaseUrlExt, filepath.Ext(docBaseUrlExt))
        re              := regexp.MustCompile(`[0-9]{8}`)
        dateStringInURL := re.FindString(docBaseUrl)
        docDate, _      := time.Parse("20060102", dateStringInURL)
        docContent      := converter.Convert(e.DOM)

        parts := strings.Split(path.Dir(docUrl), "/")
        savePath := filepath.Join(parts[3],parts[5],parts[6])
        // fmt.Println("savePath:", savePath)
        os.MkdirAll(savePath, 0750)

        fileName := colly.SanitizeFileName(docDate.Format("2006-01-02")+"_"+docBaseUrl+".md")
        // fmt.Println("fileName: ", fileName)
        if err := os.WriteFile(filepath.Join(savePath, fileName), []byte(docContent), 0666); err != nil {
            fmt.Println(err)
        }
    })

    c.Visit("https://www.vatican.va/content/francesco/fr.html")
}
