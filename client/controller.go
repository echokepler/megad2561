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
	host     string
	password string
	service  core.ServiceAdapter
	Ports    ports.Ports
	configs.MainConfig
	configs.MegadIDConfig

	configs configs.Configs
}

// NewController создает инстанс контроллера.
// В дальнейшем будет служить для инициализации http сервиса и mqtt соединения
func NewController(opts OptionsController) (*Controller, error) {
	var service core.ServiceAdapter

	if opts.ServiceAdapter == nil {
		service = &adapter.HTTPAdapter{
			Host: fmt.Sprintf("http://%v/%v", opts.Host, opts.Password),
		}
	}

	controller := Controller{
		host:     opts.Host,
		password: opts.Password,
		service:  service,
		Ports: ports.Ports{
			Service: service,
			Records: map[int]ports.PortReader{},
		},
	}

	configList := configs.Configs{
		&controller.MainConfig,
		&controller.MegadIDConfig,
	}

	controller.configs = configList

	err := configList.Read(controller.service)
	if err != nil {
		return nil, err
	}

	err = controller.Ports.Read()
	if err != nil {
		return nil, err
	}

	return &controller, nil
}

// ApplyConfigsChanges применяет изменения, отправляея измененные конфиги в сервис.
func (c *Controller) ApplyConfigsChanges() error {
	return c.configs.Write(c.service)
}
