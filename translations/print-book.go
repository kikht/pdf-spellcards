package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	name = iota
	school
	level
	cast_time
	distance
	target
	duration
	save
	resist
	component
	area
	effect
	description
)

var (
	param_name = []string{
		"Школа", "Круг", "Время сотворения", "Компоненты",
		"Дистанция", "Цель", "Область", "Эффект", "Длит.",
		"Испытание", "Устойчивость к магии",
	}
	param_id = []int{
		school, level, cast_time, component, distance, target, area, effect,
		duration, save, resist,
	}
	param_suffix = []string{";", "", "", "", "", "", "", "", "", "", ""}
)

func main() {
	reader := csv.NewReader(os.Stdin)
	reader.Read()
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("<h3 id=\"%s\">%s</h3>\n",
			strings.Replace(strings.ToLower(r[name]), " ", "-", -1),
			r[name])

		for i := range param_id {
			if r[param_id[i]] != "" {
				fmt.Printf("<p><strong>%s: </strong>%s%s</p>\n",
					param_name[i], r[param_id[i]], param_suffix[i])
			}
		}
		fmt.Println(r[description])
		fmt.Println()
	}
}
