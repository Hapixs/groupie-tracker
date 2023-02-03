package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
)

func DownloadPicture(url string, localPath string) error {
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer file.Close()

	io.Copy(file, resp.Body)
	println("Downloaded " + file.Name())

	return nil
}

func CalculatStringHash(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}
