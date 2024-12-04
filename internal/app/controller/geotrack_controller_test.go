package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCheckEntryData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		method      string
		input       string
		body        string
		query       string
		expectedErr string
	}{
		{http.MethodGet, "ip", "", "ip=8.8.8.8", ""},
		{http.MethodGet, "ip", "", "::1", "invalid input: campo 'ip' é obrigatório"},
		{http.MethodGet, "ip", "", "ip=::1", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", "", "ip=999.999.999.999", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", `{"ip":"8.8.8.8"}`, "ip=8.8.8.8", "invalid input: solicitações " + http.MethodGet + " não devem ter dados enviados via body"},

		{http.MethodGet, "country", "", "country=United States", ""},
		{method: http.MethodGet, input: "country", body: "", query: "br", expectedErr: "invalid input: campo 'country' é obrigatório"},
		{method: http.MethodGet, input: "country", body: "", query: "country=123", expectedErr: "invalid input: nome ou código de país invalido"},
		{http.MethodGet, "country", `{"ip":"8.8.8.8"}`, "country=br", "invalid input: solicitações " + http.MethodGet + " não devem ter dados enviados via body"},

		{http.MethodPost, "ip", `{"ip":"8.8.8.8"}`, "", ""},
		{http.MethodPost, "ip", `{"id":"8.8.8.8"}`, "", "invalid input: campo 'ip' é obrigatório"},
		{http.MethodPost, "ip", `{"ip":"999.999.999.999"}`, "", "invalid input: formato de ip inválido"},
		{http.MethodPost, "ip", "", "ip=8.8.8.8", "invalid input: solicitações POST não devem ter dados enviados via url"},
		{http.MethodPost, "ip", `{"country":"8.8.8.8"}`, "", "invalid input: campo 'ip' é obrigatório"},

		{http.MethodDelete, "ip", "", "ip=8.8.8.8", ""},
		{http.MethodDelete, "ip", "", "::1", "invalid input: campo 'ip' é obrigatório"},
		{http.MethodDelete, "ip", "", "ip=::1", "invalid input: formato de ip inválido"},
		{http.MethodDelete, "ip", "", "ip=999.999.999.999", "invalid input: formato de ip inválido"},
		{http.MethodDelete, "ip", `{"ip":"8.8.8.8"}`, "ip=8.8.8.8", "invalid input: solicitações " + http.MethodDelete + " não devem ter dados enviados via body"},

		{http.MethodPut, "ip", "", "teste=teste", "invalid input"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.input, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(tt.method, "/?"+tt.query, nil)
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.body != "" {
				c.Request.Body = io.NopCloser(strings.NewReader(tt.body))
				c.Request.ContentLength = int64(len(tt.body))
			}
			_, err := CheckEntryData(tt.input, c)
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("Expected error %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}
