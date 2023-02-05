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

	for k, v := range groupieLocationRequest.DatesLocations {
		val, ok := locationMap[k]
		l := []*objects.Location{new(objects.Location)}
		l[0].DateTime = v
		l[0].Name = k

		splited := strings.Split(k, "-")
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

		if len(geolocRequest) > 1 {
			l[0].Logitude = geolocRequest[0].Long
			l[0].Latitude = geolocRequest[0].Lat
		}

		if !ok {
			l = append(l, val...)
		}
		mutex.Lock()
		LocationByName[k] = l // pas sure que ca soit bon ca mon pote
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
