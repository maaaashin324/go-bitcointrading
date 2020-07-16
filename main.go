package main

import (
	"fmt"
	"gotrading/bitflyer"
	"gotrading/config"
	"gotrading/utils"
	"time"
)

func main() {
	utils.LoggingSettings(config.Config.LOGFile)
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)

	tickerChannel := make(chan bitflyer.Ticker)
	go apiClient.GetRealtimeTicker(config.Config.ProductCode, tickerChannel)
	for ticker := range tickerChannel {
		fmt.Println(ticker)
		fmt.Println(ticker.GetMidPrice())
		fmt.Println(ticker.DateTime())
		fmt.Println(ticker.TruncateDateTime(time.Second))
		fmt.Println(ticker.TruncateDateTime(time.Minute))
	}
}
