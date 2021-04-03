package megad2561

import "fmt"

type OptionsController struct {
	Host     string
	Password string
	ServiceAdapter
}

type Controller struct {
	host     string
	password string
	service  ServiceAdapter
	Ports    Ports
	MainConfig
	MegadIDConfig

	configs Configs
}

/**
* NewController создает инстанс контроллера.
*
* В дальнейшем будет служить для инициализации http сервиса и mqtt соединения
**/
func NewController(opts OptionsController) (*Controller, error) {
	var service ServiceAdapter

	if opts.ServiceAdapter == nil {
		service = &HTTPAdapter{
			Host: fmt.Sprintf("%v/%v", opts.Host, opts.Password),
		}
	}

	controller := Controller{
		host:     opts.Host,
		password: opts.Password,
		service:  service,
		Ports: Ports{
			service: service,
			Records: map[int]PortReader{},
		},
	}

	configs := Configs{
		&controller.MainConfig,
		&controller.MegadIDConfig,
	}

	controller.configs = configs

	err := configs.Read(controller.service)
	if err != nil {
		return nil, err
	}

	err = controller.Ports.Read()
	if err != nil {
		return nil, err
	}

	return &controller, nil
}

/**
* ApplyConfigsChanges применяет изменения, отправляея измененные конфиги в сервис.
*
**/
func (c *Controller) ApplyConfigsChanges() error {
	return c.configs.Write(c.service)
}
