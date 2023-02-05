package main

import (
	"handlers"
	"logger"
	"net/http"
	"objects"
	"os"
	"utils"
	"workers"
)

func main() {
	objects.InitConfig()
	objects.ProcessArguments(os.Args[1:])
	logger.PrepareLogger(objects.Static_FlagConfig_Verbose)
	utils.PrepareFolders()
	workers.Init()
	handlers.InitHandlers()
	serverport := objects.GetConfigValue[string](objects.ServerPort)
	logger.Log("[WEB] Start server on port " + serverport)
	http.ListenAndServe(":"+serverport, nil)
}
