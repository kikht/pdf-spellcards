package main

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"strings"
)

const ARROW = "→"

var (
	style_regexp     = regexp.MustCompile(` style="[^>]*"`)
	font_size_regexp = regexp.MustCompile(`font-size:[^;" ]*;?`)
	a_open_regexp    = regexp.MustCompile(`<a[^>]*>`)
	a_close_regexp   = regexp.MustCompile(`</a[^>]*>`)
	img_regexp       = regexp.MustCompile(`<img[^>]*>`)
)

func cleanHtml(html string) string {
	html = font_size_regexp.ReplaceAllString(html, "")
	html = img_regexp.ReplaceAllString(html, "")
	html = a_open_regexp.ReplaceAllString(html, "<i>")
	html = a_close_regexp.ReplaceAllString(html, "</i>")
	return html
}

func main() {
	doc, err := goquery.NewDocumentFromReader(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	name_ru := doc.Find("h3#sites-page-title-header span#sites-page-title").Text()
	src := doc.Find("#sites-canvas-main-content > table > tbody > tr > td > div > div:last-child > ul > li > font").Text()
	name_en := strings.TrimSpace(src[strings.Index(src, ARROW)+len(ARROW):])
	descr := ""
	doc.Find("#sites-canvas-main-content > table > tbody > tr > td > div > div > div > div").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "ОПИСАНИЕ" {
			s.NextAll().Each(func(i int, c *goquery.Selection) {
				html, err := goquery.OuterHtml(c)
				if err != nil {
					log.Println(err)
				}
				descr += cleanHtml(html)
			})
		}
	})

	w := csv.NewWriter(os.Stdout)
	w.Write([]string{name_en, name_ru, descr})
	w.Flush()
}
