package configs

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
)

const (
	MegadIDConfigPath = "2"
)

// MegaIDSettings основные настройки контроллера
type MegaIDSettings struct {
	MegadID string `qs:"mdid"`
	SrvLoop bool   `qs:"sl"`
}

type MegadIDConfig struct {
	service    core.ServiceAdapter
	attributes MegaIDSettings
}

func (config *MegadIDConfig) ChangeSettings(cb func(config MegaIDSettings) MegaIDSettings) error {
	config.attributes = cb(config.attributes)

	return config.write()
}

func (config *MegadIDConfig) read() error {
	params := core.ServiceValues{}

	params.Add("cf", MegadIDConfigPath)

	values, err := config.service.Get(params)
	if err != nil {
		return err
	}

	return qsparser.UnMarshal(values, &config.attributes)
}

func (config *MegadIDConfig) write() error {
	values := qsparser.Marshal(config.attributes, qsparser.MarshalOptions{})

	values.Add("cf", MegadIDConfigPath)

	err := config.service.Post(values)
	if err != nil {
		return err
	}

	return config.read()
}

func (config *MegadIDConfig) setService(adapter core.ServiceAdapter) {
	config.service = adapter
}
