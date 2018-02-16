package common

import (
	"os"
	"bufio"
	"fmt"
	"encoding/binary"
	"math/rand"
	"io"
	"net"
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

func WriteSocket(fi os.FileInfo, conn net.Conn, chunk []byte, ) {
	var length uint64            // 报文长度
	msgFlag := [1]byte{0x00}     // 报文类型
	fileName := [64]byte{}       //[64]byte 发送文件名字，含路径
	fileSize := int64(fi.Size()) // int64   文件大小
	fileHash := [16]byte{}       // 文件大小hash
	fileType := [1]byte{0x00}    // 文件类型：正常，压缩
	chunkSize := len(c)          // int64   // 块大小
	chunkSeq := [1]byte{0x00}    // 块在文件中的序号

	length = 1 + 64 + 8 + 16 + 1 + 8 + 1 + uint64(len(c))
	fmt.Println("WriteFileToSocket header length:", length, fi.Name())

	// Set header buffer
	msg := []byte{}
	lengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lengthBytes, length)
	msg = append(msg, lengthBytes[:]...)
	msg = append(msg, msgFlag[:]...)
	msg = append(msg, fileName[:]...)
	fileSizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(fileSizeBytes, uint64(fileSize))
	msg = append(msg, fileSizeBytes...)
	msg = append(msg, fileHash[:]...)
	msg = append(msg, fileType[:]...)
	chunkSizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(chunkSizeBytes, uint64(chunkSize))
	msg = append(msg, chunkSizeBytes...)
	//chunkSeqBytes := make([]byte, 1)
	//binary.BigEndian.PutUint64(chunkSeqBytes, uint64(chunkSeq))
	msg = append(msg, chunkSeq[:]...)
	msg = append(msg, chunk...)

	w := bufio.NewWriter(conn)
	w.Write(msg)
	w.Flush()
}

func ReadSocket(r io.Reader) []byte {
	header := ReadSource(r, 107)
	// message header
	length := header[:8]
	msgType := header[8:9]
	fileName := header[9:9+64]
	fileSize := header[73:73+8]
	fileHash := header[81:81+16]
	fileType := header[97:98]
	chunkSize := header[98:98+8]
	chunkSeq := header[106:107]

	// print message header
	fmt.Println("-----------------------")
	fmt.Printf("length:%x %d\n", length, binary.BigEndian.Uint64(length))
	fmt.Println("msgType:", msgType)
	fmt.Println("filename:", fileName)
	fmt.Println("fileSize:", fileSize, binary.BigEndian.Uint64(fileSize))
	fmt.Printf("fileHash:%x\n", fileHash)
	fmt.Println("fileType:", fileType)
	fmt.Printf("chunkSize:%x\n", chunkSize, )
	fmt.Println("chunkSeq:", chunkSeq)
	//fmt.Printf("chunk:%x\n", ReadSource(r, int(binary.BigEndian.Uint64(chunkSize))))

	// read message body
	msgBody := ReadSource(r, int(binary.BigEndian.Uint64(chunkSize)))
	return msgBody
}
