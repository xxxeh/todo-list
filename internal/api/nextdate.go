package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nextDateHandler обрабатывает запрос вычисление следующей даты повторения задачи.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(dateFormat, r.FormValue("now"))
	if err != nil {
		now = time.Now()
	}
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	date, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(date))
}

// NextDate вычисляет дату следующего повторения задачи.
// Параметры:
//
//	now - текущая дата.
//	date - дата выполнения задачи.
//	repeat - правило повторения задачи.
//
// Возвращаемые значения:
//
//	string - рассчитанная относительно правила, следующая дата повторения задачи.
//
//	error - ошибка, которая могла возникнуть в ходе работы.
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

// nextYear рассчитывает следющую дату, если задача выполняется ежегодно.
// Параметры:
//
//	date - дата выполнения задачи.
//	now - текущая дата.
//	params - интервал повторения задачи.
//
// Возвращаемые значения:
//
//	time.Time - следующая дата повторения задачи.
//	error - ошибка, которая могла возникнуть в ходе работы.
func nextYear(date, now time.Time, params []string) (time.Time, error) {
	if len(params) > 1 {
		return date, fmt.Errorf("Недопустимый интервал")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if after(date, now) {
			break
		}
	}
	return date, nil
}

// nextDay рассчитывает следющую дату, если задача выполняется с интервалом указанным в днях.
// Параметры:
//
//	date - дата выполнения задачи.
//	now - текущая дата.
//	params - интервал повторения задачи.
//
// Возвращаемые значения:
//
//	time.Time - следующая дата повторения задачи.
//	error - ошибка, которая могла возникнуть в ходе работы.
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
		if after(date, now) {
			break
		}
	}

	return date, nil
}

// nextDay рассчитывает следющую дату, если задача выполняется с интервалом указанным в днях недели.
// Параметры:
//
//	date - дата выполнения задачи.
//	now - текущая дата.
//	params - интервал повторения задачи.
//
// Возвращаемые значения:
//
//	time.Time - следующая дата повторения задачи.
//	error - ошибка, которая могла возникнуть в ходе работы.
func nextDayOfWeek(date, now time.Time, params []string) (time.Time, error) {
	if len(params) == 1 {
		return date, fmt.Errorf("Не указан интервал")
	}

	var day [7]bool
	//Парсим переданный интервал и отмечаем "true" дни недели в массиве day.
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
			//Если weekday = 7, то отмечаем нулевой элемент массива day.
			//Связано с особенностями нумерации дней недели (в date.Weekday() 0 - воскресенье 1 - понедельник,
			// 												в params 	     7 - воскресенье 1 - понедельник).
			day[0] = true
		}

	}

	for {
		//Двигаемся с шагом в один день.
		date = date.AddDate(0, 0, 1)
		if after(date, now) && day[date.Weekday()] {
			//Если день недели имеет значение true и дата стала позже текущей, то значит следующая дата найдена, выходим из цикла.
			break
		}
	}

	return date, nil
}

// nextDay рассчитывает следющую дату, если задача выполняется с интервалом указанным в днях месяца.
// Интервал может быть указан как в виде дней месяца, так и с указанием конкретных месяцев.
// Параметры:
//
//	date - дата выполнения задачи.
//	now - текущая дата.
//	params - интервал повторения задачи.
//
// Возвращаемые значения:
//
//	time.Time - следующая дата повторения задачи.
//	error - ошибка, которая могла возникнуть в ходе работы.
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

	//Парсим переданный интервал и отмечаем дни "true" в массиве day.
	for _, val := range strings.Split(params[1], ",") {
		d, err := strconv.Atoi(val)
		if err != nil {
			return date, err
		}
		if d < -2 || d > 31 || d == 0 {
			return date, fmt.Errorf("Недопустимое значение дня %d", d)
		}

		//Если в интервале указаны дни в виде -1, -2 - рассчитываем последний и предпоследний день текущего месяца.
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

	//Если в интервале указаны и конкретные месяцы повторения задачи, то отмечаем их "true" в массиве month.
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
		//Если в интервале не указаны конкретные месяцы, то отметим как подходящие все месяцы.
		for i := range month {
			month[i] = true
		}
	}

	for {
		if lastDay != 0 || penultDay != 0 {
			if date.Month() != currentMonth {
				//Если задача должна повторятся в последний и/или предпоследний день месяца и вычисление уже перешло на следующий месяц, то
				//Пересчитываем последний и/или предпоследний день и обновляем значения массива day.
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
		//Двигаемся с шагом в один день.
		date = date.AddDate(0, 0, 1)
		if after(date, now) && day[date.Day()] && month[date.Month()] {
			//Если день и месяц имеют значение true и дата стала позже текущей, то значит следующая дата найдена, выходим из цикла.
			break
		}
	}
	return date, nil
}
