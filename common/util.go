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
	//req := config.NewReqMsg()
	//req.Length =
	//copy(req.Filename[:], []byte("foo.in.out"))
	//req.FileSize = uint64(len(chunk))
	//req.FileHash = [16]byte{}k]]]]][
	//req.fileType = [1]byte{}
	//req.Chunk = &chunk
	//req.ChunkSize = uint64(len(chunk))

	length := make([]byte, 8) // 报文长度
	flag := make([]byte, 1)   // 报文类型
	//filename := make([]byte, 64)    // 发送文件名字， 含路径
	//fileSize := make([]byte, 8)     // 文件大小
	//fileHash := make([]byte, 16)    // 文件大小hash
	//fileType := make([]byte, 1)     // 文件类型：正常，压缩
	chunkSize := uint64(len(chunk)) // 块大小
	//chunkSeq := make([]byte, 4)     // 块在文件中的序号

	binary.BigEndian.PutUint64(length, chunkSize+1+64+8+16+1+4)

	// Write header
	w.Write(length)
	// Chunk sequeue
	w.Write(flag)
	// Write body message
	w.Write(chunk)
}

func ReadSocket(r io.Reader) []byte {
	header := ReadSource(r, 8)
	msgLen := binary.BigEndian.Uint64(header)
	fmt.Println("ReadSocket header:", msgLen)
	return ReadSource(r, int(msgLen))
}
