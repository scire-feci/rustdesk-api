package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

// StripRemoteAddrZone drops the zone from link-local IPv6 peers ("[fe80::1%eth0]:port").
// gin's ClientIP runs net.ParseIP on the host, which rejects zoned addresses and returns "",
// so such clients were recorded without an IP and all shared one limiter/ban bucket.
func StripRemoteAddrZone() gin.HandlerFunc {
	return func(c *gin.Context) {
		if host, port, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
			if i := strings.IndexByte(host, '%'); i >= 0 {
				c.Request.RemoteAddr = net.JoinHostPort(host[:i], port)
			}
		}
		c.Next()
	}
}
