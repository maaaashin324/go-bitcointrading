package main

import (
	"gotrading/config"
	"gotrading/utils"
	"log"
)

func main() {
	utils.LoggingSettings(config.Config.LOGFile)
	log.Println("test")
}
