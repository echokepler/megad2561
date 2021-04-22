package configs

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
)

type (
	SrvType int
	UART    int
)

const (
	HTTP SrvType = iota
	MQTT
)

const (
	DISABLED UART = iota
	GSM
	RS485
)

const (
	MainConfigPath     = "1"
	MatchLengthSrvMask = 4
)

// MainSettings основные настройки контроллера
type MainSettings struct {
	IP           string  `qs:"eip"`
	Pwd          string  `qs:"pwd"`
	Gateway      string  `qs:"gw"`
	Srv          string  `qs:"sip"`
	SrvType      SrvType `qs:"srvt"`
	ScriptPath   string  `qs:"sct"`
	MqttPassword string  `qs:"auth"`
	Wdog         string  `qs:"pr"`
	UART         `qs:"gsm"`
}

type MainConfig struct {
	service    core.ServiceAdapter
	attributes MainSettings
}

// GetSettings возвращает текущее состояние конфига
func (config *MainConfig) GetSettings() MainSettings {
	return config.attributes
}

// Update в коллбеке принимаем текущие параметры конфига и возвращаем новые
//
// Возвращенное состояние параметров будет обновлено
func (config *MainConfig) Update(cb func(settings MainSettings) MainSettings) error {
	updatedSettings := cb(config.attributes)

	config.attributes = updatedSettings
	return config.write()
}

func (config *MainConfig) SetMQTTServer(ip string, password string) error {
	config.attributes.Srv = ip
	config.attributes.SrvType = MQTT
	config.attributes.MqttPassword = password

	return config.write()
}

// SetHTTPServer Выставляет HTTP сервер.
//
// Обрати внимание, что при этом контроллер отключится от MQTT сервера.
func (config *MainConfig) SetHTTPServer(ip string) error {
	config.attributes.Srv = ip
	config.attributes.SrvType = HTTP
	config.attributes.MqttPassword = ""

	return config.write()
}

// DisableSrv Устанавливает управление портами на уровень контроллера.
func (config *MainConfig) DisableSrv() error {
	config.attributes.Srv = "255.255.255.255"

	return config.write()
}

/**
* Private
**/

func (config *MainConfig) read() error {
	params := core.ServiceValues{}

	params.Add("cf", MainConfigPath)

	values, err := config.service.Get(params)
	if err != nil {
		return err
	}

	return qsparser.UnMarshal(values, &config.attributes)
}

// write Отправляет значения в контроллер
//
// Обрати внимание, что после каждого вызова write мы синхронизируем
// конфиги, т.к внутри самого контролера может измениться логика валидации
// и поле, которое мы хотели изменить может остаться прежним.
func (config *MainConfig) write() error {
	values := qsparser.Marshal(config.attributes, qsparser.MarshalOptions{})

	values.Add("cf", MainConfigPath)

	err := config.service.Post(values)
	if err != nil {
		return err
	}

	return config.read()
}

func (config *MainConfig) setService(adapter core.ServiceAdapter) {
	config.service = adapter
}
