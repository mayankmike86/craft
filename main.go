package main

import (
	"crafts/config"
	"fmt"
)

func main() {
	fmt.Println("hello mayank1")
	config := config.GetConfig()
	fmt.Println("config :", config)
}
