package ports

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
)

type ModeTypeOUT int

const (
	SW ModeTypeOUT = iota
	PWM
	SWL
	DS2413
)

type OutputPort struct {
	*Port

	IsEnabledByDefault bool `qs:"d"`

	Mode ModeTypeOUT `qs:"m"`

	Group string `qs:"grp"`
}

func NewOutputPort(id int, service core.ServiceAdapter) *OutputPort {
	return &OutputPort{
		Port: &Port{
			id:      id,
			service: service,
			t:       OutputType,
		},
	}
}

func (port *OutputPort) ChangeSettings(cb func(p OutputPort) OutputPort) error {
	updatedPort := cb(*port)
	values := qsparser.Marshal(updatedPort, qsparser.MarshalOptions{})

	err := port.service.Post(values)
	if err != nil {
		return err
	}

	return port.read(values)
}

func (port *OutputPort) read(values core.ServiceValues) error {
	return qsparser.UnMarshal(values, port)
}

func (port *OutputPort) write() (core.ServiceValues, error) {
	values := qsparser.Marshal(*port, qsparser.MarshalOptions{})

	return values, nil
}
