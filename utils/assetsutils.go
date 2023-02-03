package utils

import (
	"os"
)

func PrepareFolders() error {
	err := os.Mkdir("static/assets", os.ModePerm)

	if err != nil {
		return err
	}

	err = os.Mkdir("static/assets/groups", os.ModePerm)
	return err
}
