package main

import (
	"gotrading/app/controllers"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LOGFile)
	controllers.StreamIngestionData()
	controllers.StartWebServer()
}
