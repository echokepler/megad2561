package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_isAuthorized(t *testing.T) {
	t.Parallel()
	t.Run("Should be return true when credentials is valid", func(t *testing.T) {
		var server *httptest.Server

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/sec/" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}))

		defer server.Close()

		opts := OptionsController{
			Host:     server.URL,
			Password: "sec",
		}

		_, err := NewController(opts)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should be return false when credentials is invalid", func(t *testing.T) {
		var server *httptest.Server

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		defer server.Close()

		opts := OptionsController{
			Host:     server.URL,
			Password: "s",
		}

		_, err := NewController(opts)
		if err != nil {
			assert.EqualError(t, err, "unauthorized")
		}

		assert.Error(t, err)
	})
}
