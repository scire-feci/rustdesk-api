package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// LocalOnly admits only requests made directly over loopback, i.e. by a process on this
// host such as Home Assistant with host networking. Requests relayed by a proxy (Tailscale
// Serve also connects from loopback, but adds X-Forwarded-For) or coming from the network
// are refused, independently of the trusted-proxy setting.
func LocalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isDirectLoopback(c.Request) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func isDirectLoopback(r *http.Request) bool {
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP", "Forwarded"} {
		if r.Header.Get(h) != "" {
			return false
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
