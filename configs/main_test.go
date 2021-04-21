package configs

import (
	"github.com/echokepler/megad2561/adapter"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainConfig(t *testing.T) {
	t.Run("MainConfig implement RemoteConfig interface", func(t *testing.T) {
		assert.Implements(t, (*ConfigReader)(nil), new(MainConfig))
	})
}

func TestMainConfig_Sync(t *testing.T) {
	t.Run("Should correct parse and write variables to config", func(t *testing.T) {
		t.Parallel()
		var server *httptest.Server
		var config MainConfig

		file, err := ioutil.ReadFile("./__mocks__/mainconfig.html")
		if err != nil {
			panic("Cant read mainconfig.html")
		}

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, string(file))
			if err != nil {
				panic("Cant respond html mainconfig")
			}
		}))
		service := adapter.HTTPAdapter{Host: server.URL}
		config.service = &service

		err = config.read()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "192.168.88.14", config.attributes.IP)
		assert.Equal(t, "sec", config.attributes.Pwd)
		assert.Equal(t, "255.255.255.255", config.attributes.Gateway)
		assert.Equal(t, "255.255.255.255:1883", config.attributes.Srv)
		assert.Equal(t, HTTP, config.attributes.SrvType)
		assert.Equal(t, "test value", config.attributes.ScriptPath)
		assert.Equal(t, "tst", config.attributes.Wdog)
		assert.Equal(t, GSM, config.attributes.UART)
	})
}

func TestMainConfig_Apply(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		cb       func(config *MainConfig)
		expected string
	}{
		{
			name: "Should be return correct ip",
			cb: func(config *MainConfig) {
				config.attributes.IP = "192.168.88.14"
				config.attributes.Pwd = "sec"
			},
			expected: "/?cf=1&eip=192.168.88.14&gsm=0&gw=&pr=&pwd=sec&sip=&srvt=0",
		},
		{
			name: "Should be return params with enabled mqtt",
			cb: func(config *MainConfig) {
				err := config.SetMQTTServer("192.168.88.14", "")
				if err != nil {
					t.Error(err)
				}
			},
			expected: "/?auth=&cf=1&eip=&gsm=0&gw=&pr=&pwd=&sip=192.168.88.14&srvt=1",
		},
		{
			name: "Should be assign mqtt password when enabled",
			cb: func(config *MainConfig) {
				err := config.SetMQTTServer("192.168.88.14", "password")
				if err != nil {
					t.Error(err)
				}
			},
			expected: "/?auth=password&cf=1&eip=&gsm=0&gw=&pr=&pwd=&sip=192.168.88.14&srvt=1",
		},
		{
			name: "Should be switch to http server",
			cb: func(config *MainConfig) {
				err := config.SetMQTTServer("192.168.88.14", "")
				if err != nil {
					t.Error(err)
				}
				err = config.SetHTTPServer("192.168.88.1")
				if err != nil {
					t.Error(err)
				}
			},
			expected: "/?cf=1&eip=&gsm=0&gw=&pr=&pwd=&sip=192.168.88.1&srvt=0",
		},
		{
			name: "Should be disable mqtt and http server",
			cb: func(config *MainConfig) {
				err := config.SetMQTTServer("192.168.88.14", "")
				if err != nil {
					t.Error(err)

					return
				}

				err = config.SetHTTPServer("192.168.88.1")
				if err != nil {
					t.Error(err)

					return
				}
				err = config.DisableSrv()
				if err != nil {
					t.Error(err)

					return
				}
			},
			expected: "/?cf=1&eip=&gsm=0&gw=&pr=&pwd=&sip=255.255.255.255",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var server *httptest.Server
			var config MainConfig

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" {
					assert.Equal(t, tc.expected, r.URL.String())
				}
			}))

			service := adapter.HTTPAdapter{Host: server.URL}

			config.service = &service

			tc.cb(&config)

			err := config.write()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
