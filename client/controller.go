package client

import (
	"fmt"
	"github.com/echokepler/megad2561/adapter"
	"github.com/echokepler/megad2561/configs"
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/ports"
)

type OptionsController struct {
	Host     string
	Password string
	core.ServiceAdapter
}

type Controller struct {
	host       string
	password   string
	service    core.ServiceAdapter
	Ports      ports.Ports
	MqttClient core.MqttService
	configs.MainConfig
	configs.MegadIDConfig
}

// NewController создает инстанс контроллера.
// В дальнейшем будет служить для инициализации http сервиса и mqtt соединения
func NewController(opts OptionsController) (*Controller, error) {
	var service core.ServiceAdapter

	if opts.ServiceAdapter == nil {
		service = &adapter.HTTPAdapter{
			Host: fmt.Sprintf("%v/%v", opts.Host, opts.Password),
		}
	} else {
		service = opts.ServiceAdapter
	}

	portHub := ports.NewPorts(service) // 0_-

	controller := Controller{
		host:     opts.Host,
		password: opts.Password,
		service:  service,
		Ports:    *portHub,
	}

	configList := configs.NewConfigs([]configs.ConfigReader{
		&controller.MainConfig,
		&controller.MegadIDConfig,
	}, service)

	err := configList.Read()
	if err != nil {
		return nil, err
	}
	err = controller.Ports.Read()
	if err != nil {
		return nil, err
	}

	if controller.MainConfig.IsMQTTEnabled() {
		MQTTCreds := controller.MainConfig.GetMQTTCredentials()
		megadID := controller.MegadIDConfig.GetMegadID()

		MQTTOpts := adapter.MQTTClientOptions{
			Address:  MQTTCreds.Host,
			ClientID: &megadID,
			Password: MQTTCreds.Password,
		}

		MQTTClient, err := adapter.NewMqttClient(MQTTOpts)
		if err != nil {
			return nil, err
		}

		controller.MqttClient = MQTTClient

		for id, port := range controller.Ports.Records {
			if port.GetType() == ports.InputType {
				err := controller.MqttClient.SubscribePortIn(uint8(id), func(msg core.MegadPortInMessage) {
					fmt.Println(msg)
				})

				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &controller, nil
}
