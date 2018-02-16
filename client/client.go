package main

import (
	"net"
	"fastftp/common"
	"os"
	"math"
	"bufio"
	"sync"
	"fmt"
)

func main() {
	const filename = "testdata/foo.in"
	const count = 2000
	const chunkCount = 13

	// create file
	common.WriteToFile(filename, common.RandomSource(count), 0)

	// Put Source
	createPipeline(filename, chunkCount)

}

func createPipeline(filename string, chunkCount int) {
	var wg sync.WaitGroup
	wg.Add(chunkCount)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	chunkSize := int(math.Ceil(float64(fi.Size()) / float64(chunkCount)))
	for i := 0; i < chunkCount; i++ {
		SendChunk(i, f, fi, chunkSize, &wg)
	}
	wg.Wait()
}

func SendChunk(seq int, f *os.File, fi os.FileInfo, chunkSize int, wg *sync.WaitGroup) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	offset := seq * chunkSize

	f.Seek(int64(offset), 0)
	chunk := common.ReadSource(
		bufio.NewReader(f), chunkSize)

	common.WriteSocket(fi, conn, chunk, uint64(offset))

	fmt.Printf("goroutine %d done\n", seq)
	wg.Done()
}
