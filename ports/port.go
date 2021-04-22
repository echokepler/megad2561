package ports

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
	"strconv"
)

type PortType int

const (
	InputType PortType = iota
	OutputType
	NC = 255
)

type PortReader interface {
	read(values core.ServiceValues) error
	write() (core.ServiceValues, error)
	GetType() PortType
	GetID() int
}

type Port struct {
	id      int
	t       PortType
	service core.ServiceAdapter
}

func (p *Port) GetID() int {
	return p.id
}

func (p *Port) GetType() PortType {
	return p.t
}

func (p *Port) change(port PortReader) error {
	values := qsparser.Marshal(port, qsparser.MarshalOptions{OnlySetters: true})

	values.Add("pt", strconv.FormatInt(int64(p.id), 10))

	_, err := p.service.Get(values)

	return err
}

func NewPort(id int, values core.ServiceValues, service core.ServiceAdapter) (PortReader, error) {
	var port PortReader
	portTypeInt, err := strconv.ParseInt(values.Get("pty"), 10, 64)
	if err != nil {
		return nil, err
	}
	basePort := Port{
		id:      id,
		service: service,
	}

	switch PortType(portTypeInt) {
	case InputType:
		port = PortReader(&InputPort{
			Port: &basePort,
		})

	case OutputType:
		port, err = createOutputPort(id, service, values)
	case NC:
		return nil, nil

	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	err = port.read(values)
	if err != nil {
		return nil, err
	}

	return port, nil
}

func createOutputPort(id int, service core.ServiceAdapter, values core.ServiceValues) (PortReader, error) {
	mode, err := strconv.ParseInt(values.Get("m"), 10, 64)
	if err != nil {
		return nil, err
	}

	switch ModeTypeOUT(mode) {
	case PWM:
		return NewPortPWM(id, service), nil
	default:
		return NewOutputPort(id, service), nil
	}
}
