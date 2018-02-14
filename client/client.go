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
	const count = 20
	const chunkCount = 1

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

	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chunkSize := int(math.Ceil(float64(fi.Size()) / float64(chunkCount)))
	for i := 0; i < chunkCount; i++ {
		f.Seek(int64(i*chunkSize), 0)
		chunk := common.ReadSource(
			bufio.NewReader(f), chunkSize)
		fmt.Printf(">>>[%02d] %x\n", i, chunk[:len(chunk)/8])

		w := bufio.NewWriter(conn)
		common.WriteSocket(w, chunk)
		w.Flush()
	}
}
