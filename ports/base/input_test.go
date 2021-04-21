package base

import (
	"github.com/echokepler/megad2561/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPortInput_Read(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   core.ServiceValues
		expected InputPort
	}{
		{
			name: "Should be PortMode Click after read",
			values: core.ServiceValues{
				"m":  []string{"3"},
				"af": []string{"true"},
			},
			expected: InputPort{
				settings: InputSettings{
					Mode:           CLICK,
					ForceSendToNet: true,
				},
			},
		},
		{
			name: "Should be Port is muted after read",
			values: core.ServiceValues{
				"mt": []string{"true"},
			},
			expected: InputPort{
				settings: InputSettings{
					IsMute: true,
				},
			},
		},
		{
			name: "Should be Port raw mode disabled after read",
			values: core.ServiceValues{
				"d": []string{"true"},
			},
			expected: InputPort{
				settings: InputSettings{
					IsRaw: true,
				},
			},
		},
		{
			name: "Must be a property commands after reading",
			values: core.ServiceValues{
				"ecmd": []string{"21:2;g0:0"},
			},
			expected: InputPort{
				settings: InputSettings{
					Commands: "21:2;g0:0",
				},
			},
		},
		{
			name: "Must be a property net commands after reading",
			values: core.ServiceValues{
				"eth": []string{"0.0.0.0/megad.php"},
			},
			expected: InputPort{
				settings: InputSettings{
					NetCommandAddress: "0.0.0.0/megad.php",
				},
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actualPort := InputPort{}

			err := actualPort.read(tCase.values)
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
		actual   InputPort
		expected core.ServiceValues
	}{
		{
			name: "Should be return correct values",
			actual: InputPort{
				Port: &Port{
					id: 0,
				},
				settings: InputSettings{
					Commands: "22:2",
					IsRaw:    true,
					IsMute:   true,
					Mode:     PR,
				},
			},
			expected: core.ServiceValues{
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
			actual: InputPort{
				Port: &Port{
					id: 0,
				},
				settings: InputSettings{
					Mode: P,
				},
			},
			expected: core.ServiceValues{
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
			actual: InputPort{
				Port: &Port{
					id: 0,
				},
				settings: InputSettings{
					Mode: R,
				},
			},
			expected: core.ServiceValues{
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

			values, err := tCase.actual.write()
			if err != nil {
				t.Error(err)

				return
			}

			assert.EqualValues(t, tCase.expected, values)
		})
	}
}
