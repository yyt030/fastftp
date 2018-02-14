package main

import (
	"net"
	"fmt"
	"fastftp/common"
)

func main() {
	addr := ":8888"
	createServer(addr)
}

func createServer(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		// Handle conntions in a new goroutine
		go handleRequest(conn)
	}
}

// Handle request connection
func handleRequest(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("welcome connection from %s\n", conn.RemoteAddr().String())
	msg := common.ReadSocket(conn)

	// Save chunk to file
	common.WriteToFile("testdata/foo.in.out", msg, 0)
}
