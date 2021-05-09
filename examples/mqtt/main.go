package main

import (
	"github.com/echokepler/megad2561/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	options := client.OptionsController{
		Host:     "http://192.168.88.14",
		Password: "sec",
	}

	_, err := client.NewController(options)
	if err != nil {
		panic(err)
	}

	//err = controller.MainConfig.SetMQTTServer("192.168.88.242:1883", ""); if err != nil {
	//	panic(err)
	//}

	<-c
}
