package common

import (
	"os"
	"bufio"
	"fmt"
	"net"
	"encoding/binary"
	"math/rand"
	"crypto/md5"
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
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	f.Seek(offset, 0)

	_, err = f.Write(chunk)
	if err != nil {
		panic(err)
	}
}

// Read one chunk from a file
func ReadChunk(filename string, chunkSize, offset int64) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(offset, 0)
	chunk := []byte{}
	buf := make([]byte, chunkSize)
	reader := bufio.NewReader(file)
	var byteSize int64
	for byteSize < chunkSize {
		n, err := reader.Read(buf)
		if n > 0 {
			chunk = append(chunk, buf[:n]...)
			byteSize += int64(n)
		}

		if err != nil {
			fmt.Println("ReadChunk:", err)
			break
		}
	}
	return chunk
}

// Write message to socket
func WriteSocket(conn net.Conn, msg []byte) {
	msgLen := uint64(len(msg))
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, msgLen)

	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	// Write header
	writer.Write(header)
	// Write body message
	writer.Write(msg)
}

// Read message from socket
func ReadSocket(conn net.Conn) []byte {
	// Using buffer io
	reader := bufio.NewReader(conn)

	// First read header
	header := []byte{}
	headerBuf := make([]byte, 8)
	h := 0
	for h < 8 {
		n, err := reader.Read(headerBuf)
		fmt.Println("ReadSocket header loop:", n)
		if n > 0 {
			h += n
			header = append(header, headerBuf[:n]...)
		}
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	msgLen := binary.BigEndian.Uint64(header)
	fmt.Println("ReadSocket header:", msgLen)

	// Second read body
	msg := []byte{}
	buf := make([]byte, 1024)
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

// Calc hash by MD5
func CalcMD5(data []byte) [16]byte {
	return md5.Sum(data)
}
