package workers

import (
	"api"
	"objects"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

func buildLocationFromGroupId(groupId int) map[string][]*objects.Location {
	type groupieLocation struct {
		DatesLocations map[string]([]string) `json:"datesLocations"`
	}
	type geoloc struct {
		Long string `json:"lon"`
		Lat  string `json:"lat"`
	}

	locationMap := make(map[string][]*objects.Location)
	var groupieLocationRequest groupieLocation
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(groupId)

	api.GetFromApi(
		url,
		&groupieLocationRequest,
		false,
		time.Millisecond,
		nil)

	for location, dates := range groupieLocationRequest.DatesLocations {
		val, ok := locationMap[location]
		loc := new(objects.Location)

		loc.DateTime = dates
		loc.Name = location

		splited := strings.Split(location, "-")
		cityName := strings.Replace(splited[0], "_", " ", -1)
		countryName := strings.Replace(splited[1], "_", " ", -1)
		url := "https://forward-reverse-geocoding.p.rapidapi.com/v1/forward?format=json&city=" + cityName + "&country=" + countryName
		var geolocRequest []geoloc
		header := map[string]string{
			"X-RapidAPI-Key":  "6bb723da35msh5e7240d75ecfeabp127d5cjsn611360aa1f80",
			"X-RapidAPI-Host": "forward-reverse-geocoding.p.rapidapi.com"}

		api.GetFromApi(
			url,
			&geolocRequest,
			false,
			time.Second,
			header)

		if len(geolocRequest) < 1 {
			url := "https://forward-reverse-geocoding.p.rapidapi.com/v1/forward?state=" + cityName
			api.GetFromApi(
				url,
				&geolocRequest,
				false,
				time.Second,
				header)
		}

		if len(geolocRequest) >= 1 {
			loc.Longitude = geolocRequest[0].Long
			loc.Latitude = geolocRequest[0].Lat
		}

		list := make([]*objects.Location, 0)
		list = append(list, loc)
		if !ok {
			list = append(list, val...)
		}
		mutex.Lock()
		LocationByName[location] = list
		locationMap[location] = list
		mutex.Unlock()
	}

	return locationMap
}

func GetListOfLocation() []string {
	return maps.Keys(LocationByName)
}

func GetDatesFromLocation(location string) []*objects.Location {
	return LocationByName[location]
}
