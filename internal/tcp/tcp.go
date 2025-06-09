package tcp

import (
	"fmt"
	"net"
)

func StartTCPServer(deps Dependencies) {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(fmt.Sprintf("TCP server failed: %v", err))
	}

	fmt.Println("TCP server running on port 9090")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		go HandleConnection(conn, deps)
	}
}
