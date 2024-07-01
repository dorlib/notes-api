package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNote(t *testing.T) {
	requestBody := []byte(`{"title": "Test Note", "content": "This is a test note."}`)
	request, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(requestBody))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
