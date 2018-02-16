package main

import (
	"net"
	"fastftp/common"
	"os"
	"math"
	"fmt"
	"bufio"
)

func main() {
	const filename = "testdata/foo.in"
	const count = 20000
	const chunkCount = 15

	// create file
	common.WriteToFile(filename, common.RandomSource(count), 0)

	// Put Source
	createPipeline(filename, chunkCount)

}

func createPipeline(filename string, chunkCount int) {
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
		conn, err := net.Dial("tcp", ":8888")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		offset := i * chunkSize
		f.Seek(int64(offset), 0)
		chunk := common.ReadSource(
			bufio.NewReader(f), chunkSize)
		fmt.Printf(">>>[%02d] %x\n", i, chunk[:30])

		w := bufio.NewWriter(conn)
		common.WriteSocket(fi, conn, chunk, uint64(offset))
		w.Flush()
	}
}
