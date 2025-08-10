package iputils

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func GetPublicIP() string {
	// var resp *http.Response
	// err := fmt.Errorf("it is error")
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Always use "tcp4" to force IPv4
			return (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, "tcp4", addr)
		},
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://ifconfig.me/ip")
	if err != nil {
		resp, err = http.Get("https://ipinfo.io/ip")
		if err != nil {
			resp, err = http.Get("https://api.ipify.org")
		}
		log.Println("Error fetching IP:", err)
		if err != nil {
			return "8.8.8.8"
		}
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return "8.8.8.8"
	}

	log.Println("Your public IP is:", string(ip))
	return string(ip)
}

func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, block := range privateBlocks {
		_, cidr, _ := net.ParseCIDR(block)
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
