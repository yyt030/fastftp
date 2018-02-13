package common

import (
	"os"
	"bufio"
	"fmt"
	"net"
	"encoding/binary"
	"math/rand"
)

// Create one gaven fileSize file with null
func CreateNullFile(filename string, fileSize int64) *os.File {
	if fileSize <= 0 {
		return nil
	}
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(fileSize-1, 0)
	file.Write([]byte{0})

	return file
}

// Write chunk to a file
func WriteChunk(filename string, chunk []byte, offset int64) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(offset, 0)
	_, err = file.Write(chunk)
	if err != nil {
		panic(err)
	}
}

// Read one chunk from a file
func ReadChunk(filename string, chunkSize, offset int64) []byte {
	chunk := make([]byte, chunkSize)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(offset, 0)
	reader := bufio.NewReader(file)
	byteSize, err := reader.Read(chunk)
	if err != nil {
		panic(err)
	}

	fmt.Println("ReadChunk: read byteSize:", byteSize)
	return chunk
}

// Write message to socket
func WriteSocket(conn net.Conn, msg []byte) {
	msgLen := uint64(len(msg))
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, msgLen)

	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	writer.Write(header)
	writer.Write(msg)
}

// Read message from socket
func ReadSocket(conn net.Conn) []byte {
	header := make([]byte, 8)
	reader := bufio.NewReader(conn)
	_, err := reader.Read(header)
	if err != nil {
		panic(err)
	}
	msgLen := binary.BigEndian.Uint64(header)

	buf := make([]byte, 1024)
	msg := []byte{}
	var byteSize uint64
	for uint64(byteSize) < msgLen {
		n, err := reader.Read(buf)
		if n > 0 {
			byteSize += uint64(n)
			msg = append(msg, buf[:n]...)
		}

		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return msg
}

// Random int
func RandomSource(count int) []byte {
	out := []byte{}
	for i := 0; i < count; i++ {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, rand.Uint64())
		out = append(out, b...)
	}
	return out
}

// Write buf to file
func WriteToFile(filename string, buf []byte) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.Write(buf)
	defer w.Flush()
}
