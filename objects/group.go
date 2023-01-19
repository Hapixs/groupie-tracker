package objects

import (
	"api"
	"os"
	"utils"
)

type Group struct {
	Id             int
	ImageLink      string
	Name           string
	Members        []Artist
	CreationYear   int
	FirstAlbumDate string
	DateLocations  []api.Date

	DZInformations    api.DeezerInformations
	GroupAlternatives []Group
	MostValuableGenre api.DeezerGenre
}

func (group *Group) UpdatePicture() {
	fileHash := utils.CalculatStringHash(group.Name)
	_, err := os.OpenFile("static/assets/groups/"+fileHash+".jpeg", os.O_RDONLY, os.ModePerm)
	if err != nil {
		utils.DownloadPicture(group.ImageLink, "static/assets/groups/"+fileHash+".jpeg")
	}
	group.ImageLink = "/static/assets/groups/" + fileHash + ".jpeg"
}

func (group *Group) InitFromApiArtist(apiartist api.ApiArtist) {
	group.Id = apiartist.Id
	group.ImageLink = apiartist.Image
	group.Name = apiartist.Name
	group.CreationYear = apiartist.CreationDate
	group.FirstAlbumDate = apiartist.FirstAlbum

	group.UpdatePicture()

	for _, member := range apiartist.Members {
		artist := Artist{
			Name:      member,
			GroupName: group.Name,
		}
		artist.Load()
		group.Members = append(group.Members, artist)
	}
}

func (group *Group) DefineMostValuableGenreForGroup() {
	var top api.DeezerGenre = api.DeezerGenre{Name: ""}
	table := map[api.DeezerGenre](int){top: 0}
	for _, track := range group.DZInformations.TrackList.List {
		i := 1
		val, ok := table[track.Album.Genre]
		if ok {
			i += val
		}
		table[track.Album.Genre] = i
	}
	for k, v := range table {
		if v > table[top] {
			top = k
		}
	}
	group.MostValuableGenre = top
}