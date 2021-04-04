package main

import (
	"github.com/echokepler/megad2561/client"
)

func main() {
	options := client.OptionsController{
		Host:     "192.168.88.14",
		Password: "sec",
	}
	controller, err := client.NewController(options)
	if err != nil {
		panic(err)
	}

	controller.MegadID = "new"

	err = controller.ApplyConfigsChanges()
	if err != nil {
		panic(err)
	}
}
