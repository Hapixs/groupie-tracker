package workers

import (
	"api"
	"logger"
	"objects"
	"strconv"
	"strings"
	"time"
	"utils"

	"golang.org/x/exp/maps"
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

		DefineMostValuableGenreForGroup(group)
		mutex.Lock()
		GroupList = append(GroupList, group)
		GroupById[group.Id] = group
		ok := false
		for k := range GroupByGenre {
			if k.Name == group.MostValuableGenre.Name {
				GroupByGenre[k] = append(GroupByGenre[k], group)
				ok = !ok
				break
			}
		}
		if !ok {
			GroupByGenre[group.MostValuableGenre] = append(make([]*objects.Group, 0), group)
		}
		mutex.Unlock()
	}
	logger.Log("All groups are loaded")
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
	for k, v := range GroupByGenre {
		if k.Name == group.MostValuableGenre.Name {
			group.GroupAlternatives = v
			break
		}
	}
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

func DefineMostValuableGenreForGroup(group *objects.Group) {
	top := new(objects.MusicGenre)
	top.Name = "none"
	top.Id = -1

	table := map[int](int){0: 0}
	for _, track := range group.TrackList {
		i := 1
		val, ok := table[track.Genre.Id]
		if ok {
			i += val
		}
		table[track.Genre.Id] = i
	}
	for k, v := range table {
		if v > table[top.Id] {
			top = GenreById[k]
		}
	}
	id := strconv.Itoa(top.Id)
	if id == "0" {
		top.Id = 0
	}
	group.MostValuableGenre = top
}

func AdvencedFilter(year int, members []int, location string) []*objects.Group {
	tlist := []*objects.Group{}
	for _, v := range GroupList {
		ok := true
		if year > 0 {
			ok = v.CreationYear <= year
		}
		if len(members) > 0 {
			m := false
			for _, i := range members {
				println(i)
				m = len(v.Members) == i || m
			}
			if !m {
				continue
			}
		}
		if location != "" {
			ok = ok && slices.Contains(maps.Keys(v.LocationMap), location)
		}
		if ok {
			tlist = append(tlist, v)
		}
	}
	return tlist

}
