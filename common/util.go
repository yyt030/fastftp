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
func WriteToFile(filename string, chunk []byte, offset uint64) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	f.Seek(int64(offset), 0)

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

func WriteSocket(fi os.FileInfo, conn net.Conn, chunk []byte, offset uint64) {
	var length uint64              // 报文长度
	msgFlag := [1]byte{0x00}       // 报文类型
	fileName := [64]byte{}         //[64]byte 发送文件名字，含路径
	fileSize := int64(fi.Size())   // int64   文件大小
	fileHash := [16]byte{}         // 文件大小hash
	fileType := [1]byte{0x00}      // 文件类型：正常，压缩
	chunkSize := len(chunk)        // int64   // 块大小
	chunkSeq := make([]byte, 2)    // 块在文件中的序号
	chunkOffset := make([]byte, 8) // 块所在文件的偏移量

	length = 1 + 64 + 8 + 16 + 1 + 8 + 1 + uint64(len(chunk))
	fmt.Printf("WriteFileToSocket header length:[%d], name:[%s], chunkSeq:[%d], chunk:%x\n", length, fi.Name(), offset, chunk[:len(chunk)/8])

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
	msg = append(msg, chunkSeq[:]...)
	binary.BigEndian.PutUint64(chunkOffset, offset)
	msg = append(msg, chunkOffset[:]...)
	msg = append(msg, chunk...)

	w := bufio.NewWriter(conn)
	w.Write(msg)
	w.Flush()
}

func ReadSocket(r io.Reader) ([]byte, uint64) {
	header := ReadSource(r, 116)
	// message header
	chunkSize := header[98:106]
	chunkOffset := header[108:116]

	// read message body
	msgBody := ReadSource(r, int(binary.BigEndian.Uint64(chunkSize)))
	fmt.Printf("message body:%x\n", msgBody[:len(msgBody)/8])
	return msgBody, binary.BigEndian.Uint64(chunkOffset)
}

func PrintMsgHeader(header []byte) {
	length := header[:8]
	msgType := header[8:9]
	fileName := header[9:73]
	fileSize := header[73:81]
	fileHash := header[81:97]
	fileType := header[97:98]
	chunkSize := header[98:106]
	chunkSeq := header[106:108]
	chunkOffset := header[108:116]

	// print message header
	fmt.Printf("length:%x %d\n", length, binary.BigEndian.Uint64(length))
	fmt.Println("msgType:", msgType)
	fmt.Println("filename:", fileName)
	fmt.Println("fileSize:", fileSize, binary.BigEndian.Uint64(fileSize))
	fmt.Printf("fileHash:%x\n", fileHash)
	fmt.Println("fileType:", fileType)
	fmt.Printf("chunkSize:%x\n", chunkSize, )
	fmt.Println("chunkSeq:", chunkSeq)
	fmt.Println("chunkOffset:", chunkOffset)
}
