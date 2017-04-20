package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

var headers = []string{
	"Дальность:",
	"Компоненты:",
	"Наведение Заклинания:",
	"Продолжительность:",
	"Скорость сотворения:",
	"Сопротивляемость Заклинаниям:",
	"Спасбросок:",
	"Уровень:",
	"Школа:",
	"Школа(Подшкола):",
}

func isHeader(str string) bool {
	for _, h := range headers {
		if strings.HasPrefix(str, h) {
			return true
		}
	}
	return false
}

func main() {
	writer := csv.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	for {
		name_en := strings.Trim(scanner.Text(), "() ")
		scanner.Scan()
		name_ru := strings.TrimSpace(scanner.Text())
		scanner.Scan()
		for isHeader(scanner.Text()) {
			scanner.Scan()
		}

		descr_ru := ""
		more_data := false
		for {
			if scanner.Text() != "" {
				descr_ru = descr_ru + "<p>" + scanner.Text() + "</p>"
			}

			more_data = scanner.Scan()
			new_spell := strings.HasPrefix(scanner.Text(), "(")
			if !more_data || new_spell {
				break
			}
		}

		writer.Write([]string{name_en, name_ru, descr_ru})
		if !more_data {
			break
		}
	}
	writer.Flush()
}
