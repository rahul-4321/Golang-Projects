package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/length", http.StatusSeeOther)
}

// Length
func lengthHandler(w http.ResponseWriter, r *http.Request) {
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
		result = convertLength(value, from, to)
	}
	templates.ExecuteTemplate(w, "length.html", result)
}

func convertLength(value float64, from string, to string) string {
	unitMap := map[string]float64{
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1.0,
		"kilometer":  1000.0,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.34,
	}

	meterValue := value * unitMap[from]
	result := meterValue / unitMap[to]
	return strconv.FormatFloat(result, 'f', 4, 64)
}

// Weight
func weightHandler(w http.ResponseWriter, r *http.Request) {
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
		result = convertWeight(value, from, to)
	}
	templates.ExecuteTemplate(w, "weight.html", result)
}

func convertWeight(value float64, from string, to string) string {
	//convert all units to grams first
	unitMap := map[string]float64{
		"milligram": 0.001,
		"gram":      1.0,
		"kilogram":  1000.0,
		"ton":       1000000.0,
		"ounce":     28.3495,
		"pound":     453.592,
	}

	gramValue := value * unitMap[from]
	result := gramValue / unitMap[to]
	return strconv.FormatFloat(result, 'f', 4, 64)
}

// Temperature
func tempHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/length", lengthHandler)
	http.HandleFunc("/weight", weightHandler)
	http.HandleFunc("/temp", tempHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
