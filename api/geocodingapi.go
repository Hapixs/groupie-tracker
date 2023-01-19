package api

import (
	"strings"
	"time"
)

type Geoloc struct {
	Long string `json:"lon"`
	Lat  string `json:"lat"`
}

func GetInformationForCity(city string) Geoloc {
	var request []Geoloc
	splited := strings.Split(city, "-")
	cityName := strings.Replace(splited[0], "_", " ", -1)
	countryName := strings.Replace(splited[1], "_", " ", -1)
	url := "https://forward-reverse-geocoding.p.rapidapi.com/v1/forward?format=json&city=" + cityName + "&country=" + countryName
	GetFromApi(url, &request, false, time.Second, map[string]string{
		"X-RapidAPI-Key":  "6bb723da35msh5e7240d75ecfeabp127d5cjsn611360aa1f80",
		"X-RapidAPI-Host": "forward-reverse-geocoding.p.rapidapi.com"})
	if len(request) < 1 {
		return GetInformationForState(cityName)
	}
	return request[0]
}

func GetInformationForState(state string) Geoloc {
	var request []Geoloc
	url := "https://forward-reverse-geocoding.p.rapidapi.com/v1/forward?state=" + state
	GetFromApi(url, &request, false, time.Second, map[string]string{
		"X-RapidAPI-Key":  "6bb723da35msh5e7240d75ecfeabp127d5cjsn611360aa1f80",
		"X-RapidAPI-Host": "forward-reverse-geocoding.p.rapidapi.com"})
	if len(request) < 1 {
		return Geoloc{}
	}
	return request[0]
}
