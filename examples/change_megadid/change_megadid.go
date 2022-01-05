package main

import (
	"encoding/json"
	"fmt"
	"github.com/echokepler/megad2561/client"
	"github.com/echokepler/megad2561/configs"
	"github.com/echokepler/megad2561/core"
)

func main() {
	options := client.OptionsController{
		Host:     "http://192.168.88.14",
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

	data := map[string]core.ServiceValues{}

	data["cf1"] = core.ServiceValues{"cf1": []string{"1"}, "setting": []string{"2"}}

	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
}
