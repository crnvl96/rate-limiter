package utilities

import (
	"net"
	"net/http"
	"strings"
)

func GetRealIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")

	if xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			trimmedIP := strings.TrimSpace(ip)
			if isValidIP(trimmedIP) {
				return trimmedIP
			}
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" && isValidIP(xRealIP) {
		return xRealIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
