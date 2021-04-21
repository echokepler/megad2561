package base

type PortType int

const (
	InputType PortType = iota
	OutputType
)

type Port struct {
	id int
	t  PortType
}

func (p *Port) GetID() int {
	return p.id
}

func (p *Port) GetType() PortType {
	return p.t
}
