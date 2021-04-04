package megad2561

import "strconv"

type PortType int

const (
	InputType PortType = iota
	OutputType
)

type PortReader interface {
	Read(values ServiceValues) error
	Write() (ServiceValues, error)
	GetType() PortType
	GetID() int
}

type BasePort struct {
	ID int
	// t is port type
	t PortType
}

func (p *BasePort) GetID() int {
	return p.ID
}

func (p *BasePort) GetType() PortType {
	return p.t
}

func NewPort(id int, values ServiceValues) (PortReader, error) {
	var port PortReader
	portTypeInt, err := strconv.ParseInt(values.Get("pty"), 10, 64)
	if err != nil {
		return nil, err
	}
	basePort := BasePort{
		ID: id,
		t:  PortType(portTypeInt),
	}

	switch basePort.t {
	case InputType:
		port = PortReader(&PortInput{
			BasePort: &basePort,
		})

	case OutputType:
		port = PortReader(&PortInput{
			BasePort: &basePort,
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
