package myTime

import (
	"math"
	"time"
)

//DayToInt converte una durata composta da giorni, ore e miuti in minuti
func DayToInt(gg, hh, mm int) int {
	return 1440*gg + 60*hh + mm
}

//Now calcola i minuti passati dall'inizio della settimana
func Now() int {
	date := time.Now()
	return DayToInt(day(date.Weekday().String()), date.Hour(), date.Minute())
}

//Near calcola se le due date differeiscono di n minuti (n=45)
func Near(data1, data2 int) bool {
	return math.Abs(float64(data2-data1)) < 45.0
}

//Distance calcola quanri minuti mancano ad un evento
func Distance(data1, data2 int) int {
	return data2 - data1
}

//day converte un giorno sottoforma di stringa in un intero che lo rappresenta 0-6 Lun-Dom
func day(gg string) int {

	var day int = -1

	switch gg {
	case "Monday":
		day = 0
	case "Tuesday":
		day = 1
	case "Wednesday":
		day = 2
	case "Thursday":
		day = 3
	case "Friday":
		day = 4
	case "Saturday":
		day = 5
	case "Sunday":
		day = 6
	}

	return day
}
