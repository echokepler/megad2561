package megad2561

import (
	"strconv"
)

const (
	MegadIDConfigPath = "2"
)

type MegadIDConfig struct {
	MegadID string
	SrvLoop bool
}

func (config *MegadIDConfig) Read(service ServiceAdapter) error {
	params := ServiceValues{}

	params.Add("cf", MegadIDConfigPath)

	values, err := service.Get(params)
	if err != nil {
		return err
	}

	srvLoop, err := strconv.ParseBool(values.Get("sl"))
	if err != nil {
		return err
	}

	config.MegadID = values.Get("mdid")
	config.SrvLoop = srvLoop

	return nil
}

func (config *MegadIDConfig) Write(service ServiceAdapter) error {
	values := ServiceValues{}

	values.Add("cf", MegadIDConfigPath)
	values.Add("mdid", config.MegadID)
	values.Add("sl", strconv.FormatBool(config.SrvLoop))

	return service.Post(values)
}
