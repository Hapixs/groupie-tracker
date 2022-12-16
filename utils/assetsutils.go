package utils

import (
	"api"
	"os"
)

func PrepareFolders() {

	assetsErr := os.Mkdir("static/assets", os.ModePerm)

	if assetsErr != nil {
		println(assetsErr.Error())
	}

	albumErr := os.Mkdir("static/assets/album", os.ModePerm)

	if albumErr != nil {
		println(assetsErr.Error())
	}

	artistErr := os.Mkdir("static/assets/artistes", os.ModePerm)

	if artistErr != nil {
		println(assetsErr.Error())
	}
}

func UpdateAllAlbumPics() {
	Artists := api.GetAllArtist()
	for _, a := range Artists {
		fileHash := CalculatStringHash(a.Name)
		DownloadPicture(a.Image, "static/assets/album/"+fileHash+".png")
	}
}

func UpdateAllArtistsPics() {
	for _, a := range api.GetAllArtist() {
		for _, m := range a.Members {
			go UpdateArtistPic(m)
		}
	}
}

func UpdateArtistPic(name string) {
	fileHash := CalculatStringHash(name)
	DownloadPicture(api.GetArtistPictureLink(name), "static/assets/artistes/"+fileHash+"x")

}
