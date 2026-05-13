// Package amd64 implements a simple amd64 assembler.
package arm64

import (
	"encoding/binary"
	"errors"
)

var ErrBufferTooSmall = errors.New("buffer is too small")

// Assembler implements a simple amd64 assembler. All methods on
// Assembler will emit code to Buf[Off:] and advances Off. Buf will
// never be reallocated, and attempts to assemble off the end of Buf
// will panic.
type Assembler struct {
	Buf []byte
	Off int

	err error
}

func New(size int) (*Assembler, error) {
	buf, e := Alloc(size)
	if e != nil || len(buf) == 0{
		return nil, e
	}

	return &Assembler{Buf: buf}, nil
}

func (a *Assembler) Release() {
	Release(a.Buf)
}

func (a *Assembler) Error() error {
	err := a.err
	a.err = nil
	return err
}

func (a *Assembler) addInst(inst uint32) {
	if a.Off + 3 > len(a.Buf) {
		a.err = ErrBufferTooSmall
		return
	}

	binary.LittleEndian.PutUint32(a.Buf[a.Off:], uint32(inst))
	a.Off += 4
}
