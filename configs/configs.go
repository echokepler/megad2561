package configs

import "github.com/echokepler/megad2561/core"

type ConfigReader interface {
	setService(adapter core.ServiceAdapter)
	read() error
	write() error
}

type Configs []ConfigReader

func NewConfigs(configs []ConfigReader, adapter core.ServiceAdapter) *Configs {
	var hub Configs

	for _, config := range configs {
		config.setService(adapter)

		hub = append(hub, config)
	}

	return &hub
}

func (cs Configs) Read() error {
	for _, config := range cs {
		err := config.read()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs Configs) Write() error {
	for _, config := range cs {
		err := config.write()
		if err != nil {
			return err
		}
	}

	return nil
}
