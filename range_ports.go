package megad2561

import (
	"strconv"
	"strings"
)

type RangePortsModule string

const (
	XP1       RangePortsModule = "1,14"
	XP2       RangePortsModule = "15,29"
	XT2       RangePortsModule = "30,35"
	XP5andXP6 RangePortsModule = "36,37"
)

func (rpm RangePortsModule) GetRange() ([]int, error) {
	rangePorts := []int{0, 0}

	if len(rpm) > 0 {
		for i, value := range strings.Split(string(rpm), ",") {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}

			rangePorts[i] = int(intValue)
		}
	}

	if rangePorts[0] > rangePorts[1] {
		rangePorts[1] = rangePorts[0]
	}

	return rangePorts, nil
}
