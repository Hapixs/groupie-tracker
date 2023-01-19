package workers

import (
	"api"
	"objects"
	"strconv"
	"sync"

	"golang.org/x/exp/maps"
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
	println("Asyncron, loading all deezer informations for groups")
	go UpdateAllDeezerInformations(false)
	println(strconv.Itoa(len(GroupMap)) + " groups have been loaded in cache!")
}

func transformAndCacheGroup(v api.ApiArtist) {
	defer wg.Done()
	group := objects.Group{}
	group.InitFromApiArtist(v)

	mutex.Lock()
	GroupMap[v.Id] = group
	mutex.Unlock()
}

func UpdateAllDeezerInformations(forceUpdate bool) {
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
	}
	api.SaveApiCacheToFile()
	println("Updated all deezer information !")
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