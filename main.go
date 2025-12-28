package main

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/router"
	"fmt"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Printf("config initialize error %v", err)
	}
	router.Initialize()

}
