package media

import (
	"encoding/binary"
	"fmt"
)

type Box struct {
	head             Header
	data             []byte
	containerBoxList []*Box
	parentBox        *Box
}

func (b *Box) String() string {
	return fmt.Sprintf("%s", b.head)
}

func (b *Box) Print(prefix string) {
	fmt.Printf("%s%s\n", prefix, b)
	prefix = prefix + "    "
	for _, box := range b.containerBoxList {
		box.Print(prefix)
	}
}

func (b *Box) Type() string {
	return b.head.Type()
}

func (b *Box) ParentBox() *Box {
	return b.parentBox
}

func (b *Box) Length() uint64 {
	return b.head.Length()
}

func (b *Box) AddLength(l uint64) {
	b.head.AddLength(l)
	if b.parentBox != nil {
		b.parentBox.AddLength(l)
	}
}

func (b *Box) Load(bytes []byte) []byte {
	b.head = new(BoxHeader)
	b.head.Load(bytes)
	b.data = bytes[b.head.DataSize():b.head.Length()]
	if IsContainerBox(b) {
		b.LoadContainerBox()
	}
	return bytes[b.head.Length():]
}

func (b *Box) Dump() (bytes []byte) {
	bytes = append(bytes, b.head.Data()...)
	if IsContainerBox(b) {
		for _, box := range b.containerBoxList {
			//fmt.Printf("Dump Container: %s\n", box)
			bytes = append(bytes, box.Dump()...)
		}
	} else {
		bytes = append(bytes, b.data...)
	}
	return
}

func (b *Box) LoadContainerBox() {
	bytes := b.data
	for len(bytes) > 0 {
		box := new(Box)
		bytes = box.Load(bytes)
		box.parentBox = b
		b.containerBoxList = append(b.containerBoxList, box)
	}
}

func (b *Box) Is(type_ string) bool {
	return b.head.Type() == type_
}

func (b *Box) GetBox(type_ string) *Box {
	if b.Is(type_) {
		return b
	}
	for _, box := range b.containerBoxList {
		box = box.GetBox(type_)
		if box != nil {
			return box
		}
	}
	return nil
}

func (b *Box) GetBoxList(type_ string) (rets []*Box) {
	if b.Is(type_) {
		rets = append(rets, b)
	}
	for _, box := range b.containerBoxList {
		box = box.GetBox(type_)
		if box != nil {
			rets = append(rets, box)
		}
	}
	return
}

func (b *Box) AddContainerBox(box *Box) {
	i := 0
	for i2, b2 := range b.containerBoxList {
		if b2.Type() == box.Type() {
			i = i2 + 1
		}
	}
	if i == len(b.containerBoxList) {
		b.containerBoxList = append(b.containerBoxList, box)
	} else {
		b.containerBoxList = append(b.containerBoxList[:i+1], b.containerBoxList[i:]...)
		b.containerBoxList[i] = box
	}
}

type FullBox struct {
	Box
}

func (b *FullBox) Load(bytes []byte) []byte {
	b.head = new(FullBoxHeader)
	b.head.Load(bytes)
	b.data = bytes[b.head.DataSize():b.head.Length()]
	return bytes[b.head.Length():]
}

func (b *FullBox) Version() uint32 {
	var _b = []byte{0}
	_b = append(_b, b.head.Data()[:3]...)
	return binary.BigEndian.Uint32(_b)
}

func (b *FullBox) Flags() uint16 {
	var _b = []byte{0}
	_b = append(_b, b.head.Data()[3:4]...)
	return binary.BigEndian.Uint16(_b)
}
