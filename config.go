package config

const MsgLength = 8

type ReqMsg struct {
	Length int64   // 报文长度
	MsgFlag   [1]byte // 报文类型
	FileInfo       // 文件信息
	ChunkInfo      // 块信息
}

// 文件信息
type FileInfo struct {
	FileName [64]byte // 发送文件名字，含路径
	FileSize int64    // 文件大小
	FileHash [16]byte // 文件大小hash
	FileType [1]byte  // 文件类型：正常，压缩
}

// 块信息
type ChunkInfo struct {
	ChunkSize int64   // 块大小
	ChunkSeq  [4]byte // 块在文件中的序号
	Chunk     *[]byte // chunk包信息
}
