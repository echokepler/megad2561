package megad2561

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPortInput_Read(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   ServiceValues
		expected PortInput
	}{
		{
			name: "Should be PortMode Click after read",
			values: ServiceValues{
				"m":  []string{"3"},
				"af": []string{"true"},
			},
			expected: PortInput{
				Mode:           CLICK,
				ForceSendToNet: true,
			},
		},
		{
			name: "Should be Port is muted after read",
			values: ServiceValues{
				"mt": []string{"true"},
			},
			expected: PortInput{
				IsMute: true,
			},
		},
		{
			name: "Should be Port raw mode disabled after read",
			values: ServiceValues{
				"d": []string{"true"},
			},
			expected: PortInput{
				IsRaw: true,
			},
		},
		{
			name: "Must be a property commands after reading",
			values: ServiceValues{
				"ecmd": []string{"21:2;g0:0"},
			},
			expected: PortInput{
				Commands: "21:2;g0:0",
			},
		},
		{
			name: "Must be a property net commands after reading",
			values: ServiceValues{
				"eth": []string{"0.0.0.0/megad.php"},
			},
			expected: PortInput{
				NetCommandAddress: "0.0.0.0/megad.php",
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actualPort := PortInput{}

			err := actualPort.Read(tCase.values)
			if err != nil {
				t.Error(err)
			}

			assert.EqualValues(t, tCase.expected, actualPort)
		})
	}
}

func TestPortInput_Write(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		actual   PortInput
		expected ServiceValues
	}{
		{
			name: "Should be return correct values",
			actual: PortInput{
				BasePort: &BasePort{
					ID: 0,
				},
				Commands: "22:2",
				IsRaw:    true,
				IsMute:   true,
				Mode:     PR,
			},
			expected: ServiceValues{
				"pn":   []string{"0"},
				"pty":  []string{"0"},
				"ecmd": []string{"22:2"},
				"eth":  []string{""},
				"m":    []string{"2"},
				"af":   []string{"false"},
				"naf":  []string{"false"},
				"d":    []string{"true"},
				"mt":   []string{"true"},
			},
		},
		{
			name: "Should be return correct values with mode P",
			actual: PortInput{
				BasePort: &BasePort{
					ID: 0,
				},
				Mode: P,
			},
			expected: ServiceValues{
				"pn":   []string{"0"},
				"pty":  []string{"0"},
				"ecmd": []string{""},
				"eth":  []string{""},
				"m":    []string{"0"},
				"af":   []string{"false"},
				"naf":  []string{"false"},
				"d":    []string{"false"},
				"mt":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode R",
			actual: PortInput{
				BasePort: &BasePort{
					ID: 0,
				},
				Mode: R,
			},
			expected: ServiceValues{
				"pn":   []string{"0"},
				"pty":  []string{"0"},
				"ecmd": []string{""},
				"eth":  []string{""},
				"m":    []string{"1"},
				"af":   []string{"false"},
				"naf":  []string{"false"},
				"d":    []string{"false"},
				"mt":   []string{"false"},
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			values, err := tCase.actual.Write()
			if err != nil {
				t.Error(err)

				return
			}

			assert.EqualValues(t, tCase.expected, values)
		})
	}
}
