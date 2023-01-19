package objects

import (
	"api"

	"golang.org/x/exp/maps"
)

type Artist struct {
	Name      string
	GroupName string
	ImageLink string
	WikiUrl   string
}

func (artist *Artist) updatePicture() {
	request := api.GetWikipediaImage(artist.Name)
	for _, page := range maps.Values(request.Query.Page) {
		switch page.Thumbnail.Source {
		case "":
			if artist.GroupName == artist.Name {
				artist.ImageLink = "https://cdn-icons-png.flaticon.com/512/32/32297.png"
			} else {
				tempArtist := Artist{Name: artist.GroupName, GroupName: artist.GroupName}
				tempArtist.updatePicture()
				artist.ImageLink = tempArtist.ImageLink
			}
		default:
			artist.ImageLink = page.Thumbnail.Source
		}
	}
}

func (artist *Artist) Load() {
	artist.WikiUrl = api.GetWikipediaPageLink(artist.Name)
	artist.updatePicture()
}
