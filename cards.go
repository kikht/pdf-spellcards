package main

import (
	"encoding/csv"
	"errors"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var translation = map[string]string{
	//Cast time
	"1 immediate action": "быстрое действие",
	"1 swift action":     "быстрое действие",

	"1 standard action or immediate action; see text": "основное или быстрое действие",

	"standard action":               "основное действие",
	"1 standard":                    "основное действие",
	"1 standard action":             "основное действие",
	"1 standard action or see text": "основное действие, см. текст",

	"full-round action":                      "полный раунд",
	"1 full round":                           "полный раунд",
	"1 full-round action":                    "полный раунд",
	"1 round":                                "полный раунд",
	"1 round; see text":                      "полный раунд, см. текст",
	"1 full-round action, special see below": "полный раунд, см. текст",

	"2 rounds":      "2 раунда",
	"3 rounds":      "3 раунда",
	"3 full rounds": "3 раунда",
	"6 rounds":      "6 раундов",

	"1 minute":              "1 минута",
	"Casting time 1 minute": "1 минута",
	"1 minute per page":     "1 минута на страницу",
	"1 minute/HD of target": "1 минута/HD цели",
	"1 minute/lb. created":  "1 минута/фунт",

	"2 minutes": "2 минуты",

	"10 minutes":                                      "10 минут",
	"10 minute/HD of target":                          "10 минут/HD цели",
	"10 minutes (see text)":                           "10 минут, см. текст",
	"10 minutes; see text":                            "10 минут, см. текст",
	"10 minutes or more; see text":                    "10 минут или больше",
	"at least 10 minutes; see text":                   "10 минут или больше",
	"10 minutes, plus length of memory to be altered": "10 минут + длина памяти",

	"30 minutes": "30 минут",

	"1 hour":   "1 час",
	"2 hours":  "2 часа",
	"4 hours":  "4 часа",
	"6 hours":  "6 часов",
	"8 hours":  "8 часов",
	"12 hours": "12 часов",
	"24 hours": "24 часа",
	"1 day":    "1 день",
	"1 week":   "1 неделя",
	"see text": "см. текст",

	//Duration
	"instant":       "мгновенное",
	"instantaneous": "мгновенное",
	"Instantaneous": "мгновенное",

	//Range
	"personal": "на себя",
	"touch":    "касание",
	"close (25 ft. + 5 ft./2 levels)": "близкая (25 фт. + 5 фт./2 уровня)",
	"medium (100 ft. + 10 ft./level)": "средняя (100 фт. + 10 фт./уровнень)",
	"10 ft.": "10 футов",
	"60 ft.": "60 футов",

	//Area
}

func T(str string) string {
	if res, ok := translation[str]; ok {
		return res
	} else {
		return str
	}
}

type Component struct {
	Name, Image string
}

type Spell struct {
	Name, School, Effect      string
	CastTime, Duration, Range string
	Area, AreaImg             string
	Descriptor                string
	Description               template.HTML
	Level                     int
	Components                []Component
}

type TemplateData struct {
	Spells []Spell
	Class  string
}

var tmpl = template.Must(template.ParseFiles("cards.template"))
var classes = map[string]string{
	"sor":         "чародей",
	"wiz":         "волшебник",
	"cleric":      "жрец",
	"druid":       "друид",
	"ranger":      "рейнджер",
	"bard":        "бард",
	"paladin":     "паладин",
	"alchemist":   "алхимик",
	"summoner":    "заклинатель",
	"witch":       "ведьма",
	"inquisitor":  "инквизитор",
	"oracle":      "оракул",
	"antipaladin": "антипаладин",
	"magus":       "магус",
	"adept":       "адепт",
}

func GetLevel(str string) (int, error) {
	if str == "NULL" {
		return -1, nil
	}
	return strconv.Atoi(str)
}

func GetAreaImg(area string) string {
	return ""
}

func GetComponents(str string) []Component {
	return nil
}

func GenerateCards(writer io.Writer, class string, level int) error {
	var (
		data TemplateData
		ok   bool
	)
	data.Class, ok = classes[class]
	if !ok {
		return errors.New("no such class: " + class)
	}

	file, err := os.Open("spells.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return err
	}
	var line []string
	value := func(keys ...string) string {
		for _, key := range keys {
			for i, col := range header {
				if key == col {
					return line[i]
				}
			}
		}
		return ""
	}

	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		curLevel, err := GetLevel(value(class))
		if err != nil {
			return errors.New("can't parse level value: " + err.Error())
		}
		if curLevel != level || value("source") != "PFRPG Core" {
			continue
		}

		spell := Spell{
			Name:        value("name_ru", "name"),
			School:      value("school"),
			Descriptor:  strings.Join(strings.Split(value("descriptor"), ", "), " "),
			Effect:      value("effect_ru", "effect"),
			CastTime:    T(value("casting_time")),
			Duration:    T(value("duration")),
			Range:       T(value("range")),
			Area:        T(value("area")),
			AreaImg:     GetAreaImg(value("area")),
			Description: template.HTML(value("description_formated_ru", "description_formated")),
			Level:       level,
			Components:  GetComponents(value("components")),
		}
		data.Spells = append(data.Spells, spell)
	}

	err = tmpl.Execute(writer, data)
	return err
}

func main() {
	err := GenerateCards(os.Stdout, "wiz", 0)
	if err != nil {
		log.Fatalln(err)
	}
}
