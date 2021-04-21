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

func TestMegadIDConfig(t *testing.T) {
	t.Run("should MegadIDConfig implement RemoteConfig interface", func(t *testing.T) {
		assert.Implements(t, (*ConfigReader)(nil), new(MegadIDConfig))
	})
}

func TestMegadIDConfig_Apply(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		cb       func(config *MegadIDConfig)
		expected string
	}{
		{
			name: "Should be return correct megad Id",
			cb: func(config *MegadIDConfig) {
				config.attributes.MegadID = "test"
			},
			expected: "/?cf=2&mdid=test",
		},
		{
			name: "Should be return enabled srv loop",
			cb: func(config *MegadIDConfig) {
				config.attributes.SrvLoop = true
			},
			expected: "/?cf=2&mdid=&sl=1",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var server *httptest.Server
			var config MegadIDConfig

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" {
					assert.Equal(t, tc.expected, r.URL.String())
				}
			}))

			service := adapter.HTTPAdapter{
				Host: server.URL,
			}

			tc.cb(&config)

			config.service = &service

			err := config.write()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMegadIDConfig_Sync(t *testing.T) {
	t.Run("Should correct parse and write variables to config", func(t *testing.T) {
		t.Parallel()
		var server *httptest.Server
		var config MegadIDConfig

		file, err := ioutil.ReadFile("./__mocks__/megadidconfig.html")
		if err != nil {
			t.Error(err)
		}

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, string(file))
			if err != nil {
				t.Error(err)
			}
		}))

		service := adapter.HTTPAdapter{
			Host: server.URL,
		}

		config.service = &service

		err = config.read()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "megad", config.attributes.MegadID)
		assert.Equal(t, true, config.attributes.SrvLoop)
	})
}
