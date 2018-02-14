package common

import (
	"os"
	"bufio"
	"fmt"
	"encoding/binary"
	"math/rand"
	"io"
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

// Write chunk to a file
func WriteToFile(filename string, chunk []byte, offset int64) {
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

// Read source data
func ReadSource(r io.Reader, chuckSize int) []byte {
	chunk := []byte{}
	buf := make([]byte, chuckSize)
	bytesRead := 0
	for {
		n, err := r.Read(buf)
		bytesRead += n
		if n > 0 {
			chunk = append(chunk, buf[:n]...)
		}
		if err != nil ||
			(chuckSize != -1 && bytesRead >= chuckSize) {
			break
		}
	}
	return chunk
}

func WriteSocket(w io.Writer, chunk []byte) {
	chunkSize := uint64(len(chunk))
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, chunkSize)
	// Write header
	w.Write(header)
	// Write body message
	w.Write(chunk)
}

func ReadSocket(r io.Reader) []byte {
	header := ReadSource(r, 8)
	msgLen := binary.BigEndian.Uint64(header)
	fmt.Println("ReadSocket header:", msgLen)
	return ReadSource(r, int(msgLen))
}
