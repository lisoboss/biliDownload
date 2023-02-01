package media

import (
	"bytes"
	"encoding/binary"
)

func IsContainerBox(b *Box) bool {
	switch b.Type() {
	case "mdia",
		"minf",
		"moof",
		"moov",
		"mvex",
		"stbl",
		"traf",
		"trak":
		return true
	default:
		return false
	}
}

func BinaryBigEndianInt32(b []byte, data any) {
	_ = binary.Read(bytes.NewBuffer(b), binary.BigEndian, data)
}

func BinaryBigEndianPutInt32(data int32) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes()
}
