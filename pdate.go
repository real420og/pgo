package pgo

import (
	"time"
	"regexp"
	"strconv"
)

var specCases map[string]interface{}

type GoDate struct {
	parsedSymbols   []string
	inputDateFormat string
	t               time.Time
	unix            int64
	phpToGoFormat   map[string]string
}

const phpDateFormatSymbols = "[\\D]"

func Date(args ...interface{}) string {
	var date GoDate

	if len(args) == 0 {
		panic("At least 1st parameter format must be passed")
	}

	for i, p := range args {
		switch i {
		case 0:
			param, ok := p.(string)
			isOk(ok, "You must provide format string parameter as string")
			date.inputDateFormat = param
			break
		case 1:
			param, ok := p.(int64)
			isOk(ok, "You must provide timestamp as int64 type")
			date.unix = param
			break
		}
	}

	date.t = time.Now()
	if date.unix > 0 {
		// parse unix representation with int type to Time object
		date.t = time.Unix(date.unix, 0)
	}

	date.initMapping()
	return date.t.Format(date.parse())
}

func (date *GoDate) parse() string {
	date.convert()

	var convertedString string
	for _, v := range date.parsedSymbols {
		if val, ok := date.phpToGoFormat[v]; ok {
			convertedString += val
		} else if sVal, ok := specCases[v]; ok {
			v, ok := sVal.(int)
			if ok {
				convertedString += strconv.Itoa(v)
			} else {
				convertedString += sVal.(string)
			}
		} else {
			convertedString += v
		}
	}

	return convertedString
}

func (date *GoDate) convert() {
	r, _ := regexp.Compile(phpDateFormatSymbols)

	for _, chSlice := range r.FindAllStringSubmatch(date.inputDateFormat, -1) {
		for _, ch := range chSlice {
			date.parsedSymbols = append(date.parsedSymbols, ch)
		}
	}
}

func (date *GoDate) initMapping() {
	date.phpToGoFormat = map[string]string{
		"Y": "2006",
		"m": "02",
		"d": "01",
		"H": "15",
		"i": "04",
		"s": "05",
		"D": "Mon",
		"M": "Jan",
	}

	_, isoWeek := date.t.ISOWeek();
	specCases = map[string]interface{}{
		"l": date.t.Weekday().String(),
		// todo: escape double conversion for week day and day, day of year
		"N": isoWeek,
		"z": date.t.YearDay(),
		"j": date.t.Day(),
	}
}