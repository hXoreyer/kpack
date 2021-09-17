package kpack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type UnPackFunc func(r io.Reader) interface{}

type Kpack struct {
	Length int16
}

func Pack(buf *bytes.Buffer, data interface{}) {
	v := reflect.ValueOf(data).Elem()
	t := reflect.TypeOf(data).Elem()
	len := v.NumField()
	for i := 0; i < len; i++ {
		if t.Field(i).Name == "Kpack" {
			binary.Write(buf, binary.BigEndian, v.Field(i).Interface().(Kpack).Length)
			continue
		}
		if t.Field(i).Tag.Get("ksize") == "true" {
			v.Field(i).SetInt(int64(sizeof(v.Field(i + 1))))
		}
		binary.Write(buf, binary.BigEndian, v.Field(i).Interface())
	}
}
func UnPack(scanner *bufio.Scanner, bufSize int, up UnPackFunc) ([]interface{}, int) {
	var i int
	s := make([]interface{}, bufSize)
	for scanner.Scan() {
		s[i] = up(bytes.NewBuffer(scanner.Bytes()))
		i++
	}
	return s, i
}

func ScanPack(buf *bytes.Buffer) (*bufio.Scanner, int) {
	var in int16
	sz := buf.Bytes()
	/*
		fmt.Println(sz[:2])
		in = len(buf.Bytes())/int(int16(sz[0])<<8|int16(sz[1])) + 1
	*/
	bf := bytes.NewBuffer(sz[:2])
	binary.Read(bf, binary.BigEndian, &in)
	fmt.Println(in)
	scanner := bufio.NewScanner(buf)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF {
			length := int16(0)
			binary.Read(bytes.NewBuffer(data), binary.BigEndian, &length)
			if int(length) <= len(data) {
				return int(length), data[:int(length)], nil
			}
		}
		return
	})
	return scanner, int(in)
}
