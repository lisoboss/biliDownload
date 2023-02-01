package media

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Header interface {
	Type() string
	String() string
	Length() uint64
	AddLength(l uint64)
	Load(bytes []byte)
	DataSize() uint8
	Data() []byte
}

type BoxHeader struct {
	size      uint32
	largeSize uint64
	type_     string
	data      []byte
	dataSize  uint8
}

func (h *BoxHeader) Type() string {
	return h.type_
}

func (h *BoxHeader) String() string {
	return fmt.Sprintf("tpye: %s, length: %d", h.type_, h.Length())
}

func (h *BoxHeader) Length() uint64 {
	if h.size == 0 {
		return h.largeSize
	}
	return uint64(h.size)
}

func (h *BoxHeader) DataSize() uint8 {
	return h.dataSize
}

func (h *BoxHeader) Data() []byte {
	return h.data
}

func (h *BoxHeader) AddLength(l uint64) {
	l = l + h.Length()
	if l > math.MaxInt32 {
		if h.size != 0 {
			h.size = 0
			h.data[0] = 0
			h.data[1] = 0
			h.data[2] = 0
			h.data[3] = 0
		}
		h.largeSize = l
		// binary.BigEndian.PutUint64
		h.data[4+0] = byte(l >> 56)
		h.data[4+1] = byte(l >> 48)
		h.data[4+2] = byte(l >> 40)
		h.data[4+3] = byte(l >> 32)
		h.data[4+4] = byte(l >> 24)
		h.data[4+5] = byte(l >> 16)
		h.data[4+6] = byte(l >> 8)
		h.data[4+7] = byte(l)
	} else {
		v := uint32(l)
		h.size = v
		// binary.BigEndian.PutUint32
		h.data[0] = byte(v >> 24)
		h.data[1] = byte(v >> 16)
		h.data[2] = byte(v >> 8)
		h.data[3] = byte(v)
	}
}

func (h *BoxHeader) Load(bytes []byte) {
	h.dataSize = 4
	h.size = binary.BigEndian.Uint32(bytes[:h.dataSize])
	if h.size == 0 {
		h.largeSize = binary.BigEndian.Uint64(bytes[h.dataSize : h.dataSize+8])
		h.dataSize += 8
	}
	h.type_ = string(bytes[h.dataSize : h.dataSize+4])
	h.dataSize += 4
	h.data = bytes[:h.dataSize]
}

type FullBoxHeader struct {
	BoxHeader
	extend []byte
}

func (h *FullBoxHeader) Load(bytes []byte) {
	h.BoxHeader.Load(bytes)
	h.extend = bytes[h.dataSize : h.dataSize+4]
	h.dataSize += 4
	h.data = bytes[:h.dataSize]
}
