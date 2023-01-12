package utils

import (
	"api"
	"objects"
	"os"
)

func PrepareFolders() {

	assetsErr := os.Mkdir("static/assets", os.ModePerm)

	if assetsErr != nil {
		println(assetsErr.Error())
	}

	groupsErr := os.Mkdir("static/assets/groups", os.ModePerm)

	if groupsErr != nil {
		println(assetsErr.Error())
	}
}

func UpdateAllGroupsPics() {
	_, downloadImages, _ := objects.WebServerConfig.GetConfigItem(objects.DownloadPicture)
	if !downloadImages {
		return
	}
	println("Checking and updating groups images..")
	groups := api.GetCachedGroups()
	for _, a := range groups {
		fileHash := CalculatStringHash(a.Name)
		_, err := os.OpenFile("static/assets/groups/"+fileHash+".jpeg", os.O_RDONLY, os.ModePerm)
		if err != nil {
			DownloadPicture(a.ImageLink, "static/assets/groups/"+fileHash+".jpeg")
		}
		api.EditGroupImageLink(a.Id, "/static/assets/groups/"+fileHash+".jpeg")
	}
	println("All groups images are downloaded !")
}
