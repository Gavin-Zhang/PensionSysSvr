package tcp

import (
	"bytes"
	"testing"
)

func Benchmark_decode(b *testing.B) {
	buf1 := []byte{0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04, 0x00, 0x03, 0x10, 0x11, 0x12}
	buf2 := []byte{0x00, 0x05, 0x00, 0x01, 0x02}
	buf3 := []byte{0x03, 0x04}

	pkg := PackageDecoder{}
	pkg.buffer = bytes.NewBuffer(nil)

	handler := &testHandler{}

	for i := 0; i < b.N; i++ {
		pkg.Decode(buf1, handler)
		pkg.Decode(buf2, handler)
		pkg.Decode(buf3, handler)
	}
}
