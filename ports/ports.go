package ports

import (
	"errors"
	"fmt"
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/ports/base"
	"strconv"
)

type Ports struct {
	Service core.ServiceAdapter
	Records map[int]PortReader
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
		values, err := p.Service.Get(params)
		if err != nil {
			return err
		}

		port, err := NewPort(id, values)
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

func (p *Ports) GetByID(id int, t base.PortType) (PortReader, error) {
	port := p.Records[id]

	if t != port.GetType() {
		return nil, errors.New("port not found")
	}

	return port, nil
}

func (p *Ports) Set(reader PortReader) error {
	p.Records[reader.GetID()] = reader

	port := p.Records[reader.GetID()]
	values, err := port.Write()
	if err != nil {
		return err
	}

	values.Add("pn", strconv.FormatInt(int64(port.GetID()), 10))
	values.Add("pty", strconv.FormatInt(int64(port.GetType()), 10))

	return p.Service.Post(values)
}
