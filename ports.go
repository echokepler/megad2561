package megad2561

import (
	"errors"
	"fmt"
	"strconv"
)

type Ports struct {
	service ServiceAdapter
	Records map[int]PortReader
}

// Read считываем информацию со всех портов.
//
// В данный момент беру все порты до XP5andXP6, не очень хорошая реализация, но поправимо в будущем.
func (ports *Ports) Read() error {
	maxRangeInt, err := XP5andXP6.GetRange()
	if err != nil {
		return err
	}

	IDs := make([]int, maxRangeInt[1])

	for id := range IDs {
		params := ServiceValues{}

		params.Add("pt", strconv.FormatInt(int64(id), 10))
		values, err := ports.service.Get(params)
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

		ports.Records[id] = port
	}

	return nil
}

func (ports *Ports) GetByID(id int, t PortType) (PortReader, error) {
	port := ports.Records[id]

	if t != port.GetType() {
		return nil, errors.New("port not found")
	}

	return port, nil
}

func (ports *Ports) Set(reader PortReader) error {
	ports.Records[reader.GetID()] = reader

	port := ports.Records[reader.GetID()]
	values, err := port.Write()
	if err != nil {
		return err
	}

	values.Add("pn", strconv.FormatInt(int64(port.GetID()), 10))
	values.Add("pty", strconv.FormatInt(int64(port.GetType()), 10))

	return ports.service.Post(values)
}
