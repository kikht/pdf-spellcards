package main

import (
	"encoding/csv"
	"errors"
	"flag"
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

	"1 min.":                "1 минута",
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
	"instant":                    "мгновенное",
	"instantaneous":              "мгновенное",
	"Instantaneous":              "мгновенное",
	"instantaneous; see text":    "мгновенное, см. текст",
	"permanent":                  "постоянное",
	"permanent (D)":              "постоянное (до отмены)",
	"permanent until discharged": "постоянное (до отмены)",
	"concentration":              "концентрация",
	"Concentration + 2 rounds":   "концентрация + 2 раунда",

	"concentration, up to 1 minute/level":        "концентрация, 1 мин./ур.",
	"concentration, up to 1 minute/ level":       "концентрация, 1 мин./ур.",
	"concentration, up to 1 min./level":          "концентрация, 1 мин./ур.",
	"concentration, up to 10 min./level":         "концентрация, 10 мин./ур.",
	"concentration, up to 10 min./ level":        "концентрация, 10 мин./ур.",
	"1 round/level":                              "1 раунд/уровень",
	"1 round + 1 round per three levels":         "1 раунд + 1 раунд/3 ур.",
	"3 rounds/level":                             "3 раунда/уровень",
	"1 minute or until discharged":               "1 минута или до отмены",
	"1 min./level":                               "1 минута/уровень",
	"1 min./level (D)":                           "1 мин./ур. или до отмены",
	"10 min./level":                              "10 минут/уровень",
	"30 minutes or until discharged":             "30 минут или до отмены",
	"1d4 rounds or 1 round; see text":            "1d4 раунда или 1 раунд, см. текст",
	"2d4 rounds":                                 "2d4 раунда",
	"1d6+2 rounds":                               "1d6+2 раунда",
	"1 hour/level":                               "1 час/уровень",
	"1 hour/level (D)":                           "1 час/уровень или до отмены",
	"1 hour/level or until discharged; see text": "1 час/уровень или до отмены; см. текст",
	"2 hours/level":                              "2 часа/уровень",
	"until landing or 1 round/level":             "до приземл. или 1 раунд/ур.",
	"1 day/level":                                "1 день/уровень",

	//Range
	"personal":                                    "на себя",
	"touch":                                       "касание",
	"personal or touch":                           "на себя или касание",
	"personal or close (25 ft. + 5 ft./2 levels)": "на себя или близкая (25 фт. + 5 фт./2 уровня)",

	"close (25 ft. + 5 ft./2 levels)": "близкая (25 фт. + 5 фт./2 уровня)",
	"medium (100 ft. + 10 ft./level)": "средняя (100 фт. + 10 фт./уровнень)",
	"medium (100 ft. + 10 ft. level)": "средняя (100 фт. + 10 фт./уровнень)",
	"long (400 ft. + 40 ft./level)":   "дальняя (400 фт. + 40 фт./уровень)",

	"10 ft.": "10 футов",
	"15 ft.": "15 футов",
	"20 ft.": "20 футов",
	"30 ft.": "30 футов",
	"50 ft.": "50 футов",
	"60 ft.": "60 футов",

	//Area
	"cone-shaped emanation":                "коническая эманация",
	"cone-shaped burst":                    "конический всплеск",
	"50-ft.-radius burst, centered on you": "50 футов вокруг колдующего",

	"20-ft.-radius emanation centered on a point in space":                    "20 футовая сфера с центром в заданной точке",
	"one creature, one object, or a 5-ft. cube":                               "одно создание, один объект или 5 футовый куб",
	"The caster and all allies within a 50-ft. burst, centered on the caster": "Колдующий и союзники в пределах 50 футов",
	"one or more living creatures within a 10-ft.-radius burst":               "одно или несколько живых созданий внутри 10 футовой сферы",
	"several living creatures, no two of which may be more than 30 ft. apart": "несколько живых созданий, никакие два из которых не стоят дальше 30 футов друг от друга",
}

func T(str string) string {
	if res, ok := translation[str]; ok {
		return res
	} else {
		if str != "" {
			log.Println("untranslated:", str)
		}
		return str
	}
}

type Component struct {
	Name, Image string
}

var (
	COMP_VERBAL       = Component{"верб.", "verbal"}
	COMP_SOMATIC      = Component{"сомат.", "somatic"}
	COMP_MATERIAL     = Component{"матер.", "material"}
	COMP_FOCUS        = Component{"фокус", "focus"}
	COMP_DIVINE_FOCUS = Component{"бож.", "divine_focus"}
	COMP_F_DF         = Component{"фокус/бож.", "f_df"}
	COMP_M_DF         = Component{"матер./бож.", "m_df"}
)

type Spell struct {
	Name, School, Effect      string
	CastTime, Duration, Range string
	Area, AreaImg             string
	Descriptors               []string
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

func GetAreaImg(rng, area string) string {
	switch {
	case strings.HasPrefix(area, "cone-shaped"):
		switch rng {
		case "15 ft.":
			return "cone-15"
		case "30 ft.":
			return "cone-30"
		case "60 ft.":
			return "cone-60"
		default:
			return ""
		}
	case strings.Contains(area, "5-ft.-radius"):
		return "radius-5"
	case strings.Contains(area, "10-ft.-radius"):
		return "radius-10"
	case strings.Contains(area, "20-ft.-radius"):
		return "radius-20"
	case strings.Contains(area, "30-ft.-radius"):
		return "radius-30"
	default:
		return ""
	}
}

func GetComponents(str string) []Component {
	var comp []Component
	for _, item := range strings.Split(str, ", ") {
		item = strings.TrimSpace(item)
		switch {
		case strings.HasPrefix(item, "M/DF") || strings.HasPrefix(item, "DF/M"):
			comp = append(comp, COMP_M_DF)
		case strings.HasPrefix(item, "F/DF") || strings.HasPrefix(item, "DF/F"):
			comp = append(comp, COMP_F_DF)
		case strings.HasPrefix(item, "M"):
			comp = append(comp, COMP_MATERIAL)
		case strings.HasPrefix(item, "F"):
			comp = append(comp, COMP_FOCUS)
		case strings.HasPrefix(item, "V"):
			comp = append(comp, COMP_VERBAL)
		case strings.HasPrefix(item, "S"):
			comp = append(comp, COMP_SOMATIC)
		case strings.HasPrefix(item, "DF"):
			comp = append(comp, COMP_DIVINE_FOCUS)
		default:
			//log.Println("Unknown component: ", item)
		}
	}
	return comp
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
			Effect:      value("effect_ru", "effect"),
			CastTime:    T(value("casting_time")),
			Duration:    T(value("duration")),
			Range:       T(value("range")),
			Area:        T(value("area")),
			AreaImg:     GetAreaImg(value("range"), value("area")),
			Description: template.HTML(value("description_formated_ru", "description_formated")),
			Level:       level,
			Components:  GetComponents(value("components")),
		}
		if descStr := value("descriptor"); descStr != "" {
			spell.Descriptors = strings.Split(descStr, ", ")
		}
		data.Spells = append(data.Spells, spell)
	}

	err = tmpl.Execute(writer, data)
	return err
}

func main() {
	var (
		class = flag.String("class", "wiz", "Character class")
		level = flag.Int("level", 0, "Spell level")
	)
	flag.Parse()
	err := GenerateCards(os.Stdout, *class, *level)
	if err != nil {
		log.Fatalln(err)
	}
}
