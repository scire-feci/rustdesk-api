package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func clientIPFor(remoteAddr string, handlers ...gin.HandlerFunc) string {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.Use(handlers...)
	var ip string
	g.GET("/", func(c *gin.Context) { ip = c.ClientIP() })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remoteAddr
	g.ServeHTTP(httptest.NewRecorder(), req)
	return ip
}

func TestZonedRemoteAddrHasNoClientIPWithoutMiddleware(t *testing.T) {
	if ip := clientIPFor("[fe80::1%eth0]:21114"); ip != "" {
		t.Fatalf("gin now parses zoned addresses (got %q); StripRemoteAddrZone may be unnecessary", ip)
	}
}

func TestStripRemoteAddrZone(t *testing.T) {
	cases := map[string]string{
		"[fe80::6e1f:f7ff:fec1:db9a%br0]:21114": "fe80::6e1f:f7ff:fec1:db9a",
		"[2001:db8::1]:21114":                   "2001:db8::1",
		"192.168.1.65:21114":                    "192.168.1.65",
	}
	for remoteAddr, want := range cases {
		if got := clientIPFor(remoteAddr, StripRemoteAddrZone()); got != want {
			t.Errorf("ClientIP for %s = %q, want %q", remoteAddr, got, want)
		}
	}
}
