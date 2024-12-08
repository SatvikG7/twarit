package server

import (
	"fmt"
	"net/http"

	"github.com/satvikg7/twarit/internal/router"
	"github.com/satvikg7/twarit/pkg/qr"
)

func SetupServer() {
	router.Route()
}

func Listen(port string) {
	var loc string = ":"
	loc = loc + port
	qr.GenerateQR(port)
	fmt.Println("	    Scan the QR")
	fmt.Println("		OR")
	fmt.Print("    Visit http://", qr.GetOutboundIP()+loc)
	fmt.Println()
	println("          Type 'q' to quit")
	http.ListenAndServe(loc, nil)
}
