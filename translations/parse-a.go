package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"regexp"
	"strings"
)

var (
	note_regexp         = regexp.MustCompile(`{{\?{1,2}\|[\pL\pN ]*}}`)
	untranslated_regexp = regexp.MustCompile(`{{\?{3}\|([\pL\pN ]*)}}`)
	italic_regexp       = regexp.MustCompile(`''([\pL\pN\ ]*)''`)
)

func norm(str string) string {
	str = note_regexp.ReplaceAllString(str, "")
	str = untranslated_regexp.ReplaceAllString(str, "<i>$1</i>")
	str = italic_regexp.ReplaceAllString(str, "<i>$1</i>")
	return str
}

func main() {
	writer := csv.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	for {
		name_ru := norm(strings.Trim(scanner.Text(), "= "))
		scanner.Scan()
		name_en := norm(strings.Trim(scanner.Text(), "{?|}"))
		scanner.Scan()
		for strings.HasPrefix(scanner.Text(), ": '''") {
			scanner.Scan()
		}

		descr_ru := "<p>"
		scanner.Scan()
		if scanner.Text() != "" {
			descr_ru += norm(scanner.Text())
		}
		new_par := false
		more_data := false
		for {
			more_data = scanner.Scan()
			new_spell := strings.HasPrefix(scanner.Text(), "===")
			if !more_data || new_spell {
				break
			}

			if scanner.Text() == "" {
				descr_ru += "</p>"
				new_par = true
			} else {
				if new_par {
					descr_ru += "<p>"
					new_par = false
				}
				descr_ru += norm(scanner.Text())
			}
		}

		writer.Write([]string{name_en, name_ru, descr_ru})
		if !more_data {
			break
		}
	}
	writer.Flush()
}
