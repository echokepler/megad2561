package ports

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
)

type PortPWM struct {
	*OutputPort

	Value uint8 `qs:"setter,pwm"`
}

func NewPortPWM(id int, service core.ServiceAdapter) *PortPWM {
	return &PortPWM{
		OutputPort: &OutputPort{
			Port: &Port{
				id:      id,
				service: service,
				t:       OutputType,
			},
		},
	}
}

func (port *PortPWM) Set(value uint8) error {
	port.Value = value

	return port.change(port)
}

func (port *PortPWM) ChangeSettings(cb func(p PortPWM) PortPWM) error {
	updatedPort := cb(*port)
	values := qsparser.Marshal(updatedPort, qsparser.MarshalOptions{})

	err := port.service.Post(values)
	if err != nil {
		return err
	}

	return port.read(values)
}

func (port *PortPWM) read(values core.ServiceValues) error {
	return qsparser.UnMarshal(values, port)
}

func (port *PortPWM) write() (core.ServiceValues, error) {
	values := qsparser.Marshal(*port, qsparser.MarshalOptions{})

	return values, nil
}
