package main

import (
	"fmt"
	"gotrading/config"
)

func main() {
	fmt.Println(config.Config.APIKey)
}
