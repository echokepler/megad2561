package ports

import (
	"errors"
	"fmt"
	"github.com/echokepler/megad2561/core"
	"strconv"
)

type Ports struct {
	service core.ServiceAdapter
	Records map[int]PortReader
}

func NewPorts(service core.ServiceAdapter) *Ports {
	return &Ports{
		service: service,
		Records: map[int]PortReader{},
	}
}

// Read считываем информацию со всех портов.
//
// В данный момент беру все порты до XP5andXP6, не очень хорошая реализация, но поправимо в будущем.
func (p *Ports) Read() error {
	maxRangeInt, err := core.XP5andXP6.GetRange()
	if err != nil {
		return err
	}

	IDs := make([]int, maxRangeInt[1])

	for id := range IDs {
		params := core.ServiceValues{}

		params.Add("pt", strconv.FormatInt(int64(id), 10))
		values, err := p.service.Get(params)

		if values != nil && values.IsEmpty() {
			continue
		}

		if err != nil {
			return err
		}

		port, err := NewPort(id, values, p.service)
		if err != nil {
			fmt.Println(err)
			break
		}

		if port == nil {
			continue
		}

		p.Records[id] = port
	}

	return nil
}

func (p *Ports) GetByID(id int, t PortType) (PortReader, error) {
	port := p.Records[id]

	if t != port.GetType() {
		return nil, errors.New("port not found")
	}

	return port, nil
}

func (p *Ports) Set(reader PortReader) error {
	p.Records[reader.GetID()] = reader

	port := p.Records[reader.GetID()]
	values, err := port.write()
	if err != nil {
		return err
	}

	values.Add("pn", strconv.FormatInt(int64(port.GetID()), 10))
	values.Add("pty", strconv.FormatInt(int64(port.GetType()), 10))

	return p.service.Post(values)
}

func (p *Ports) ChangeService(adapter core.ServiceAdapter) {
	p.service = adapter
}
