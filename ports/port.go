package ports

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/ports/base"
	"strconv"
)

type PortReader interface {
	Read(values core.ServiceValues) error
	Write() (core.ServiceValues, error)
	GetType() base.PortType
	GetID() int
}

func NewPort(id int, values core.ServiceValues) (PortReader, error) {
	var port PortReader
	portTypeInt, err := strconv.ParseInt(values.Get("pty"), 10, 64)
	if err != nil {
		return nil, err
	}
	basePort := base.Port{
		ID: id,
	}

	switch base.PortType(portTypeInt) {
	case base.InputType:
		port = PortReader(&base.InputPort{
			Port: &basePort,
		})

	case base.OutputType:
		port = PortReader(&base.OutputPort{
			Port: &basePort,
		})
	default:
		return nil, nil
	}

	err = port.Read(values)
	if err != nil {
		return nil, err
	}

	return port, nil
}
