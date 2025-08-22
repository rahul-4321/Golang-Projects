package converter

import(
	"strconv"
	"net/http"
)


func WeightHandler(w http.ResponseWriter, r*http.Request){
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
		result=convertWeight(value,from,to)
	}
	templates.ExecuteTemplate(w,"weight.html",result)
}

func convertWeight(value float64, from string, to string) string{
	//convert all units to grams first
	unitMap:=map[string]float64{
		"milligram":0.001,
		"gram":1.0,
		"kilogram":1000.0,
		"ton":1000000.0,
		"ounce":28.3495,
		"pound":453.592,
	}

	gramValue:=value * unitMap[from]
	result:=gramValue / unitMap[to]	
	return strconv.FormatFloat(result, 'f', 4, 64)
}
