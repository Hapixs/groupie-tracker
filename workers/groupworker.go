package workers

import (
	"api"
	"logger"
	"objects"
	"strings"
	"time"
	"utils"

	"golang.org/x/exp/slices"
)

func prepareGroups() {
	type groupieArtist struct {
		Id             int                   `json:"id"`
		Image          string                `json:"image"`
		Name           string                `json:"name"`
		Members        []string              `json:"members"`
		CreationDate   int                   `json:"creationDate"`
		FirstAlbum     string                `json:"firstAlbum"`
		DatesLocations map[string]([]string) `json:"datesLocations"`
	}
	var groupieArtistList []groupieArtist
	groupieArtistListUrl := "https://groupietrackers.herokuapp.com/api/artists"

	type deezerQuerySearch struct {
		Data []struct {
			SearchArtist struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"data"`
	}

	api.GetFromApi(
		groupieArtistListUrl,
		&groupieArtistList,
		false,
		time.Millisecond,
		nil)

	for _, groupieArtist := range groupieArtistList {
		group := new(objects.Group)

		// Simple transformation
		group.Id = groupieArtist.Id
		group.ImageLink = groupieArtist.Image
		group.Name = groupieArtist.Name
		group.CreationYear = groupieArtist.CreationDate
		group.FirstAlbumDate = groupieArtist.FirstAlbum

		println("working for " + group.Name)

		// Complex transformation
		//todo members
		//todo dates
		url := "https://api.deezer.com/search?q=" + utils.FormatArtistName(group.Name)
		deezerQuerySearchRequest := deezerQuerySearch{}

		api.GetFromApi(
			url,
			&deezerQuerySearchRequest,
			false,
			time.Second/9,
			nil)

		if len(deezerQuerySearchRequest.Data) > 0 {
			deezerGroupId := deezerQuerySearchRequest.Data[0].SearchArtist.Id
			for _, data := range deezerQuerySearchRequest.Data {
				if strings.Contains(data.SearchArtist.Name, group.Name) {
					deezerGroupId = data.SearchArtist.Id
					break
				}
			}
			group.TrackList = buildTrackListFromGroupId(deezerGroupId, group.Id, group.Name)
		} else {
			logger.Log("No deezer informations found for " + group.Name)
		}

		group.LocationMap = buildLocationFromGroupId(group.Id)

		for _, member := range groupieArtist.Members {
			group.Members = append(group.Members, *buildArtist(member, group.Name, group.Id))
		}

		mutex.Lock()
		GroupList = append(GroupList, group)
		GroupById[group.Id] = group
		mutex.Unlock()
	}
}

// func LoadGroups() {
// 	api.LoadApiDataFromFile()
// 	println("Loading groups in cache for better performances")
// 	artists := api.GetGroupieArtistList()

// 	for _, v := range artists {
// 		go transformAndCacheGroup(v)
// 		waitgroup.Add(1)
// 	}
// 	waitgroup.Wait()
// 	go UpdateAllDeezerInformations(false)
// 	go UpdateAllGeolocInformation(false)
// 	println(strconv.Itoa(len(GroupMap)) + " groups have been loaded in cache!")
// }

// func transformAndCacheGroup(v api.ApiArtist) {
// 	defer waitgroup.Done()
// 	group := objects.Group{}
// 	group.InitFromApiArtist(v)
// 	mutex.Lock()
// 	GroupMap[v.Id] = group
// 	mutex.Unlock()
// }

// func UpdateAllDeezerInformations(forceUpdate bool) {
// 	println("[ASYNC] Loading informations from deezer's api")
// 	for k, v := range GroupMap {
// 		v.DZInformations = api.GetDeezerInformationsFromName(v.Name, forceUpdate)
// 		v.DefineMostValuableGenreForGroup()
// 		mutex.Lock()
// 		GroupMap[k] = v
// 		s := GroupByGenreMap[v.MostValuableGenre]
// 		if !GroupSliceContain(s, v) {
// 			s = append(s, v)
// 		}
// 		GroupByGenreMap[v.MostValuableGenre] = s
// 		mutex.Unlock()
// 	}
// 	api.SaveApiCacheToFile()
// 	println("Updated all deezer information !")
// }

// func UpdateAllGeolocInformation(forceUpdate bool) {
// 	println("[ASYNC] Loading geoloc informations")
// 	for i, value := range GroupMap {
// 		for location, dates := range value.DateLocations {
// 			for _, date := range dates {
// 				city := location
// 				geoloc := api.GetInformationForCity(city)
// 				date.Loc = geoloc
// 				val, ok := value.DateLocations[location]
// 				l := []api.Date{date}
// 				if !ok {
// 					l = append(l, val...)
// 				}
// 				mutex.Lock()
// 				value.DateLocations[location] = l
// 				GroupMap[i] = value
// 				mutex.Unlock()
// 			}
// 		}
// 	}
// 	api.SaveApiCacheToFile()
// 	println("All geoloc information are loaded")
// }

func GroupSliceContain(s []objects.Group, v objects.Group) bool {
	for _, value := range s {
		if value.Name == v.Name {
			return true
		}
	}
	return false
}

func UpdateAlternativeGroupsForGroup(group *objects.Group) {
	group.GroupAlternatives = GroupByGenre[group.MostValuableGenre]
}

func GetDeezerGenreList() []*objects.MusicGenre {
	return slices.Clone(GenreList)
}

func FilterGroupsByName(filter string) []*objects.Group {
	tlist := []*objects.Group{}
	for _, v := range GroupList {
		if strings.Contains(strings.ToUpper(v.Name), strings.ToUpper(filter)) {
			tlist = append(tlist, v)
		}
	}
	slices.SortFunc(tlist, func(a, b *objects.Group) bool {
		return strings.Index(
			strings.ToUpper(a.Name), strings.ToUpper(filter)) < strings.Index(
			strings.ToUpper(b.Name), strings.ToUpper(filter))
	})
	return tlist
}
