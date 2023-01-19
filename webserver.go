package main

import (
	"handlers"
	"net/http"
	"objects"
	"os"
	"strconv"
	"utils"
	"workers"
)

func main() {
	objects.WebServerConfig.InitConfig()
	println("Reading start arguments")
	objects.GameProcessArguments(os.Args[1:])
	utils.PrepareFolders()
	workers.LoadGroups()
	handlers.InitHandlers()
	serverport, _, _ := objects.WebServerConfig.GetConfigItem(objects.ServerPort)
	http.ListenAndServe(":"+strconv.Itoa(serverport), nil)
}
