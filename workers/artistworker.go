package workers

import (
	"objects"
	"strings"

	"golang.org/x/exp/slices"
)

func addArtist(artists []objects.Artist) {
	mutex.Lock()
	artistsList = append(artistsList, artists...)
	mutex.Unlock()
}

func FiltreAllArtistByName(filter string) []objects.Artist {
	list := []objects.Artist{}
	for _, a := range artistsList {
		if strings.Contains(strings.ToUpper(a.Name), strings.ToUpper(filter)) {
			list = append(list, a)
		}
	}

	slices.SortFunc(list, func(a, b objects.Artist) bool {
		return strings.Index(strings.ToUpper(a.Name), strings.ToUpper(filter)) < strings.Index(strings.ToUpper(b.Name), strings.ToUpper(filter))
	})

	return list
}
