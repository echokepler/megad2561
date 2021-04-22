package main

import (
	"github.com/echokepler/megad2561/client"
	"github.com/echokepler/megad2561/ports"
)

const IdInputPort = 4

func main() {
	options := client.OptionsController{
		Host:     "192.168.88.14",
		Password: "sec",
	}
	controller, err := client.NewController(options)
	if err != nil {
		panic(err)
	}

	port, err := controller.Ports.GetByID(IdInputPort, ports.InputType)
	if err != nil {
		panic(err)
	}

	portInput := port.(*ports.InputPort)
	portInput.Commands = "22:2|g0:0;g1:0;22:0"

	err = controller.Ports.Set(port)
	if err != nil {
		panic(err)
	}
}
