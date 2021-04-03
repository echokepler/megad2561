package megad2561

type ConfigReader interface {
	Read(service ServiceAdapter) error
	Write(service ServiceAdapter) error
}

type Configs []ConfigReader

func (cs Configs) Read(service ServiceAdapter) error {
	for _, config := range cs {
		err := config.Read(service)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs Configs) Write(service ServiceAdapter) error {
	for _, config := range cs {
		err := config.Write(service)
		if err != nil {
			return err
		}
	}

	return nil
}
