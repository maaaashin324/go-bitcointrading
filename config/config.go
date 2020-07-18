package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

// List can be used to get any environment variable
type List struct {
	APIKey        string
	APISecret     string
	LOGFile       string
	ProductCode   string
	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DbName        string
	SQLDriver     string
	Port          int
}

// Config is a struct of List. You can use this struct to get any environment variable
var Config List

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to load config.ini: %v", err)
		os.Exit(1)
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = List{
		APIKey:        cfg.Section("bitflyer").Key("api_key").String(),
		APISecret:     cfg.Section("bitflyer").Key("api_secret").String(),
		LOGFile:       cfg.Section("gotrading").Key("log_file").MustString("gotrading.log"),
		ProductCode:   cfg.Section("gotrading").Key("product_code").MustString("BTC_JPY"),
		Durations:     durations,
		TradeDuration: durations[cfg.Section("gotrading").Key("trade_duration").String()],
		DbName:        cfg.Section("db").Key("name").String(),
		SQLDriver:     cfg.Section("db").Key("driver").String(),
		Port:          cfg.Section("web").Key("port").MustInt(8080),
	}
}
