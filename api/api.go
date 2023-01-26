package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Date struct {
	Locations string
	DateTime  []string
	Loc       Geoloc
}

var mutex sync.Mutex
var CacheApi = map[string]string{}

func CallExternalApi[T any](url string, structure *T, sleepTime time.Duration, headerVar map[string]string) error {
	req, _ := http.NewRequest("GET", url, nil)

	for k, v := range headerVar {
		req.Header.Add(k, v)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Error when calling " + url)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("Error when reading body of " + url)
	}
	err = json.Unmarshal(content, structure)
	if err != nil {
		return errors.New("Error when pasing body of " + url)
	}

	mutex.Lock()
	CacheApi[url] = string(content)
	mutex.Unlock()

	time.Sleep(time.Duration(sleepTime))
	return nil
}

func CallCacheApi[T any](url string, structure *T) error {
	mutex.Lock()
	val, ok := CacheApi[url]
	mutex.Unlock()
	if ok {
		return json.Unmarshal([]byte(val), structure)
	}
	return errors.New("Unable to find in api cache " + url)
}

func GetFromApi[T any](url string, structure *T, update bool, sleepTime time.Duration, headerVar map[string]string) error {
	err := CallCacheApi(url, structure)
	if err != nil || update {
		err = CallExternalApi(url, structure, sleepTime, headerVar)
	}
	return err
}

func SaveApiCacheToFile() {
	mutex.Lock()
	save, err := json.MarshalIndent(CacheApi, "", "\t")
	mutex.Unlock()
	if err != nil {
		println("Error: JSON error")
		return
	}

	file, err := os.Create("data.json")
	if err != nil {
		println("Error 2 with data")
		return
	}

	file.Write(save)
	file.Close()
}

func LoadApiDataFromFile() {
	content, err := os.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(content, &CacheApi)
	} else {
		println("Seams like the first start of this web server.")
		println("Some operations may take several minutes.")
	}
}
