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
		{http.MethodGet, "ip", "", "ip=127.0.0.1", ""},
		{http.MethodGet, "ip", "", "ip=255.255.255.255", ""},
		{http.MethodGet, "ip", "", "ip=999.999.999.999", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", "", "ip=::1", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", "", "ip=127.0.0.0001", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", "", "ip=abc.def.ghi.jkl", "invalid input: formato de ip inválido"},
		{http.MethodGet, "ip", "", "", "invalid input: campo 'ip' é obrigatório"},
		{http.MethodGet, "ip", "{country=127.0.0.1}", "ip=8.8.8.8", "invalid input: solicitações GET não devem ter dados enviados via body"},
		{http.MethodGet, "country", "", "", "invalid input: campo 'country' é obrigatório"},
		{http.MethodGet, "country", "{ip=127.0.0.1}", "country=br", "invalid input: solicitações GET não devem ter dados enviados via body"},
		{http.MethodGet, "country", "", "country=127.0.0.1", "invalid input: nome ou código de país invalido"},
		{http.MethodGet, "country", "", "country=São Paulo", "invalid input: nome ou código de país invalido"},
		{http.MethodPost, "country", `{"country":"Brazil"}`, "", ""},
		{http.MethodPost, "country", "", "country=Brazil", "invalid input: solicitações POST não devem ter dados enviados via url"},
		{http.MethodDelete, "ip", "", "ip=127.0.0.1", ""},
		{http.MethodDelete, "ip", "{ip= }", "ip=8.8.8.8", "invalid input: solicitações DELETE não devem ter dados enviados via body"},
		{http.MethodPut, "ip", "", "", "invalid input: método não suportado"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.input, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(tt.method, "/?"+tt.query, nil)
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.body != "" {
				c.Request.Body = io.NopCloser(strings.NewReader(tt.body))
			}
			_, err := CheckEntryData(tt.input, c)
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("Expected error %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}
