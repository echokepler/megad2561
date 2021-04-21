package core

type PortType int

const (
	InputType PortType = iota
	OutputType
)

type PortReader interface {
	ConfigReader
	GetID() uint
	GetPortType() PortType
}
