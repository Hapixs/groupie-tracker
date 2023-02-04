package main

import (
	"handlers"
	"net/http"
	"objects"
	"os"
	"utils"
	"workers"
)

func main() {
	objects.InitConfig()
	objects.ProcessArguments(os.Args[1:])
	utils.PrepareFolders()
	workers.LoadGroups()
	handlers.InitHandlers()
	serverport := objects.GetConfigValue[string](objects.ServerPort)
	println("[WEB] Start server on port " + serverport)
	http.ListenAndServe(":"+serverport, nil)
}
