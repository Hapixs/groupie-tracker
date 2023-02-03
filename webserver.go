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
	objects.GameProcessArguments(os.Args[1:])
	utils.PrepareFolders()
	workers.LoadGroups()
	handlers.InitHandlers()
	serverport, _, _ := objects.WebServerConfig.GetConfigItem(objects.ServerPort)
	http.ListenAndServe(":"+strconv.Itoa(serverport), nil)
	println("[WEB] Server ready to use on port" + strconv.Itoa(serverport))
}
