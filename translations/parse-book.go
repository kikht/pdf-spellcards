package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	header_re   = regexp.MustCompile(`<h3 id=".*">(.*)</h3>`)
	property_re = regexp.MustCompile(`<p><strong>(.*)</strong>(.*)</p>`)
	//school_re = regexp.MustCompile(`<p><strong>Школа: </strong>(.*);</p>`)
)

type Spell struct {
	Name      string
	School    string
	Level     string
	CastTime  string
	Dist      string
	Target    string
	Duration  string
	Save      string
	Resist    string
	Descr     []string
	Component string
	Area      string
	Effect    string
}

type parseState int

const (
	PARSE_NAME = iota
	PARSE_PROPERTY
	PARSE_DESCRIPTION
)

var lastSpell string = "begin"

func printSpell(s Spell, writer *csv.Writer) {
	if s.Name != "" {
		writer.Write([]string{s.Name, s.School, s.Level, s.CastTime,
			s.Dist, s.Target, s.Duration, s.Save, s.Resist, s.Component,
			s.Area, s.Effect, strings.Join(s.Descr, "\n")})
		lastSpell = s.Name
	} else {
		log.Println("Strange unnamed spell", s, "after", lastSpell)
	}
}

func main() {
	writer := csv.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	writer.Write([]string{"name", "school", "level", "cast_time",
		"distance", "target", "duration", "save", "resist", "component",
		"area", "effect", "description"})
	var s Spell
	var state = PARSE_DESCRIPTION
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if p := header_re.FindStringSubmatch(line); p != nil {
			printSpell(s, writer)
			s = Spell{}
			s.Name = strings.TrimSpace(p[1])

			if state != PARSE_DESCRIPTION {
				log.Println(s.Name, " - Unexpected name in state", state)
			}
			state = PARSE_NAME
		} else if p := property_re.FindStringSubmatch(line); p != nil {
			v := strings.TrimSpace(p[2])
			switch p[1] {
			case "Школа: ":
				s.School = strings.TrimRight(v, ";")
			case "Круг: ":
				s.Level = v
			case "Время сотворения: ":
				s.CastTime = v
			case "Дистанция: ":
				s.Dist = v
			case "Цель: ", "Цели: ":
				s.Target = v
			case "Цель или область: ", "Область или цель: ":
				s.Target = v
				s.Area = v
			case "Цель, эффект или область: ":
				s.Target = v
				s.Area = v
				s.Effect = v
			case "Длит.: ", "Длительность: ":
				s.Duration = v
			case "Испытание: ":
				s.Save = v
			case "Устойчивость к магии: ":
				s.Resist = v
			case "Компоненты: ":
				s.Component = v
			case "Область: ":
				s.Area = v
			case "Эффект: ":
				s.Effect = v
			default:
				log.Println("Unparsed property:", p[1:], "for spell", s.Name)
			}

			if state != PARSE_NAME && state != PARSE_PROPERTY {
				log.Println(s.Name, " - Unexpected property",
					p[1], "=", v, "in state", state)
			}
			state = PARSE_PROPERTY
		} else if (strings.HasPrefix(line, "<p>") && strings.HasSuffix(line, "</p>")) ||
			strings.HasPrefix(line, "<ol") || line == "</ol>" ||
			(strings.HasPrefix(line, "<li>") && strings.HasSuffix(line, "</li>")) {
			s.Descr = append(s.Descr, line)

			if state != PARSE_PROPERTY && state != PARSE_DESCRIPTION {
				log.Println(s.Name, " - Unexpected description in state", state)
			}
			state = PARSE_DESCRIPTION
		} else if line == "" {
		} else {
			log.Println("Unparsed line:", line)
		}
	}
	printSpell(s, writer)
	writer.Flush()
}
