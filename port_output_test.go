package megad2561

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPortOutput_Read(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   ServiceValues
		expected PortOutput
	}{
		{
			name: "Should be PortMode SW after read",
			values: ServiceValues{
				"m":  []string{"0"},
				"af": []string{"true"},
			},
			expected: PortOutput{
				Mode: SW,
			},
		},
		{
			name: "Should be Port is enabled by default after read",
			values: ServiceValues{
				"d": []string{"true"},
			},
			expected: PortOutput{
				IsEnabledByDefault: true,
			},
		},
		{
			name: "Should be Port Group is 1 after read",
			values: ServiceValues{
				"grp": []string{"1"},
			},
			expected: PortOutput{
				Group: "1",
			},
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			actualPort := PortOutput{}

			err := actualPort.Read(tCase.values)
			if err != nil {
				t.Error(err)
			}

			assert.EqualValues(t, tCase.expected, actualPort)
		})
	}
}

func TestPortOutput_Write(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		actual   PortOutput
		expected ServiceValues
	}{
		{
			name: "Should be return correct values",
			actual: PortOutput{
				BasePort: &BasePort{
					ID: 0,
				},
				Group: "2",
				Mode:  SW,
			},
			expected: ServiceValues{
				"m":   []string{"0"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode PWM",
			actual: PortOutput{
				BasePort: &BasePort{
					ID: 0,
				},
				Mode: PWM,
			},
			expected: ServiceValues{
				"m":   []string{"1"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode SWL",
			actual: PortOutput{
				BasePort: &BasePort{
					ID: 0,
				},
				Mode: SWL,
			},
			expected: ServiceValues{
				"m":   []string{"2"},
				"grp": []string{""},
				"d":   []string{"false"},
			},
		},
		{
			name: "Should be return correct values with mode DS2413",
			actual: PortOutput{
				BasePort: &BasePort{
					ID: 0,
				},
				Mode: DS2413,
			},
			expected: ServiceValues{
				"m":   []string{"3"},
				"grp": []string{""},
				"d":   []string{"false"},
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
