package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// List can be used to get any environment variable
type List struct {
	APIKey      string
	APISecret   string
	LOGFile     string
	ProductCode string
}

// Config is a struct of List. You can use this struct to get any environment variable
var Config List

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to load config.ini: %v", err)
		os.Exit(1)
	}

	Config = List{
		APIKey:      cfg.Section("bitflyer").Key("api_key").String(),
		APISecret:   cfg.Section("bitflyer").Key("api_secret").String(),
		LOGFile:     cfg.Section("gotrading").Key("log_file").MustString("gotrading.log"),
		ProductCode: cfg.Section("gotrading").Key("product_code").MustString("BTC_JPY"),
	}
}
