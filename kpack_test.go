package kpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"
	"time"
)

func BenchmarkKpack(b *testing.B) {
	// TODO: Initialize
	pack := Package{
		Msg: []byte("现在时间是:" + time.Now().Format("2006-01-02 15:04:05")),
		Sex: []byte("男"),
		Age: 23,
	}
	pack.Length = int16(SizeOf(pack))
	buf := new(bytes.Buffer)
	// 写入四次，模拟TCP粘包效果
	for i := 0; i < 4; i++ {
		Pack(buf, &pack)
	}
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		buft := new(bytes.Buffer)
		buft.Write(buf.Bytes())
		scanner, size := ScanPack(buft)
		//解包
		UnPack(scanner, size, func(r io.Reader) interface{} {
			p := &Package{}
			binary.Read(r, binary.BigEndian, &p.Length)
			binary.Read(r, binary.BigEndian, &p.MsgLen)
			p.Msg = make([]byte, p.MsgLen)
			binary.Read(r, binary.BigEndian, &p.Msg)
			binary.Read(r, binary.BigEndian, &p.SexLen)
			p.Sex = make([]byte, p.SexLen)
			binary.Read(r, binary.BigEndian, &p.Sex)
			binary.Read(r, binary.BigEndian, &p.Age)
			return p
		})
	}
}

type Package struct {
	Kpack
	MsgLen int16 `ksize:"true"`
	Msg    []byte
	SexLen int16 `ksize:"true"`
	Sex    []byte
	Age    int16
}

func (p *Package) String() string {
	return fmt.Sprintf("length:%d msglen:%d msg:%s sexlen:%d sex:%s age:%d", p.Length, p.MsgLen, p.Msg, p.SexLen, p.Sex, p.Age)
}
