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

func TestPortInput_read(t *testing.T) {
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
				Mode:           CLICK,
				ForceSendToNet: true,
			},
		},
		{
			name: "Should be Port is muted after read",
			values: core.ServiceValues{
				"mt": []string{"true"},
			},
			expected: InputPort{
				IsMute: true,
			},
		},
		{
			name: "Should be Port raw mode disabled after read",
			values: core.ServiceValues{
				"d": []string{"true"},
			},
			expected: InputPort{
				IsRaw: true,
			},
		},
		{
			name: "Must be a property commands after reading",
			values: core.ServiceValues{
				"ecmd": []string{"21:2;g0:0"},
			},
			expected: InputPort{
				Commands: "21:2;g0:0",
			},
		},
		{
			name: "Must be a property net commands after reading",
			values: core.ServiceValues{
				"eth": []string{"0.0.0.0/megad.php"},
			},
			expected: InputPort{
				NetCommandAddress: "0.0.0.0/megad.php",
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

func TestPortInput_write(t *testing.T) {
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
				Commands: "22:2",
				IsRaw:    true,
				IsMute:   true,
				Mode:     PR,
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
				Mode: P,
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
				Mode: R,
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

func TestInputPort_ChangeSettings(t *testing.T) {
	t.Parallel()

	t.Run("Should update settings in remote service", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, "1", params.Get("d"))
		}))
		service := adapter.HTTPAdapter{Host: server.URL}
		port := NewInputPort(0, &service)

		err := port.ChangeSettings(func(p InputPort) InputPort {
			p.IsRaw = true

			return p
		})

		assert.True(t, port.IsRaw)

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Must not update properties if a negative response is received from the service", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		service := adapter.HTTPAdapter{Host: server.URL}
		port := NewInputPort(0, &service)

		err := port.ChangeSettings(func(p InputPort) InputPort {
			p.IsRaw = true

			return p
		})

		assert.Error(t, err)
		assert.False(t, port.IsRaw)

	})
}
