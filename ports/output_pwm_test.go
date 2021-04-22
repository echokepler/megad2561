package ports

import (
	"github.com/echokepler/megad2561/adapter"
	"github.com/echokepler/megad2561/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPortPWM_Set(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "0", params.Get("pt"))
		assert.Equal(t, "255", params.Get("pwm"))
	}))

	service := adapter.HTTPAdapter{Host: server.URL}
	port := NewPortPWM(0, &service)

	err := port.Set(255)
	if err != nil {
		t.Error(err)
	}
}

func TestPortPWM_write(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		actual   PortPWM
		expected core.ServiceValues
	}{
		{
			name: "Should be return correct values",
			actual: PortPWM{
				Value: 211,
				OutputPort: &OutputPort{
					Port: &Port{
						id: 0,
					},
					Group: "2",
					Mode:  SW,
				},
			},
			expected: core.ServiceValues{
				"m":   []string{"0"},
				"grp": []string{""},
				"d":   []string{"false"},
				"pwm": []string{"211"},
			},
		},
		{
			name: "Should be return correct values with mode PWM",
			actual: PortPWM{
				OutputPort: &OutputPort{
					Port: &Port{
						id: 0,
					},
					Mode: PWM,
				},
			},
			expected: core.ServiceValues{
				"m":   []string{"1"},
				"grp": []string{""},
				"d":   []string{"false"},
				"pwm": []string{"0"},
			},
		},
		{
			name: "Should be return correct values with mode SWL",
			actual: PortPWM{
				OutputPort: &OutputPort{
					Port: &Port{
						id: 0,
					},
					Mode: SWL,
				},
			},
			expected: core.ServiceValues{
				"m":   []string{"2"},
				"grp": []string{""},
				"d":   []string{"false"},
				"pwm": []string{"0"},
			},
		},
		{
			name: "Should be return correct values with mode DS2413",
			actual: PortPWM{
				OutputPort: &OutputPort{
					Port: &Port{
						id: 0,
					},
					Mode: DS2413,
				},
			},
			expected: core.ServiceValues{
				"m":   []string{"3"},
				"grp": []string{""},
				"d":   []string{"false"},
				"pwm": []string{"0"},
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
