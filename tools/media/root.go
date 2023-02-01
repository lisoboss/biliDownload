package media

import (
	"fmt"
	"os"
)

type Root struct {
	boxList []*Box
}

func (r *Root) AddBox(box *Box) {
	r.boxList = append(r.boxList, box)
}

func (r *Root) BoxList() []*Box {
	return r.boxList
}

func (r *Root) Index(i int) *Box {
	return r.boxList[i]
}

func (r *Root) Print() {
	for _, box := range r.boxList {
		box.Print("[+] ")
	}
}

func (r *Root) Load(bytes []byte) {
	for len(bytes) > 0 {
		box := new(Box)
		bytes = box.Load(bytes)
		r.boxList = append(r.boxList, box)
	}
}

func (r *Root) Dump() (bytes []byte) {
	for _, box := range r.boxList {
		bytes = append(bytes, box.Dump()...)
	}
	return
}

func (r *Root) GetBox(type_ string) *Box {
	for _, box := range r.boxList {
		box = box.GetBox(type_)
		if box != nil {
			return box
		}
	}
	return nil
}

func (r *Root) GetBoxByReverse(type_ string) *Box {
	n := len(r.boxList)
	for i := n - 1; i > -1; i-- {
		box := r.boxList[i]
		box = box.GetBox(type_)
		if box != nil {
			return box
		}
	}
	return nil
}

func (r *Root) GetBoxList(type_ string) (rets []*Box) {
	for _, box := range r.boxList {
		box = box.GetBox(type_)
		if box != nil {
			rets = append(rets, box)
		}
	}
	return
}

func NewRoot(filename string) (r *Root, err error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return r, fmt.Errorf("RootBox ReadFile: %s", err)
	}
	r = new(Root)
	r.Load(bytes)
	return
}
