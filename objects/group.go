package objects

import (
	"api"
	"os"
	"utils"
)

type Group struct {
	Id             int                     `json:"id"`
	ImageLink      string                  `json:"image_link"`
	Name           string                  `json:"name"`
	Members        []Artist                `json:"members"`
	CreationYear   int                     `json:"creation_year"`
	FirstAlbumDate string                  `json:"first_album_date"`
	DateLocations  map[string]([]api.Date) `json:"date_locations"`

	DZInformations    api.DeezerInformations `json:"deezer_informations"`
	GroupAlternatives []Group                `json:"group_alternatives"`
	MostValuableGenre api.DeezerGenre        `json:"most_valuable_genre"`
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
			GroupId:   group.Id,
		}
		artist.Load()
		group.Members = append(group.Members, artist)
	}

	api.UpdateGroupRelation(&apiartist)
	group.DateLocations = map[string]([]api.Date){}
	for k, v := range apiartist.DatesLocations {
		val, ok := group.DateLocations[k]
		l := []api.Date{{
			Locations: k,
			DateTime:  v,
		}}
		if !ok {
			l = append(l, val...)
		}
		group.DateLocations[k] = l
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
	switch group.MostValuableGenre.Name {
	case "Classique":
		group.MostValuableGenre.FontName = "Nicoone"
	case "Dance":
		group.MostValuableGenre.FontName = "Orbitron"
	case "Rap/Hip Hop":
		group.MostValuableGenre.FontName = "Lacquer"
	case "Pop":
		group.MostValuableGenre.FontName = "Reem_Kufi_Ink"
	case "Reggae":
		group.MostValuableGenre.FontName = "Shadows_Into_Light"
	case "Metal":
		group.MostValuableGenre.FontName = "Metal_Mania"
	case "Rock":
		group.MostValuableGenre.FontName = "Rock_Salt"
	case "Alternative":
		group.MostValuableGenre.FontName = "Unbounded"
	}
}
