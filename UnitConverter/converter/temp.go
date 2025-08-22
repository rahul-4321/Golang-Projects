package converter

import (
	"net/http"
	"strconv"
)

func TempHandler(w http.ResponseWriter, r *http.Request) {
	var result string

	if r.Method == http.MethodPost {
		valuestr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")

		value, err := strconv.ParseFloat(valuestr, 64)
		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}
		result = convertTemperature(value, from, to)
	}
	templates.ExecuteTemplate(w, "temp.html", result)
}

func convertTemperature(value float64, from string, to string) string {

	var celcius float64
	switch from {
	case "Celsius":
		celcius = value
	case "Fahrenheit":
		celcius = (value - 32) * 5 / 9
	case "Kelvin":
		celcius = value - 273.15
	default:
		return "Invalid unit"
	}
	var result float64
	switch to {
	case "Celsius":
		result = celcius
	case "Fahrenheit":
		result = (celcius * 9 / 5) + 32
	case "Kelvin":
		result = celcius + 273.15
	default:
		return "Invalid unit"
	}
	return strconv.FormatFloat(result, 'f', 4, 64)

}
