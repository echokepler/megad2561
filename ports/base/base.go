package base

type PortType int

const (
	InputType PortType = iota
	OutputType
)

type Port struct {
	ID int
	t  PortType
}

func (p *Port) GetID() int {
	return p.ID
}

func (p *Port) GetType() PortType {
	return p.t
}
