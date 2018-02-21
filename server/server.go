package main

import (
	"net"
	"fmt"
	"fastftp/common"
)

func main() {
	addr := "0.0.0.0:8888"
	fmt.Println("Listening on:", addr)
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

		handleRequest(conn)
	}
}

// Handle request connection
func handleRequest(conn net.Conn) {
	defer conn.Close()
	msg, offset := common.ReadChunkFromSocket(conn)
	// Save chunk to file
	common.WriteChunkToFile("testdata/foo.out", msg, int64(offset))
}
