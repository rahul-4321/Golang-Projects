package converter

import "html/template"

var templates * template.Template

func InitTemplates(t *template.Template){
	templates=t
}