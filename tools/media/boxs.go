package media

import (
	"encoding/binary"
	"fmt"
)

type ChunkOffsetBox struct {
	FullBox
	EntryCount  uint32
	ChunkOffset []uint32
}

func (b *ChunkOffsetBox) String() string {
	return fmt.Sprintf("%s\nEntryCount: %d, ChunkOffset: %v", b.FullBox.String(), b.EntryCount, b.ChunkOffset)
}

func (b *ChunkOffsetBox) Load(bytes []byte) []byte {
	var (
		size        = 4
		i    uint32 = 0
	)

	bytes = b.FullBox.Load(bytes)
	// load EntryCount
	b.EntryCount = binary.BigEndian.Uint32(b.data[:size])
	size += 4
	// load ChunkOffset
	for ; i < b.EntryCount; i++ {
		b.ChunkOffset = append(b.ChunkOffset, binary.BigEndian.Uint32(b.data[size:size+4]))
		size += 4
	}
	return bytes
}

func (b *ChunkOffsetBox) Dump() (bytes []byte) {
	bytes = append(bytes, b.head.Data()...)
	binary.BigEndian.PutUint32(bytes, b.EntryCount)
	for _, v := range b.ChunkOffset {
		binary.BigEndian.PutUint32(bytes, v)
	}
	return
}

type TrackRunBox struct {
	FullBox
	SampleCount uint32
	DataOffset  int32
}

func (b *TrackRunBox) String() string {
	return fmt.Sprintf("%s\nSampleCount: %d, DataOffset: %d", b.FullBox.String(), b.SampleCount, b.DataOffset)
}

func (b *TrackRunBox) Load(bytes []byte) []byte {
	bytes = b.FullBox.Load(bytes)
	// load SampleCount
	b.SampleCount = binary.BigEndian.Uint32(b.data[:4])
	BinaryBigEndianInt32(b.data[4:8], &b.DataOffset)
	return bytes
}

type MovieFragmentHeaderBox struct {
	FullBox
	sequenceNumber uint32
}

func (b *MovieFragmentHeaderBox) String() string {
	return fmt.Sprintf("%s\nsequenceNumber: %d", b.FullBox.String(), b.sequenceNumber)
}

func (b *MovieFragmentHeaderBox) Load(bytes []byte) []byte {
	bytes = b.FullBox.Load(bytes)
	// load sequenceNumber
	b.sequenceNumber = binary.BigEndian.Uint32(b.data[:4])
	return bytes
}

func (b *MovieFragmentHeaderBox) SetSequenceNumber(v uint32) {
	b.sequenceNumber = v
	b.data[0] = byte(v >> 24)
	b.data[1] = byte(v >> 16)
	b.data[2] = byte(v >> 8)
	b.data[3] = byte(v)
}

func (b *MovieFragmentHeaderBox) SequenceNumber() uint32 {
	return b.sequenceNumber
}

type TrackHeaderBox struct {
	FullBox
	trackID      uint32
	trackIDIndex int
}

func (b *TrackHeaderBox) String() string {
	return fmt.Sprintf("%s\ntrackID: %d", b.FullBox.String(), b.trackID)
}

func (b *TrackHeaderBox) Load(bytes []byte) []byte {
	bytes = b.FullBox.Load(bytes)
	// load trackID
	fmt.Printf("TrackHeaderBox Version %d", b.Version())
	if b.Version() == 1 {
		b.trackIDIndex = 16
	} else {
		b.trackIDIndex = 8
	}
	b.trackID = binary.BigEndian.Uint32(b.data[b.trackIDIndex : b.trackIDIndex+4])
	return bytes
}

func (b *TrackHeaderBox) SetTrackID(v uint32) {
	b.trackID = v
	b.data[b.trackIDIndex+0] = byte(v >> 24)
	b.data[b.trackIDIndex+1] = byte(v >> 16)
	b.data[b.trackIDIndex+2] = byte(v >> 8)
	b.data[b.trackIDIndex+3] = byte(v)
}

func (b *TrackHeaderBox) TrackID() uint32 {
	return b.trackID
}
