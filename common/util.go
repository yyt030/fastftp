package common

import (
	"os"
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

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

func ReadChunkFromFile(filename string, chunkSize int, startOffset int64) ([]byte, int64) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Seek(startOffset, 0)

	chunk := []byte{}
	buf := make([]byte, chunkSize)
	bytesRead := 0

	for {
		n, err := f.Read(buf)
		if n > 0 {
			chunk = append(chunk, buf[:n]...)
			bytesRead += n
		}
		if err != nil || bytesRead >= chunkSize {
			break
		}
	}

	return chunk, startOffset
}

func WriteChunkToFile(filename string, chunk []byte, offset int64) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteAt(chunk, offset)
	if err != nil {
		panic(err)
	}
}

func WriteChunkToSocket(conn net.Conn, filename string, chunk []byte, offset int64) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	msgLen := make([]byte, 8)    // 报文长度
	msgType := make([]byte, 1)   // 报文类型
	fileName := make([]byte, 64) // 发送文件名字，含路径
	copy(fileName[:], f.Name())

	fileSize := make([]byte, 8) // 发送源文件大小
	binary.BigEndian.PutUint64(fileSize, uint64(info.Size()))

	fileHash := make([]byte, 16) // 文件大小hash
	fileType := make([]byte, 1)  // 文件类型：正常，压缩
	chunkSize := make([]byte, 8) // 块大小
	binary.BigEndian.PutUint64(chunkSize, uint64(len(chunk)))
	chunkOffset := make([]byte, 8) // 块所在文件的偏移量
	binary.BigEndian.PutUint64(chunkOffset, uint64(offset))

	// 发送的buf
	buf := []byte{}
	buf = append(buf, msgLen...)
	buf = append(buf, msgType...)
	buf = append(buf, fileName...)
	buf = append(buf, fileSize...)
	buf = append(buf, fileHash...)
	buf = append(buf, fileType...)
	buf = append(buf, chunkSize...)
	buf = append(buf, chunkOffset...)
	buf = append(buf, chunk...)

	binary.BigEndian.PutUint64(buf[:8], uint64(len(buf)-8)) // 赋值报文长度
	// Write socket
	w := bufio.NewWriter(conn)
	w.Write(buf)
	w.Flush()
}

func ReadChunkFromSocket(r io.Reader) ([]byte, uint64) {
	header := ReadSource(r, 114)
	// message header
	chunkSize := header[98:106]
	chunkOffset := header[106:114]

	// message body
	msgBody := ReadSource(r, int(binary.BigEndian.Uint64(chunkSize)))
	return msgBody, binary.BigEndian.Uint64(chunkOffset)
}
