package config

import "fmt"

const MsgLength = 8

type ReqMsg struct {
	Length   [8]byte  // 报文长度
	Flag     [1]byte  // 报文类型
	Filename [64]byte // 发送文件名字， 含路径
	FileSize [8]byte  // 文件大小
	FileHash [16]byte // 文件大小hash
	FileType [1]byte  // 文件类型：正常，压缩
	Chunk             // 块信息
}

type Chunk struct {
	ChunkSize [8]byte // 块大小
	ChunkSeq  [4]byte // 块在文件中的序号
	Content   *[]byte // chunk包信息
}

func NewReqMsg() *ReqMsg {
	return &ReqMsg{}
}

func (req *ReqMsg) ReadSource() *ReqMsg {
	newReq := NewReqMsg()
	fmt.Println(">>>", newReq)
	return req
}
