package netutil

//for getting the actual ip either form x-forwarded or x-realip headers falls back to remoteip
import (
	"net"
	"net/http"
	"strings"
)

// IsPrivate checks whether an IP is private
func IsPrivate(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	privateBlocks := []*net.IPNet{
		// IPv4 private ranges
		mustCIDR("10.0.0.0/8"),
		mustCIDR("172.16.0.0/12"),
		mustCIDR("192.168.0.0/16"),
		mustCIDR("127.0.0.0/8"),
		mustCIDR("169.254.0.0/16"),
		// IPv6 localhost + link-local
		mustCIDR("::1/128"),
		mustCIDR("fc00::/7"),
		mustCIDR("fe80::/10"),
	}

	for _, block := range privateBlocks {
		if block.Contains(parsedIP) {
			return true
		}
	}
	return false
}

func mustCIDR(cidr string) *net.IPNet {
	_, block, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("invalid CIDR: " + cidr)
	}
	return block
}

// GetClientIP extracts the best IP address from headers or RemoteAddr
func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" && !IsPrivate(ip) {
				return ip
			}
		}
	}

	xrip := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xrip != "" && !IsPrivate(xrip) {
		return xrip
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && !IsPrivate(ip) {
		return ip
	}

	return ""
}
