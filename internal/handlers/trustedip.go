package handlers

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckTrusted(trustedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if trustedSubnet == "" {
			c.Next()
		}
		ip := c.GetHeader("X-Real-IP")
		_, subnet, _ := net.ParseCIDR(trustedSubnet)

		if !subnet.Contains(net.ParseIP(ip)) {

			c.Status(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
