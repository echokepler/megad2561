package main

import (
	"fmt"
	"github.com/echokepler/megad2561/client"
	"github.com/echokepler/megad2561/ports"
)

const IdPortPWM = 13

func main() {
	options := client.OptionsController{
		Host:     "192.168.88.14",
		Password: "sec",
	}
	controller, err := client.NewController(options)
	if err != nil {
		panic(err)
	}

	fmt.Println(controller.Ports)
	port, err := controller.Ports.GetByID(IdPortPWM, ports.OutputType)
	if err != nil {
		panic(err)
	}

	portPWM := port.(*ports.PortPWM)
	_ = portPWM.Set(0)
}
