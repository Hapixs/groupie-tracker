package workers

import (
	"objects"
	"strings"

	"golang.org/x/exp/slices"
)

func buildArtist(name, groupName string, groupId int) *objects.Artist {
	artist := new(objects.Artist)
	artist.Name, artist.GroupName, artist.GroupId = name, groupName, groupId
	artist.Load()

	mutex.Lock()
	ArtistList = append(ArtistList, artist)
	mutex.Unlock()

	return artist
}

func FiltreAllArtistByName(filter string) []*objects.Artist {
	list := []*objects.Artist{}
	for _, a := range ArtistList {
		if strings.Contains(strings.ToUpper(a.Name), strings.ToUpper(filter)) {
			list = append(list, a)
		}
	}

	slices.SortFunc(list, func(a, b *objects.Artist) bool {
		return strings.Index(strings.ToUpper(a.Name), strings.ToUpper(filter)) < strings.Index(strings.ToUpper(b.Name), strings.ToUpper(filter))
	})

	return list
}
