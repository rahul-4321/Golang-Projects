package converter

import (
	"strconv"
	"net/http"
)

//Length
func LengthHandler(w http.ResponseWriter, r*http.Request){
	 var result string

	if r.Method==http.MethodPost{
		valuestr:=r.FormValue("value")
		from:=r.FormValue("from")
		to:=r.FormValue("to")

		value,err:=strconv.ParseFloat(valuestr,64)
		if err!=nil{
			http.Error(w,"Invalid value",http.StatusBadRequest)
			return
		}
		result=convertLength(value,from,to)
	}
	templates.ExecuteTemplate(w,"length.html",result)
}

func convertLength(value float64, from string, to string) string{
	unitMap:=map[string]float64{
		"millimeter":0.001,
		"centimeter":0.01,
		"meter":1.0,
		"kilometer":1000.0,
		"inch":0.0254,
		"foot":0.3048,
		"yard":0.9144,
		"mile":1609.34,
	}

	meterValue:=value * unitMap[from]
	result:=meterValue / unitMap[to]	
	return strconv.FormatFloat(result, 'f', 4, 64)
}
