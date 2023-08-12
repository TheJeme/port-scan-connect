package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"net/http"
)

func main() {
	http.HandleFunc("/", getPing)
	fmt.Printf("Server started on %s:7229\n", getLocalIP())
	err := http.ListenAndServe(":7229", nil)
	if err != nil {
		fmt.Println("Error starting server : ", err)
	}
}

func getPing(w http.ResponseWriter, r *http.Request) {
	requestIp := getRequestIP(r)
	if requestIp == "" {
		fmt.Printf("Ping\n")
		} else {
		fmt.Printf("Ping from %s\n", requestIp)
	}
	io.WriteString(w, "Ping\n")
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return fmt.Sprintf("%s", localAddr.IP)
}

func getRequestIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String()
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		localIp := getLocalIP()
		if (ip == localIp) {
			return fmt.Sprintf("%s (you)", localIp)
		}
		if ip == "::1" {
			return fmt.Sprintf("%s (you)", localIp)
		}
		return ip
	}

	return ""
}