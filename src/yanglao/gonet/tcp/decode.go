package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type PackageHandler interface {
	Process(message []byte)
	DecodeError()
}

type PackageDecoder struct {
	lastLen uint16
	buffer  *bytes.Buffer
}

func (self *PackageDecoder) Decode(buf []byte, handle PackageHandler) error {
	self.buffer.Write(buf)

	for {
		// 未有中断的数据
		if self.lastLen == 0 {
			// 缓冲区包长度未读完
			if self.buffer.Len() < 2 {
				return nil
			}

			// 读取包长
			// var packageLength uint16 = 0
			binary.Read(self.buffer, binary.BigEndian, &self.lastLen)

			// 包长度非法
			if self.lastLen == 0 {
				self.buffer.Reset()
				// close
				handle.DecodeError()
				break
			}
		}

		// 数据区未读完
		if int(self.lastLen) > self.buffer.Len() {
			return nil
		}

		// 读取包数据
		srcMessage := self.buffer.Next(int(self.lastLen))
		newMessage := make([]byte, len(srcMessage))
		copy(newMessage, srcMessage)
		self.lastLen = 0

		handle.Process(newMessage)
	}

	return nil
}

func Decode16(cacheBuffer *bytes.Buffer, lastLenPtr *uint32) (message []byte, err error) {
	// 未有中断的数据
	if *lastLenPtr == 0 {
		// 缓冲区包长度未读完
		if cacheBuffer.Len() < 2 {
			return nil, nil
		}

		// 读取包长
		var lastLen uint16 = 0
		binary.Read(cacheBuffer, binary.BigEndian, &lastLen)
		*lastLenPtr = uint32(lastLen)

		// 包长度非法
		if *lastLenPtr == 0 {
			cacheBuffer.Reset()
			return nil, errors.New("Decode16 package length failed.")
		}
	}

	// 数据区未读完
	if int(*lastLenPtr) > cacheBuffer.Len() {
		return nil, nil
	}

	// 读取包数据
	srcMessage := cacheBuffer.Next(int(*lastLenPtr))
	message = make([]byte, len(srcMessage))
	copy(message, srcMessage)
	*lastLenPtr = 0

	return message, nil
}

func Decode32(cacheBuffer *bytes.Buffer, lastLen *uint32) (message []byte, err error) {
	// 未有中断的数据
	if *lastLen == 0 {
		// 缓冲区包长度未读完
		if cacheBuffer.Len() < 4 {
			return nil, nil
		}

		// 读取包长
		binary.Read(cacheBuffer, binary.BigEndian, lastLen)

		// 包长度非法
		if *lastLen == 0 {
			cacheBuffer.Reset()
			return nil, errors.New("Decode32 package length failed.")
		}
	}

	// 数据区未读完
	if int(*lastLen) > cacheBuffer.Len() {
		return nil, nil
	}

	// 读取包数据
	srcMessage := cacheBuffer.Next(int(*lastLen))
	message = make([]byte, len(srcMessage))
	copy(message, srcMessage)
	*lastLen = 0

	return message, nil
}

func NewPackageDecoder() *PackageDecoder {
	return &PackageDecoder{0, bytes.NewBuffer(nil)}
}
