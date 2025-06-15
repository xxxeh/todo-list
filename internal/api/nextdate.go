package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		now = time.Now()
	}
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	date, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(date))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	params := strings.Split(repeat, " ")
	datepart := params[0]

	switch datepart {
	case "y":
		date, err = nextYear(date, now, params)
		if err != nil {
			return "", err
		}

	case "d":
		date, err = nextDay(date, now, params)
		if err != nil {
			return "", err
		}

	case "w":
		date, err = nextDayOfWeek(date, now, params)
		if err != nil {
			return "", err
		}

	case "m":
		date, err = nextDayOfMonth(date, now, params)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("Недопустимый символ %s", datepart)
	}

	return date.Format(dateFormat), nil
}

func nextYear(date, now time.Time, params []string) (time.Time, error) {
	if len(params) > 1 {
		return date, fmt.Errorf("Недопустимый интервал")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			break
		}
	}
	return date, nil
}

func nextDay(date, now time.Time, params []string) (time.Time, error) {
	if len(params) == 1 {
		return date, fmt.Errorf("Не указан интервал")
	}

	days, err := strconv.Atoi(params[1])
	if err != nil {
		return date, err
	}

	if days > 400 {
		return date, fmt.Errorf("Превышен максимально допустимый интервал %d (%d)", 400, days)
	}

	for {
		date = date.AddDate(0, 0, days)
		if date.After(now) {
			break
		}
	}

	return date, nil
}

func nextDayOfWeek(date, now time.Time, params []string) (time.Time, error) {
	if len(params) == 1 {
		return date, fmt.Errorf("Не указан интервал")
	}

	var day [7]bool
	for _, val := range strings.Split(params[1], ",") {
		weekday, err := strconv.Atoi(val)
		if err != nil {
			return date, err
		}

		if weekday < 1 || weekday > 7 {
			return date, fmt.Errorf("Недопустимое значение дня недели %d", weekday)
		}

		if weekday != 7 {
			day[weekday] = true
		} else {
			day[0] = true
		}

	}

	for {
		date = date.AddDate(0, 0, 1)
		if date.After(now) && day[date.Weekday()] {
			break
		}
	}

	return date, nil
}

func nextDayOfMonth(date, now time.Time, params []string) (time.Time, error) {
	if len(params) == 1 {
		return date, fmt.Errorf("Не указан интервал")
	}

	var day [32]bool
	var month [13]bool
	lastDay := 0
	penultDay := 0
	currentMonth := date.Month()
	lastDayOfMonth := time.Date(date.Year(), currentMonth+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)

	for _, val := range strings.Split(params[1], ",") {
		d, err := strconv.Atoi(val)
		if err != nil {
			return date, err
		}
		if d < -2 || d > 31 || d == 0 {
			return date, fmt.Errorf("Недопустимое значение дня %d", d)
		}

		if d == -1 {
			lastDay = lastDayOfMonth.Day()
			day[lastDay] = true
			continue
		}

		if d == -2 {
			penultDay = lastDayOfMonth.Day() - 1
			day[penultDay] = true
			continue
		}

		day[d] = true
	}

	if len(params) > 2 {
		for _, val := range strings.Split(params[2], ",") {
			m, err := strconv.Atoi(val)
			if err != nil {
				return date, err
			}
			if m < 1 || m > 12 {
				return date, fmt.Errorf("Недопустимое значение месяца %d", m)
			}

			month[m] = true
		}
	} else {
		for i := range month {
			month[i] = true
		}
	}

	for {
		if lastDay != 0 || penultDay != 0 {
			if date.Month() != currentMonth {
				currentMonth = date.Month()
				day[lastDay] = false
				day[penultDay] = false
				lastDayOfMonth = time.Date(date.Year(), currentMonth+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
				lastDay = lastDayOfMonth.Day()
				penultDay = lastDayOfMonth.Day() - 1
				day[lastDay] = true
				day[penultDay] = true
			}
		}
		date = date.AddDate(0, 0, 1)
		if date.After(now) && day[date.Day()] && month[date.Month()] {
			break
		}
	}
	return date, nil
}
