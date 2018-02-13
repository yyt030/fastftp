package main

import (
	"net"
	"io/ioutil"
	"fastftp/common"
)

func main() {
	const filename = "foo.in"
	const count = 20000000

	// create file
	common.WriteToFile(filename, common.RandomSource(count))

	// Read file
	content, err := ioutil.ReadFile("foo.in")
	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	common.WriteSocket(conn, content)
}
