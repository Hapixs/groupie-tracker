package workers

import (
	"api"
	"objects"

	"golang.org/x/exp/maps"
)

func AddGroupDates(group objects.Group) {
	mutex.Lock()
	for _, v := range group.DateLocations {
		for _, date := range v {
			date.GroupId = group.Id
			dates := []api.Date{date}
			val, ok := Locations[date.Locations]
			if ok {
				dates = append(dates, val...)
			}
			Locations[date.Locations] = dates
		}
	}
	mutex.Unlock()
}

func GetListOfLocation() []string {
	return maps.Keys(Locations)
}

func GetDatesFromLocation(location string) []api.Date {
	return Locations[location]
}
