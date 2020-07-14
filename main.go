package main

import (
	"fmt"
	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
	"log"
	"time"
)

func main() {
	utils.LoggingSettings(config.Config.LOGFile)
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)
	ticker, err := apiClient.GetTicker("BTC_JPY")
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(ticker)
	fmt.Println(ticker.GetMidPrice())
	fmt.Println(ticker.DateTime())
	fmt.Println(ticker.TruncateDateTime(time.Minute))
}
