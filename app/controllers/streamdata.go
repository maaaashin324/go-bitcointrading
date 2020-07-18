package controllers

import (
	"log"

	"gotrading/app/models"
	"gotrading/bitflyer"
	"gotrading/config"
)

func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.APIKey, config.Config.APISecret)
	go apiClient.GetRealtimeTicker(config.Config.ProductCode, tickerChannel)
	go func() {
		for ticker := range tickerChannel {
			log.Printf("action=StreamIngestionData, %v", ticker)
			for _, duration := range config.Config.Durations {
				isCreated := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreated && duration == config.Config.TradeDuration {
					// TODO
				}
			}
		}
	}()
}
