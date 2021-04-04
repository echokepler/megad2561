package output

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/ports/base"
	"strconv"
)

type PortPWM struct {
	*base.OutputPort

	Value uint8
}

func (port *PortPWM) Set(value uint8) {
	port.Value = value
}

func (port *PortPWM) Read(values core.ServiceValues) error {
	err := port.OutputPort.Read(values)
	if err != nil {
		return err
	}

	if values.Has("pwm") {
		pwm, err := strconv.ParseInt(values.Get("pwm"), 10, 64)
		if err != nil {
			return err
		}

		port.Value = uint8(pwm)
	}

	return nil
}
