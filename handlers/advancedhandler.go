package handlers

import (
	"net/http"
	"strconv"
	"text/template"
	"workers"
)

func advancedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(advancedTemplatePath))

	MembersNumber := []int{}

	for i := 0; i < 7; i++ {
		MembersNumber = append(MembersNumber, i)
	}

	data := HtmlData{
		ProjectName: "Chazam",
		PageName:    "advanced",

		ASearch_MembersNumber: MembersNumber,
		ASearch_Locations:     workers.GetListOfLocation(),
	}

	if r.Method == "POST" {
		r.ParseForm()
		year, _ := strconv.Atoi(r.Form.Get("date-range"))
		members := make([]int, 0)
		for i := range MembersNumber {
			_, ok := r.Form["inlineCheckbox_"+strconv.Itoa(i)]
			if ok {
				members = append(members, i)
			}
		}
		location := r.FormValue("countries")
		data.Groups = workers.AdvencedFilter(year, members, location)
	} else {
		data.Groups = workers.AdvencedFilter(0, make([]int, 0), "")
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
