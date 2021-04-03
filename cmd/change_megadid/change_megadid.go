package main

import (
	"github.com/echokepler/megad2561"
)

func main() {
	options := megad2561.OptionsController{
		Host:     "192.168.88.14",
		Password: "sec",
	}
	controller, err := megad2561.NewController(options)
	if err != nil {
		panic(err)
	}

	controller.MegadID = "new"

	err = controller.ApplyConfigsChanges()
	if err != nil {
		panic(err)
	}
}
