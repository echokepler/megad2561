package megad2561

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHTTPAdapter_Get(t *testing.T) {
	t.Run("Should correct serialize programconfig", func(t *testing.T) {
		t.Parallel()
		var server *httptest.Server

		file, err := ioutil.ReadFile("./mock/programconfig.html")
		if err != nil {
			t.Error("Cant read programconfig.html")
		}

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, string(file))
			if err != nil {
				t.Error("Cant respond html programconfig")
			}
		}))

		httpAdapter := HTTPAdapter{
			Host: server.URL,
		}

		values, _ := httpAdapter.Get(ServiceValues{})

		assert.Equal(t, "prp", values.Get("prp"), "Input PRP")
		assert.Equal(t, "0", values.Get("prc"), "Select PRC")
		assert.Equal(t, "false", values.Get("prs"), "Checkbox PRS")
		assert.Equal(t, "true", values.Get("prs2"), "Checkbox PRS2")
		assert.Equal(t, "10", values.Get("cf"), "Hidden input CF")
	})
}

func TestHTTPAdapter_formatValues(t *testing.T) {
	t.Run("Should convert strconv bool to 0 or 1", func(t *testing.T) {
		httpAdapter := HTTPAdapter{}

		originalValues := ServiceValues{}

		originalValues.Add("bool", strconv.FormatBool(true))
		formattedValues, err := httpAdapter.formatValues(originalValues)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "1", formattedValues.Get("bool"))
	})
}
