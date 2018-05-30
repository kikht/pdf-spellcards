package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"html/template"
	"io"
	"log"
	"os"
	"sort"
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

	"concentration, up to 1 minute/level":        "концентрация, 1 мин./УЗ",
	"concentration, up to 1 minute/ level":       "концентрация, 1 мин./УЗ",
	"concentration, up to 1 min./level":          "концентрация, 1 мин./УЗ",
	"concentration, up to 10 min./level":         "концентрация, 10 мин./УЗ",
	"concentration, up to 10 min./ level":        "концентрация, 10 мин./УЗ",
	"1 round/level":                              "1 раунд/УЗ",
	"1 round + 1 round per three levels":         "1 раунд + 1 раунд/3 УЗ",
	"3 rounds/level":                             "3 раунда/УЗ",
	"1 minute or until discharged":               "1 минута или до отмены",
	"1 min./level":                               "1 минута/УЗ",
	"1 min./level (D)":                           "1 мин./УЗ или до отмены",
	"10 min./level":                              "10 минут/УЗ",
	"30 minutes or until discharged":             "30 минут или до отмены",
	"1d4 rounds or 1 round; see text":            "1d4 раунда или 1 раунд, см. текст",
	"2d4 rounds":                                 "2d4 раунда",
	"1d6+2 rounds":                               "1d6+2 раунда",
	"1 hour/level":                               "1 час/УЗ",
	"1 hour/level (D)":                           "1 час/УЗ или до отмены",
	"1 hour/level or until discharged; see text": "1 час/УЗ или до отмены; см. текст",
	"2 hours/level":                              "2 часа/УЗ",
	"until landing or 1 round/level":             "до приземл. или 1 раунд/УЗ",
	"1 day/level":                                "1 день/УЗ",

	//Range
	"personal":                                    "на себя",
	"touch":                                       "касание",
	"personal or touch":                           "на себя или касание",
	"personal or close (25 ft. + 5 ft./2 levels)": "на себя или близкая (25 фт. + 5 фт./2 УЗ)",

	"close (25 ft. + 5 ft./2 levels)": "близкая (25 фт. + 5 фт./2 УЗ)",
	"medium (100 ft. + 10 ft./level)": "средняя (100 фт. + 10 фт./УЗ)",
	"medium (100 ft. + 10 ft. level)": "средняя (100 фт. + 10 фт./УЗ)",
	"long (400 ft. + 40 ft./level)":   "дальняя (400 фт. + 40 фт./УЗ)",

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
	COMP_VERBAL       = Component{"слов.", "verbal"}
	COMP_SOMATIC      = Component{"жест.", "somatic"}
	COMP_MATERIAL     = Component{"реаг.", "material"}
	COMP_FOCUS        = Component{"фок.п.", "focus"}
	COMP_DIVINE_FOCUS = Component{"сакр.", "divine_focus"}
	COMP_F_DF         = Component{"фок.п./сакр..", "f_df"}
	COMP_M_DF         = Component{"реаг./сакр.", "m_df"}
	COMP_COSTLY       = Component{"ценн.", "costly"}
)

type Spell struct {
	Name, School, ShortDescr string
	CastTime, Duration, Save string
	Range, Area, AreaImg     string
	Descriptors              []string
	Description              template.HTML
	Level                    int
	Components               []Component
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

func GetSavingThrow(str string) string {
	str = strings.Replace(str, "Fortitude", "Стойкость", -1)
	str = strings.Replace(str, "Fort.", "Стойкость", -1)
	str = strings.Replace(str, "Fort", "Стойкость", -1)
	str = strings.Replace(str, "Will", "Воля", -1)
	str = strings.Replace(str, "Reflex", "Рефлекс", -1)
	str = strings.Replace(str, "none", "нет", -1)
	str = strings.Replace(str, "negates", "отменяет", -1)
	str = strings.Replace(str, "partial", "частично", -1)
	str = strings.Replace(str, "half", "наполовину", -1)
	str = strings.Replace(str, "disbelief", "недоверие", -1)
	str = strings.Replace(str, "harmless", "безвредно", -1)
	str = strings.Replace(str, "object", "объект", -1)

	str = strings.Replace(str, "blinding only", "только ослепление", -1)
	str = strings.Replace(str, "if interacted with", "при взаимодействии", -1)

	str = strings.Replace(str, "; see text", "", -1)
	str = strings.Replace(str, "see text", "см. текст", -1)
	str = strings.Replace(str, "or", "или", -1)
	if str == "нет" {
		str = ""
	}
	return str
}

type IntSet map[int]struct{}

func (s *IntSet) Contains(i int) bool {
	_, v := (*s)[i]
	return v
}

func (s *IntSet) Add(i int) {
	(*s)[i] = struct{}{}
}

func ParseIntSet(str string) (IntSet, error) {
	res := make(IntSet)
	for _, p := range strings.Split(str, ",") {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "-") {
			lim := strings.SplitN(p, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(lim[0]))
			if err != nil {
				return res, errors.New("Can't parse int set range start: " +
					err.Error())
			}
			stop, err := strconv.Atoi(strings.TrimSpace(lim[1]))
			if err != nil {
				return res, errors.New("Can't parse int set range stop: " +
					err.Error())
			}
			for i := start; i <= stop; i++ {
				res.Add(i)
			}
		} else {
			i, err := strconv.Atoi(p)
			if err != nil {
				return res, errors.New("Can't parse int set value: " +
					err.Error())
			}
			res.Add(i)
		}
	}
	return res, nil
}

type SpellSortLevel []Spell

func (s SpellSortLevel) Len() int           { return len(s) }
func (s SpellSortLevel) Less(i, j int) bool { return s[i].Level < s[j].Level }
func (s SpellSortLevel) Swap(i, j int)      { tmp := s[i]; s[i] = s[j]; s[j] = tmp }

func GenerateCards(writer io.Writer, class string, level IntSet) error {
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
				if key == col && line[i] != "" {
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
		if !level.Contains(curLevel) || value("source") != "PFRPG Core" {
			continue
		}

		shortDescr := value("short_description_ru", "short_description")
		shortDescr = strings.TrimSpace(shortDescr)
		descr := value("description_formated_ru", "description_formated")
		castTime := T(value("casting_time"))
		if castTime == "основное действие" {
			castTime = ""
		}

		spell := Spell{
			Name:        value("name_ru", "name"),
			School:      value("school"),
			ShortDescr:  shortDescr,
			CastTime:    castTime,
			Duration:    T(value("duration")),
			Save:        GetSavingThrow(value("saving_throw")),
			Range:       T(value("range")),
			Area:        T(value("area")),
			AreaImg:     GetAreaImg(value("range"), value("area")),
			Description: template.HTML(descr),
			Level:       curLevel,
			Components:  GetComponents(value("components")),
		}
		if descStr := value("descriptor"); descStr != "" {
			spell.Descriptors = strings.Split(descStr, ", ")
		}
		data.Spells = append(data.Spells, spell)
	}
	sort.Stable(SpellSortLevel(data.Spells))
	err = tmpl.Execute(writer, data)
	return err
}

func main() {
	var (
		class = flag.String("class", "wiz", "Character class")
		level = flag.String("level", "0", "Spell level")
	)
	flag.Parse()
	levelSet, err := ParseIntSet(*level)
	if err != nil {
		log.Fatalln(err)
	}
	err = GenerateCards(os.Stdout, *class, levelSet)
	if err != nil {
		log.Fatalln(err)
	}
}
