package utils

import (
	"bytes"
	"text/template"
)

func LoadFragmentAsString(name string, data any) string {
	tmpl := template.Must(template.ParseFiles("static/templates/fragments/" + name))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, data)
	return buf.String()
}
