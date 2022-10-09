package tcp

import (
	"bytes"
	"testing"
)

type testHandler struct {
	messageList [][]byte
	decodeError bool
}

func (self *testHandler) Process(message []byte) {
	self.messageList = append(self.messageList, message)
}

func (self *testHandler) DecodeError() {
	self.decodeError = true
}

func Test_decode(t *testing.T) {
	pkg := PackageDecoder{}
	pkg.buffer = bytes.NewBuffer(nil)

	handler := &testHandler{}

	////////////////////////////////////////
	// 完整buffer
	// 1. [0x00, 0x01, 0x02, 0x03, 0x04]
	// 2. [0x10, 0x11, 0x12]
	buf1 := []byte{0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04, 0x00, 0x03, 0x10, 0x11, 0x12}
	pkg.Decode(buf1, handler)

	if (pkg.lastLen != 0) || (pkg.buffer.Len() != 0) {
		t.Errorf("协议拆分出错")
	}

	if len(handler.messageList) != 2 {
		// t.Error("解析数量不对:[标准]")
		t.Errorf("协议数量不对[目标%d:实际%d]", 2, len(handler.messageList))
	}

	message := handler.messageList[0]

	// 测试第1条协议
	if len(message) != 5 {
		t.Errorf("第%d条协议长度出错", 1)
	}

	if (message[0] != 0x00) ||
		(message[1] != 0x01) ||
		(message[2] != 0x02) ||
		(message[3] != 0x03) ||
		(message[4] != 0x04) {
		t.Errorf("第%d条协议数据出错", 1)
	}

	message = handler.messageList[1]

	// 测试第2条协议
	if len(message) != 3 {
		t.Errorf("第%d条协议长度出错", 2)
	}

	if (message[0] != 0x10) ||
		(message[1] != 0x11) ||
		(message[2] != 0x12) {
		t.Errorf("第%d条协议数据出错", 2)
	}

	handler.messageList = [][]byte{}

	////////////////////////////////////////
	// 测试不完整协议
	// 1. [0x00, 0x01, 0x02]
	buf2 := []byte{0x00, 0x05, 0x00, 0x01, 0x02}
	pkg.Decode(buf2, handler)

	if (pkg.lastLen != 5) || (pkg.buffer.Len() != 3) {
		t.Errorf("协议拆分出错")
	}

	message = pkg.buffer.Bytes()
	if (message[0] != 0x00) ||
		(message[1] != 0x01) ||
		(message[2] != 0x02) {
		t.Errorf("第%d条协议数据出错", 3)
	}

	// 追加协议
	// [0x03, 0x04]
	buf3 := []byte{0x03, 0x04}
	pkg.Decode(buf3, handler)

	if (pkg.lastLen != 0) || (pkg.buffer.Len() != 0) {
		t.Errorf("协议拆分出错")
	}

	if len(handler.messageList) != 1 {
		t.Errorf("协议数量不对[目标%d:实际%d]", 1, len(handler.messageList))
	}

	// 测试第3条协议
	message = handler.messageList[0]
	if len(message) != 5 {
		t.Errorf("第%d条协议长度出错", 3)
	}

	if (message[0] != 0x00) ||
		(message[1] != 0x01) ||
		(message[2] != 0x02) ||
		(message[3] != 0x03) ||
		(message[4] != 0x04) {
		t.Errorf("第%d条协议数据出错", 3)
	}

	////////////////////////////////////////
	// 测试协议解析失败
	buf4 := []byte{0x00, 0x00, 0x01, 0x02, 0x03}
	pkg.Decode(buf4, handler)

	if (pkg.lastLen != 0) || (pkg.buffer.Len() != 0) {
		t.Errorf("协议拆分出错")
	}

	if handler.decodeError != true {
		t.Error("错误协议识别错误")
	}
}
