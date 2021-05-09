package core

type PortType int

const (
	InputType PortType = iota
	OutputType
)

type PortReader interface {
	GetID() uint
	GetPortType() PortType
}
