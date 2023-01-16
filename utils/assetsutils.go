package utils

import (
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
