package megad2561

import "strconv"

type ModeTypeOUT int

const (
	SW ModeTypeOUT = iota
	PWM
	SWL
	DS2413
)

type PortOutput struct {
	*BasePort

	IsEnabledByDefault bool

	Mode ModeTypeOUT

	Group string
}

func (port *PortOutput) Read(values ServiceValues) error {
	var err error

	if values.Has("m") {
		mode, err := strconv.ParseInt(values.Get("m"), 10, 64)
		if err != nil {
			return err
		}

		port.Mode = ModeTypeOUT(mode)
	}

	if values.Has("d") {
		port.IsEnabledByDefault, err = strconv.ParseBool(values.Get("d"))
		if err != nil {
			return err
		}
	}

	port.Group = values.Get("grp")

	return nil
}

func (port *PortOutput) Write() (ServiceValues, error) {
	values := ServiceValues{}

	values.Add("d", strconv.FormatBool(port.IsEnabledByDefault))
	values.Add("m", strconv.FormatInt(int64(port.Mode), 10))
	values.Add("grp", port.Group)

	return values, nil
}
