package tools

import (
	"errors"
	"io"
)

var (
	ErrFuncNil               = errors.New("function nil")
	ErrArrayIndexOutOfBounds = errors.New("bytes array bounds read")
)

type Reader struct {
	bytes        []byte
	length       uint64
	addBytesOver bool
	addBytesFunc func() ([]byte, error)
}

func (r *Reader) Read(n uint64) ([]byte, error) {
	if r.addBytesFunc == nil {
		return nil, ErrFuncNil
	}
	for r.length < n {
		if r.addBytesOver {
			return nil, ErrArrayIndexOutOfBounds
		}
		if err := r.CallAddBytesFunc(); err != nil {
			return nil, err
		}
	}
	_b := r.bytes[:n]
	r.bytes = r.bytes[n:]
	r.length = uint64(len(r.bytes))
	return _b, nil
}

func (r *Reader) Over() bool {
	return r.addBytesOver
}

func (r *Reader) CallAddBytesFunc() error {
	if b, err := r.addBytesFunc(); err != nil {
		if err == io.EOF {
			r.addBytesOver = true
		} else {
			return err
		}
	} else {
		r.bytes = append(r.bytes, b...)
		r.length = uint64(len(r.bytes))
	}
	return nil
}

func NewReader(addBytesFunc func() ([]byte, error)) *Reader {
	r := new(Reader)
	r.addBytesFunc = addBytesFunc
	return r
}
