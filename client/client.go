package main

import (
	"math"
	"fastftp/common"
	"net"
	"sync"
	"fmt"
)

func main() {
	const chunkCount = 30
	const fileSize = 865075200
	const srcFilename = "testdata/ubuntu-16.04.3-server-amd64.iso"

	var wg sync.WaitGroup
	wg.Add(chunkCount)

	chunkSize := int(math.Ceil(float64(fileSize) / float64(chunkCount)))
	for i := 0; i < chunkCount; i++ {
		go ReadWriteChunk(srcFilename, chunkSize, int64(i*chunkSize), &wg)
	}
	fmt.Println(">>> goroutine waiting...")
	wg.Wait()
}

func ReadWriteChunk(src string, chunkSize int, offset int64, wg *sync.WaitGroup) {
	chunk, newOffset := common.ReadChunkFromFile(src, chunkSize, offset)

	conn, err := net.Dial("tcp", "172.16.66.132:8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	common.WriteChunkToSocket(conn, src, chunk, newOffset)

	//common.ReadSource(conn, 10)
	fmt.Println("goroutine done")

	wg.Done()
}
