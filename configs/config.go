package configs

import "github.com/echokepler/megad2561/core"

type ConfigReader interface {
	Read(service core.ServiceAdapter) error
	Write(service core.ServiceAdapter) error
}

type Configs []ConfigReader

func (cs Configs) Read(service core.ServiceAdapter) error {
	for _, config := range cs {
		err := config.Read(service)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs Configs) Write(service core.ServiceAdapter) error {
	for _, config := range cs {
		err := config.Write(service)
		if err != nil {
			return err
		}
	}

	return nil
}
