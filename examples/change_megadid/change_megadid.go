package main

import (
	"github.com/echokepler/megad2561/client"
	"github.com/echokepler/megad2561/configs"
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

	err = controller.MegadIDConfig.ChangeSettings(func(settings configs.MegaIDSettings) configs.MegaIDSettings {
		settings.MegadID = "megs"

		return settings
	})

	if err != nil {
		panic(err)
	}
}
