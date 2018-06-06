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
	"7 rounds":      "7 раундов",

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

	"1 hour":             "1 час",
	"1 hour or less":     "1 час или менее",
	"2 hours":            "2 часа",
	"4 hours":            "4 часа",
	"6 hours":            "6 часов",
	"8 hours":            "8 часов",
	"12 hours":           "12 часов",
	"24 hours":           "24 часа",
	"24 hours; see text": "24 часа, см. текст",
	"1 day":              "1 день",
	"1 week":             "1 неделя",

	"see text": "см. текст",
	"See text": "см. текст",

	//Duration
	"instant":                 "мгновенное",
	"instantaneous":           "мгновенное",
	"Instantaneous":           "мгновенное",
	"instantaneous; see text": "мгновенное, см. текст",

	"instantaneous (1 round); see text":        "мгновенное (1 раунд); см. текст",
	"instantaneous (1d4 rounds); see text":     "мгновенное (1d4 раундов); см. текст",
	"instantaneous (1d6 rounds); see text":     "мгновенное (1d6 раундов); см. текст",
	"instantaneous or 1 round/level; see text": "мгновенное или 1 раунд/УЗ; см. текст",
	"instantaneous or 10 min./level; see text": "мгновенное или 10 мин./УЗ; см. текст",
	"instantaneous/1 hour; see text":           "мгновенное / 1 час; см. текст",

	"instantaneous/10 minutes per HD of subject; see text":           "мгновенное / 10 минут за КЗ существа; см. текст",
	"instantaneous or concentration (up to 1 round/level); see text": "мгновенное или концентрация (до 1 раунда/УЗ); см. текст",

	"permanent":                      "постоянное",
	"permanent (D)":                  "постоянное (до отмены)",
	"permanent until discharged":     "постоянное (до отмены)",
	"permanent until discharged (D)": "постоянное (до отмены)",
	"permanent; see text":            "постоянное; см. текст",

	"permanent or until discharged; until released or 1d4 days + 1 day/level; see text": "постоянное или пока не сработает; до прекращения или 1d4 дня + 1 день / 1 УЗ; см. текст",
	"permanent until triggered, then 1 round/level":                                     "постоянное, пока не сработает, далее 1 раунд / 1 УЗ",

	"concentration":                                                  "концентрация",
	"concentration, up to 1 round/level":                             "концентрация, 1 раунд/УЗ",
	"concentration, up to 1 minute/level":                            "концентрация, 1 мин./УЗ",
	"concentration, up to 1 minute/ level":                           "концентрация, 1 мин./УЗ",
	"concentration, up to 1 min./level":                              "концентрация, 1 мин./УЗ",
	"concentration, up to 10 min./level":                             "концентрация, 10 мин./УЗ",
	"concentration, up to 10 min./ level":                            "концентрация, 10 мин./УЗ",
	"concentration (maximum 10 rounds) (D)":                          "концентрация (максимум 10 раундов) (П)",
	"concentration (up to 1 round per 2 levels)":                     "концентрация (до 1 раунда за 2 УЗ)",
	"concentration (up to 1 round/level) or instantaneous; see text": "концентрация (до 1 раунда / 1 УЗ) или мгновенное; см. текст",
	"concentration + 1 hour/level":                                   "концентрация + 1 час / 1 УЗ (П)",
	"concentration + 1 round/level":                                  "концентрация + 1 раунд / 1 УЗ",
	"concentration + 2 rounds":                                       "концентрация + 2 раунда",
	"Concentration + 2 rounds":                                       "концентрация + 2 раунда",
	"Concentration + 3 rounds":                                       "концентрация + 3 раунда",
	"concentration +1 hour/ level":                                   "концентрация +1 час / 1 УЗ (П)",
	"Concentration +1 round/level":                                   "концентрация +1 раунд / 1 УЗ (П)",

	"1 round/level":                      "1 раунд/УЗ",
	"1 round/level (D)":                  "1 раунд/УЗ или до отмены",
	"1 round/level ; see text":           "1 раунд/УЗ, см. текст",
	"1 round + 1 round per three levels": "1 раунд + 1 раунд/3 УЗ",
	"3 rounds/level":                     "3 раунда/УЗ",

	"1 minute or until discharged": "1 минута или до отмены",
	"1 min./level":                 "1 минута/УЗ",
	"1 min./level (D)":             "1 мин./УЗ или до отмены",
	"1 min/level":                  "1 минута/УЗ",
	"1 minute/level":               "1 минута/УЗ",
	"2 min./level (D)":             "2 мин. / 1 УЗ (П)",

	"10 min./level":                     "10 минут/УЗ",
	"10 min./level (D)":                 "10 минут/УЗ или до отмены",
	"10 min./level or until discharged": "10 минут/УЗ или до отмены",
	"10 min./level or until used":       "10 минут/УЗ или пока не будет использовано",

	"30 minutes and 2d6 rounds; see text":        "30 минут и 2d6 раундов; см. текст",
	"30 minutes or until discharged":             "30 минут или до отмены",
	"1d4 rounds or 1 round; see text":            "1d4 раунда или 1 раунд, см. текст",
	"2d4 rounds":                                 "2d4 раунда",
	"1d6+2 rounds":                               "1d6+2 раунда",
	"1 hour/level":                               "1 час/УЗ",
	"1 hour/level (D)":                           "1 час/УЗ или до отмены",
	"1 hour/level or until discharged; see text": "1 час/УЗ или до отмены; см. текст",
	"2 hours/level":                              "2 часа/УЗ",
	"2 hours/level; see text":                    "2 часа/УЗ; см. текст",
	"1 day/level":                                "1 день/УЗ",
	"1 day/level; see text":                      "1 день/УЗ; см. текст",
	"1 day/level or until discharged":            "1 день/УЗ или до отмены",
	"1 day/level (D) or until discharged":        "1 день/УЗ или до отмены",

	"1 hour plus 12 hours; see text":                "1 час плюс 12 часов; см. текст",
	"1 hour/level ; see text":                       "1 час/УЗ; см. текст",
	"1 hour/level or until completed":               "1 час/УЗ или пока не будет исполнено",
	"1 hour/level or until discharged":              "1 час/УЗ или пока не сработает",
	"1 hour/level or until expended; see text":      "1 час/УЗ или пока не будет израсходовано; см. текст",
	"1 hour/level or until you return to your body": "1 час/УЗ или пока вы не вернетесь в свое тело",
	"1 hour/level; see text":                        "1 час/УЗ; см. текст",

	"1 round /level (D)":                                       "1 раунд / 1 УЗ (П)",
	"1 round/level (D) and concentration + 3 rounds; see text": "1 раунд / 1 УЗ (П) и концентрация + 3 раунда; см. текст",
	"1 round/level or 1 hour/level; see text":                  "1 раунд / 1 УЗ или 1 час / 1 УЗ; см. текст",
	"1 round/level or 1 round; see text":                       "1 раунд / 1 УЗ или 1 раунд; см. текст",
	"1 round/level or 1 round; see text for cause fear":        "1 раунд / 1 УЗ или 1 раунд; см. текст заклинания устрашение",
	"1 round/level or until all beams are exhausted":           "1 раунд / 1 УЗ или пока не кончатся все лучи",
	"1 round/level or until discharged, whichever comes first": "1 раунд / 1 УЗ или пока не сработает",

	"1d4+1 rounds":                           "1d4+1 раунд",
	"1d4+1 rounds (apparent time); see text": "1d4+1 раунд; см. текст",
	"4d12 hours; see text":                   "4d12 часов; см. текст",
	"5 rounds or less; see text":             "5 раундов или меньше; см. текст",
	"60 days or until discharged":            "60 дней или пока не сработает",
	"7 days or 7 months ; see text":          "7 дней или 7 месяцев (П); см. текст",

	"one day/level":                   "один день / 1 УЗ (П)",
	"one usage per two levels":        "одно применение за 2 УЗ",
	"until expended or 10 min./level": "пока не будет исчерпано или 10 мин. / 1 УЗ",
	"until landing or 1 round/level":  "до приземления или 1 раунд / 1 УЗ",
	"Until triggered or broken":       "пока не сработает или не будет разрушено",
	"up to 1 round/level":             "до 1 раунда / 1 УЗ",

	"1 hour/caster level or until discharged, then 1 round/caster level; see text":  "1 час/УЗ или пока не сработает, затем 1 раунд / 1 УЗ; см. текст",
	"no more than 1 hour/level or until discharged (destination is reached)":        "пока не достигнет места назначения, но не более 1 часа / 1 УЗ",
	"1d4+1 rounds, or 1d4+1 rounds after creatures leave the smoke cloud; see text": "1d4+1 раунд или 1d4+1 раунд после того, как существа выйдут из дыма; см. текст",

	//Range
	"personal":           "на себя",
	"touch":              "касание",
	"touch; see text":    "касание, см. текст",
	"personal or touch":  "на себя или касание",
	"personal and touch": "на себя и касание",
	"unlimited":          "без ограничений",

	"close (25 ft. + 5 ft./2 levels)": "близкая (25 фт. + 5 фт./2 УЗ)",
	"medium (100 ft. + 10 ft./level)": "средняя (100 фт. + 10 фт./УЗ)",
	"medium (100 ft. + 10 ft. level)": "средняя (100 фт. + 10 фт./УЗ)",
	"long (400 ft. + 40 ft./level)":   "дальняя (400 фт. + 40 фт./УЗ)",

	"personal or close (25 ft. + 5 ft./2 levels)":       "на себя или близкая (25 фт. + 5 фт./2 УЗ)",
	"close (25 ft. + 5 ft./2 levels) or see text":       "близкая (25 футов + 5 футов / 2 УЗ) или см. текст",
	"close (25 ft. + 5 ft./2 levels); see text":         "близкая (25 футов + 5 футов / 2 УЗ); см. текст",
	"close (25 ft. + 5 ft./2 levels)/100 ft.; see text": "близкая (25 футов + 5 футов / 2 УЗ) / 100 футов; см. текст",

	"0 ft.":           "0 футов",
	"0 ft.; see text": "0 футов; см. текст",
	"10 ft.":          "10 футов",
	"15 ft.":          "15 футов",
	"20 ft.":          "20 футов",
	"30 ft.":          "30 футов",
	"40 ft.":          "40 футов",
	"50 ft.":          "50 футов",
	"60 ft.":          "60 футов",
	"120 ft.":         "120 футов",

	"1 mile":  "1 миля",
	"2 miles": "2 мили",
	"5 miles": "5 миль",

	"up to 10 ft./level":                    "до 10 футов / 1 УЗ",
	"40 ft./level":                          "40 футов / 1 УЗ",
	"1 mile/level":                          "1 миля / 1 УЗ",
	"anywhere within the area to be warded": "внутри защищаемой области",

	//Area
	"cone-shaped emanation":                "коническая эманация",
	"cone-shaped burst":                    "конический всплеск",
	"50-ft.-radius burst, centered on you": "50 футов вокруг колдующего",

	"20-ft.-radius emanation centered on a point in space":                      "20 футовая сфера с центром в заданной точке",
	"one creature, one object, or a 5-ft. cube":                                 "одно создание, один объект или 5 футовый куб",
	"The caster and all allies within a 50-ft. burst, centered on the caster":   "Колдующий и союзники в пределах 50 футов",
	"one or more living creatures within a 10-ft.-radius burst":                 "одно или несколько живых созданий внутри 10 футовой сферы",
	"several living creatures, no two of which may be more than 30 ft. apart":   "несколько живых созданий, никакие два из которых не стоят дальше 30 футов друг от друга",
	"10-ft. square/level; see text":                                             "10 кв. футов / 1 УЗ; см. текст",
	"10-ft.-radius emanation around the creature":                               "эманация радиусом 10 футов вокруг существа",
	"10-ft.-radius emanation centered on you":                                   "эманация радиусом 10 футов, с центром на вас",
	"10-ft.-radius emanation from touched creature":                             "эманация радиусом 10 футов вокруг существа",
	"10-ft.-radius emanation, centered on you":                                  "эманация радиусом 10 футов, с центром на вас",
	"10-ft.-radius spherical emanation, centered on you":                        "эманация радиусом 10 футов, с центром на вас",
	"10-ft.-radius spread":                                                      "облако радиусом 10 футов",
	"120-ft. line":                                                              "линия 120 футов",
	"2-mile-radius circle, centered on you; see text":                           "круг радиусом 2 мили с центром на вас; см. текст",
	"20-ft.-radius burst":                                                       "взрыв радиусом 20 футов",
	"20-ft.-radius emanation":                                                   "эманация радиусом 20 футов",
	"20-ft.-radius emanation centered on a creature, object, or point in space": "эманация радиусом 20 футов, центром которой является существо, предмет или точка в пространстве",
	"20-ft.-radius spread":                                                      "облако радиусом 20 футов",
	"30-ft. cube/level":                                                         "(куб 30×30×30 футов) / 1 УЗ",
	"40 ft./level radius cylinder 40 ft. high":                                  "цилиндр радиусом 40 футов / 1 УЗ и 40 футов высотой",
	"40-ft. radius emanating from the touched point":                            "эманация радиусом 40 футов",
	"40-ft.-radius emanation":                                                   "эманация радиусом 40 футов",
	"40-ft.-radius emanation centered on you":                                   "эманация радиусом 40 футов, с центром на вас",
	"5-ft.-radius emanation centered on you":                                    "эманация радиусом 5 футов с центром на вас",
	"5-ft.-radius spread; or one solid object or one crystalline creature":      "облако радиусом 5 футов, либо один твердый предмет, либо одно существо с кристаллической структурой",
	"60-ft. cube/level":                                                         "(куб 60×60×60) / 1 УЗ",
	"60-ft. line from you":                                                      "60-футовая линия, идущая от вас",
	"60-ft. line-shaped emanation from you":                                     "60-футовая эманация в форме линии, идущая от вас",
	"80-ft.-radius burst":                                                       "взрыв радиусом 80 футов",
	"80-ft.-radius spread (S)":                                                  "облако радиусом 80 футов (Ф)",
	"all allies and foes within a 40-ft.-radius burst centered on you":          "все союзники и противники в области радиусом 40 футов с вами в центре (взрыв)",
	"all metal objects within a 40-ft.-radius burst":                            "все металлические предметы в радиусе 40 футов (взрыв)",
	"barred cage (20-ft. cube) or windowless cell (10-ft. cube)":                "решетчатая клетка (куб 20×20×20 футов) или клетка без отверстий (куб 10×10×10 футов)",
	"circle, centered on you, with a radius of 400 ft. + 40 ft./level":          "круг с вами в центре и радиусом 400 футов + 40 футов / 1 УЗ",
	"cloud spreads in 20-ft. radius, 20 ft. high":                               "облако радиусом 20 футов, высотой 20 футов",
	"creatures and objects within 10-ft.-radius spread":                         "существа и предметы в радиусе 10 футов (облако)",
	"creatures and objects within a 5-ft.-radius burst":                         "существа и предметы в радиусе 5 футов (взрыв)",
	"creatures in a 20-ft.-radius spread":                                       "существа в облаке радиусом 20 футов",
	"creatures within a 20-ft.-radius spread":                                   "существа в радиусе 20 футов (облако)",
	"cylinder (10-ft. radius, 40-ft. high)":                                     "цилиндр радиусом 10 футов, высотой 40 футов",
	"cylinder (20-ft. radius, 40 ft. high)":                                     "цилиндр радиусом 20 футов, высотой 40 футов",
	"cylinder (40-ft. radius, 20 ft. high)":                                     "цилиндр радиусом 40 футов, высотой 20 футов",
	"dirt in an area up to 750 ft. square and up to 10 ft. deep (S)":            "земля в области до 750×750 футов, глубиной до 10 футов (Ф)",
	"four 40-ft.-radius spreads, see text":                                      "четыре облака радиусом 40 футов, см. текст",
	"line from your hand":                                                       "линия, идущая от вашей руки",
	"living creatures within a 10-ft.-radius burst":                             "живые существа в радиусе 10 футов (взрыв)",
	"nonchaotic creatures in a 40-ft.-radius spread centered on you":            "нехаотичные существа в радиусе 40 футов (облако) с вами в центре",
	"nonevil creatures in a 40-ft.-radius spread centered on you":               "все незлые существа в радиусе 40 футов от вас",
	"nongood creatures in a 40-ft.-radius spread centered on you":               "недобрые существа в радиусе 40 футов от вас (облако)",
	"nonlawful creatures in a 40-ft.-radius spread centered on you":             "непринципиальные существа в радиусе 40 футов от вас",
	"nonlawful creatures within a burst that fills a 30-ft. cube":               "непринципиальные существа в кубе 30×30×30 футов (взрыв)",
	"object touched or up to 5 sq. ft./level":                                   "один предмет или до 5 кв. футов / 1 УЗ",
	"one 20-ft. cube/level":                                                     "(один куб 20×20×20 футов) / 1 УЗ (Ф)",
	"one 20-ft. cube/level (S)":                                                 "(один куб 20×20×20 футов) / 1 УЗ (Ф) или один огненный волшебный предмет",
	"one 20-ft. square/level":                                                   "(один квадрат 20×20 футов) / 1 УЗ",
	"one 30-ft. cube/level":                                                     "(один куб 30×30×30 футов) / 1 УЗ (Ф)",
	"one spellcaster, creature, or object":                                      "один заклинатель, существо или предмет",
	"plants in a 40-ft.-radius spread":                                          "растения в облаке радиусом 40 футов",
	"several living creatures within a 40-ft.-radius burst":                     "несколько живых существ в радиусе 40 футов",
	"several undead creatures within a 40-ft.-radius burst":                     "несколько представителей нежити в радиусе 40 футов (взрыв)",
	"Target see text":                                                           "см. текст",
	"two 10-ft. cubes per level (S)":                                            "(два куба 10×10×10 футов) / 1 УЗ (Ф)",
	"up to 10-ft.-radius/level emanation centered on you":                       "эманация с вами в центре радиусом до 10 футов / 1 УЗ",
	"up to 200 sq. ft./level":                                                   "до 200 кв. футов / 1 УЗ (Ф)",
	"up to one 10-ft. cube/level (S)":                                           "(один куб 10×10×10 футов) / 1 УЗ (Ф)",
	"up to two 10-ft. cubes/level":                                              "(до двух кубов 10×10×10 футов) / 1 УЗ (Ф)",
	"water in a volume of 10 ft./level by 10 ft./level by 2 ft./level":          "вода объемом 10 футов / 1 УЗ на 10 футов / 1 УЗ на 2 фута / 1 УЗ (Ф)",

	"all magical effects and magic items within a 40-ft.-radius burst, or one magic item (see text)": "все магические эффекты и предметы в радиусе 40 футов либо один магический предмет (см. текст)",
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
)

type Spell struct {
	Name, School, ShortDescr string
	CastTime, Duration, Save string
	Range, Area, AreaImg     string
	Descriptors              []string
	Description              template.HTML
	Level                    int
	Components               []Component
	ComponentsText           string
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

func GetComponentsText(text, costly string) string {
	if costly != "1" {
		return ""
	}
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "С, ")
	text = strings.TrimPrefix(text, "Ж, ")
	return text
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
			ComponentsText: GetComponentsText(
				value("components_ru", "components"),
				value("costly_components")),
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
