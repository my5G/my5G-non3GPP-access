package http_wrapper_test

import (
	"github.com/stretchr/testify/assert"
	"free5gc/lib/http_wrapper"
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080?name=NCTU&location=Hsinchu", nil)
	req.Header.Set("Location", "https://www.nctu.edu.tw/")
	request := http_wrapper.NewRequest(req, 1000)
	assert.Equal(t, "https://www.nctu.edu.tw/", request.Header.Get("Location"))
	assert.Equal(t, "NCTU", request.Query.Get("name"))
	assert.Equal(t, "Hsinchu", request.Query.Get("location"))
	assert.Equal(t, 1000, request.Body)
}
