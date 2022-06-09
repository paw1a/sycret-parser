package handler

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_DocEndpoint(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{
									"URLTemplate": "https://sycret.ru/service/apigendoc/forma_025u.doc", 
									"RecordID": "20"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Bad request (no url)",
			inputBody:          `{"RecordID": "20"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Bad request (no record)",
			inputBody:          `{"URLTemplate": "https://sycret.ru/service/apigendoc/forma_025u.doc"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Bad request (empty body)",
			inputBody:          ``,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid document url",
			inputBody: `{
									"URLTemplate": "invalid url", 
									"RecordID": "20"}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	assert.Nil(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(
				"GET",
				"/api/doc",
				bytes.NewBufferString(test.inputBody),
			)
			req.Header.Set("Content-type", "application/json")
			w := httptest.NewRecorder()

			DocEndpoint(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
