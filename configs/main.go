package configs

import (
	"github.com/echokepler/megad2561/core"
	"strconv"
	"strings"
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

type MainConfig struct {
	IP           string
	Pwd          string
	Gateway      string
	Srv          string
	SrvType      SrvType
	ScriptPath   string
	MqttPassword string
	Wdog         string
	UART
}

func (config *MainConfig) Read(service core.ServiceAdapter) error {
	params := core.ServiceValues{}

	params.Add("cf", MainConfigPath)

	values, err := service.Get(params)
	if err != nil {
		return err
	}

	if len(values.Get("srvt")) > 0 {
		srvTypeInt, err := strconv.ParseInt(values.Get("srvt"), 10, 64)
		if err != nil {
			return err
		}

		config.SrvType = SrvType(srvTypeInt)
	}

	if len(values.Get("gsm")) > 0 {
		uartInt, err := strconv.ParseInt(values.Get("gsm"), 10, 64)
		if err != nil {
			return err
		}

		config.UART = UART(uartInt)
	}

	config.IP = values.Get("eip")
	config.Pwd = values.Get("pwd")
	config.Gateway = values.Get("gw")
	config.Srv = values.Get("sip")
	config.Wdog = values.Get("pr")
	config.ScriptPath = values.Get("sct")

	if config.SrvType == MQTT {
		config.MqttPassword = values.Get("auth")
	}

	return nil
}

// Write Отправляет значения в контроллер
func (config *MainConfig) Write(service core.ServiceAdapter) error {
	values := core.ServiceValues{}

	values.Add("cf", MainConfigPath)
	values.Add("eip", config.IP)
	values.Add("pwd", config.Pwd)
	values.Add("gw", config.Gateway)
	values.Add("sip", config.Srv)
	values.Add("pr", config.Wdog)
	values.Add("gsm", strconv.FormatInt(int64(config.UART), 10))

	if strings.Count(config.Srv, "255") != MatchLengthSrvMask {
		values.Add("srvt", strconv.FormatInt(int64(config.SrvType), 10))

		if config.SrvType == MQTT {
			values.Add("auth", config.MqttPassword)
		}
	}

	return service.Post(values)
}

func (config *MainConfig) SetMQTTServer(ip string, password string) {
	config.Srv = ip
	config.SrvType = MQTT
	config.MqttPassword = password
}

// SetHTTPServer Выставляет HTTP сервер.
//
// Обрати внимание, что при этом контроллер отключится от MQTT сервера.
func (config *MainConfig) SetHTTPServer(ip string) {
	config.Srv = ip
	config.SrvType = HTTP
	config.MqttPassword = ""
}

// DisableSrv Устанавливает управление портами на уровень контроллера.
func (config *MainConfig) DisableSrv() {
	config.Srv = "255.255.255.255"
}
