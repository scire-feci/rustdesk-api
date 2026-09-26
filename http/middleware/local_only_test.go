package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.GET("/", LocalOnly(), func(c *gin.Context) { c.Status(http.StatusOK) })
	cases := []struct {
		remoteAddr string
		header     string
		want       int
	}{
		{"127.0.0.1:50000", "", http.StatusOK},
		{"[::1]:50000", "", http.StatusOK},
		{"127.0.0.1:50000", "X-Forwarded-For", http.StatusForbidden}, // via Tailscale Serve
		{"127.0.0.1:50000", "X-Real-IP", http.StatusForbidden},
		{"192.168.1.65:50000", "", http.StatusForbidden},
		{"192.168.1.65:50000", "X-Forwarded-For", http.StatusForbidden},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = tc.remoteAddr
		if tc.header != "" {
			req.Header.Set(tc.header, "127.0.0.1")
		}
		w := httptest.NewRecorder()
		g.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Errorf("%s with %q: got %d, want %d", tc.remoteAddr, tc.header, w.Code, tc.want)
		}
	}
}
