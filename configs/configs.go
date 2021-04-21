package configs

type ConfigReader interface {
	read() error
	write() error
}
type Configs []ConfigReader

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
