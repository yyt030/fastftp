package main

import (
	"net"
	"fastftp/common"
	"os"
	"math"
	"fmt"
)

func main() {
	const filename = "testdata/foo.in"
	const count = 1

	// create file
	common.WriteChunk(filename, common.RandomSource(count), 0)

	createPipeline(filename, 100000)

}

func createPipeline(filename string, chunkSize int64) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chunkCount := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))
	fmt.Println(">>>", chunkCount)
	for i := 0; i < chunkCount; i++ {
		chunk := common.ReadChunk(filename, chunkSize, int64(i)*chunkSize)
		fmt.Println(">>>", chunk)
		common.WriteSocket(conn, chunk)
	}
}
