package configs

import (
	"github.com/echokepler/megad2561/core"
	"strconv"
)

const (
	MegadIDConfigPath = "2"
)

type MegadIDConfig struct {
	MegadID string
	SrvLoop bool
}

func (config *MegadIDConfig) Read(service core.ServiceAdapter) error {
	params := core.ServiceValues{}

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

func (config *MegadIDConfig) Write(service core.ServiceAdapter) error {
	values := core.ServiceValues{}

	values.Add("cf", MegadIDConfigPath)
	values.Add("mdid", config.MegadID)
	values.Add("sl", strconv.FormatBool(config.SrvLoop))

	return service.Post(values)
}
