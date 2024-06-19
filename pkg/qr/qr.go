package qr

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"log"
	"net"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// RenderString as a QR code
func RenderString(s string) {
	q, err := qrcode.New(s, qrcode.Highest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q.ToSmallString(false))
}

func GenerateQR(port string) {
	RenderString("http://" + GetOutboundIP() + ":" + port)
}
