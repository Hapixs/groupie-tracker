package workers

import (
	"api"
	"objects"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var GroupMap = map[int](objects.Group){}
var GroupByGenreMap = map[api.DeezerGenre]([]objects.Group){}
var wg sync.WaitGroup
var mutex sync.Mutex

func LoadGroups() {
	api.LoadApiDataFromFile()
	println("Loading groups in cache for better performances")
	artists := api.GetGroupieArtistList()

	for _, v := range artists {
		go transformAndCacheGroup(v)
		wg.Add(1)
	}
	wg.Wait()
	go UpdateAllDeezerInformations(false)
	go UpdateAllGeolocInformation(false)
	println(strconv.Itoa(len(GroupMap)) + " groups have been loaded in cache!")
}

func transformAndCacheGroup(v api.ApiArtist) {
	defer wg.Done()
	group := objects.Group{}
	group.InitFromApiArtist(v)
	go addArtist(group.Members)
	go AddGroupDates(group)
	mutex.Lock()
	GroupMap[v.Id] = group
	mutex.Unlock()
}

func UpdateAllDeezerInformations(forceUpdate bool) {
	println("[ASYNC] Loading informations from deezer's api")
	for k, v := range GroupMap {
		v.DZInformations = api.GetDeezerInformationsFromName(v.Name, forceUpdate)
		v.DefineMostValuableGenreForGroup()
		mutex.Lock()
		GroupMap[k] = v
		s := GroupByGenreMap[v.MostValuableGenre]
		if !GroupSliceContain(s, v) {
			s = append(s, v)
		}
		GroupByGenreMap[v.MostValuableGenre] = s
		mutex.Unlock()
		go addTrackByGroup(v.DZInformations.TrackList, v)
	}
	api.SaveApiCacheToFile()
	println("Updated all deezer information !")
}

func UpdateAllGeolocInformation(forceUpdate bool) {
	println("[ASYNC] Loading geoloc informations")
	for i, value := range GroupMap {
		for location, dates := range value.DateLocations {
			for _, date := range dates {
				city := location
				geoloc := api.GetInformationForCity(city)
				date.Loc = geoloc
				val, ok := value.DateLocations[location]
				l := []api.Date{date}
				if !ok {
					l = append(l, val...)
				}
				mutex.Lock()
				value.DateLocations[location] = l
				GroupMap[i] = value
				mutex.Unlock()
			}
		}
	}
	api.SaveApiCacheToFile()
	println("All geoloc information are loaded")
}

func GroupSliceContain(s []objects.Group, v objects.Group) bool {
	for _, value := range s {
		if value.Name == v.Name {
			return true
		}
	}
	return false
}

func UpdateAlternativeGroupsForGroup(group *objects.Group) {
	group.GroupAlternatives = GroupByGenreMap[group.MostValuableGenre]
}

func GetDeezerGenreList() []api.DeezerGenre {
	return maps.Keys(GroupByGenreMap)
}

func FilterGroupsByName(filter string) []objects.Group {
	tlist := []objects.Group{}
	for _, v := range maps.Clone(GroupByGenreMap) {
		for _, g := range v {
			if strings.Contains(strings.ToUpper(g.Name), strings.ToUpper(filter)) {
				tlist = append(tlist, g)
			}
		}
	}
	slices.SortFunc(tlist, func(a, b objects.Group) bool {
		return strings.Index(
			strings.ToUpper(a.Name), strings.ToUpper(filter)) < strings.Index(
			strings.ToUpper(b.Name), strings.ToUpper(filter))
	})
	return tlist
}
