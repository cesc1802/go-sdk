package r2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptMaxRedirects(t *testing.T) {
	assert := assert.New(t)

	var pingURL, pongURL string
	var pingCount, pongCount int
	ping := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		pingCount++
		http.Redirect(rw, r, pongURL, http.StatusTemporaryRedirect)
		return
	}))
	defer ping.Close()

	pong := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		pongCount++
		http.Redirect(rw, r, pingURL, http.StatusTemporaryRedirect)
		return
	}))
	defer pong.Close()

	pingURL = ping.URL
	pongURL = pong.URL

	_, err := New(pingURL, OptMaxRedirects(32)).Discard()
	assert.True(ErrIsTooManyRedirects(err))
	assert.Equal(32, pingCount+pongCount)
}
