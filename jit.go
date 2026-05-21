// Package gojit contains basic support for writing JITs in golang. It
// contains functions for allocating byte slices in executable memory,
// and converting between such slices and golang function types.
package gojit

import (
	"errors"
	"unsafe"
	"github.com/edsrzf/mmap-go"
)

// Alloc returns a byte slice of the specified length that is marked
// RWX -- i.e. the memory in it can be both written and executed. This
// is just a simple wrapper around syscall.Mmap.
//
// len most likely needs to be a multiple of PageSize.
func Alloc(len int) ([]byte, error) {
	b, err := mmap.MapRegion(nil, len, mmap.EXEC|mmap.RDWR, mmap.ANON, int64(0))
	return b, err
}

// Release frees a buffer allocated by Alloc
func Release(b []byte) error {
	m := mmap.MMap(b)
	return m.Unmap()
}

// Addr returns the address in memory of a byte slice, as a uintptr
func Addr(b []byte) uintptr {
    return uintptr(unsafe.Pointer(unsafe.SliceData(b)))
}

func CallJit(b uintptr) {
    callJIT(b)
}

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
