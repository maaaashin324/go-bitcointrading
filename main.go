package main

import (
	"fmt"
	"gotrading/app/models"
	"gotrading/config"
	"gotrading/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LOGFile)
	fmt.Println(models.DbConnection)
}
