package ports

import (
	"github.com/echokepler/megad2561/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseOutputPort_read(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   core.ServiceValues
		expected OutputPort
	}{
		{
			name: "Should be PortMode SW after read",
			values: core.ServiceValues{
				"m":  []string{"0"},
				"af": []string{"true"},
			},
			expected: OutputPort{
				Mode: SW,
			},
		},
		{
			name: "Should be Port is enabled by default after read",
			values: core.ServiceValues{
				"d": []string{"true"},
			},
			expected: OutputPort{
				IsEnabledByDefault: true,
			},
		},
		{
			name: "Should be Port Group is 1 after read",
			values: core.ServiceValues{
				"grp": []string{"1"},
			},
			expected: OutputPort{
				Group: "1",
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actualPort := OutputPort{}

			err := actualPort.read(tCase.values)
			if err != nil {
				t.Error(err)
			}

			assert.EqualValues(t, tCase.expected, actualPort)
		})
	}
}

func TestPortOutput_write(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		actual   OutputPort
		expected core.ServiceValues
	}{
		{
			name: "Should be return correct values",
			actual: OutputPort{
				Port: &Port{
					id: 0,
				},
				Group: "2",
				Mode:  SW,
			},
			expected: core.ServiceValues{
				"m":   []string{"0"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode PWM",
			actual: OutputPort{
				Port: &Port{
					id: 0,
				},
				Mode: PWM,
			},
			expected: core.ServiceValues{
				"m":   []string{"1"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode SWL",
			actual: OutputPort{
				Port: &Port{
					id: 0,
				},
				Mode: SWL,
			},
			expected: core.ServiceValues{
				"m":   []string{"2"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode DS2413",
			actual: OutputPort{
				Port: &Port{
					id: 0,
				},
				Mode: DS2413,
			},
			expected: core.ServiceValues{
				"m":   []string{"3"},
				"grp": []string{""},
				"d":   []string{"false"},
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
