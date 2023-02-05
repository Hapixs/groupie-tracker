package objects

import (
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
	LocationMap map[string]([]*Location) `json:"locations"`
	Members     []Artist                 `json:"members"`

	GroupAlternatives []*Group    `json:"group_alternatives"`
	MostValuableGenre *MusicGenre `json:"most_valuable_genre"`
}

func (group *Group) UpdatePicture() {
	fileHash := utils.CalculatStringHash(group.Name)
	_, err := os.OpenFile("static/assets/groups/"+fileHash+".jpeg", os.O_RDONLY, os.ModePerm)
	if err != nil {
		utils.DownloadPicture(group.ImageLink, "static/assets/groups/"+fileHash+".jpeg")
	}
	group.ImageLink = "/static/assets/groups/" + fileHash + ".jpeg"
}
