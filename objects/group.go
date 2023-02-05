package objects

import (
	"api"
	"os"
	"utils"
)

type Group struct {
	// Direct From groupie
	Id             int    `json:"id"`
	ImageLink      string `json:"image_link"`
	Name           string `json:"name"`
	CreationYear   int    `json:"creation_year"`
	FirstAlbumDate string `json:"first_album_date"`

	//Direct from deezer
	DeezerId  int      `json:"deezerId"`
	TrackList []*Track `json:"trackList"`

	//Transformed by workers
	DateLocations map[string]([]api.Date)  `json:"date_locations"` // deaprected
	LocationMap   map[string]([]*Location) `json:"locations"`
	Members       []Artist                 `json:"members"`

	DZInformations    api.DeezerInformations `json:"deezer_informations"` // to remove
	GroupAlternatives []*Group               `json:"group_alternatives"`
	MostValuableGenre *MusicGenre            `json:"most_valuable_genre"`
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
	top := new(MusicGenre)
	table := map[*MusicGenre](int){top: 0}
	for _, track := range group.TrackList {
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
