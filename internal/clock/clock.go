package clock

import (
	"fmt"
	"time"
)

type Info struct {
	Time string
	Date string
}

var weekdays = [...]string{
	"domingo",
	"lunes",
	"martes",
	"miércoles",
	"jueves",
	"viernes",
	"sábado",
}

var months = [...]string{
	"enero",
	"febrero",
	"marzo",
	"abril",
	"mayo",
	"junio",
	"julio",
	"agosto",
	"septiembre",
	"octubre",
	"noviembre",
	"diciembre",
}

func Now() Info {
	now := time.Now()

	date := weekdays[now.Weekday()] +
		", " +
		fmt.Sprintf("%d", now.Day()) +
		" de " +
		months[now.Month()-1]

	return Info{
		Time: now.Format("15:04:05"),
		Date: date,
	}
}
