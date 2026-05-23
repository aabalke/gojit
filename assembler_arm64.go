package gojit

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/aabalke/gojit/cache"
)

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

func callJIT(code uintptr)
func callJITImplAddr() uintptr


func (a *Assembler) addInst(inst uint32) {
	if a.Off + 3 > len(a.Buf) {
		a.err = ErrBufferTooSmall
		return
	}

	binary.LittleEndian.PutUint32(a.Buf[a.Off:], uint32(inst))
	a.Off += 4
}

func getTaggedLabelAddr(tagIdx uint8) uintptr {
	impl := callJITImplAddr()
	bts := unsafe.Slice((*uint8)(unsafe.Pointer(impl)), 0x100)
	tagBytes := []uint8{tagIdx, 0xBE, 0xAD, 0xDE}
	offset := bytes.Index(bts, tagBytes)
	offset += 4 // past offset
	return impl + uintptr(offset)
}

func (a *Assembler) Mov64(rd Reg, v uint64) {
	a.Movz(rd, uint32(v>> 0) & 0xFFFF, HW_00, true)
	a.Movk(rd, uint32(v>>16) & 0xFFFF, HW_16, true)
	a.Movk(rd, uint32(v>>32) & 0xFFFF, HW_32, true)
	a.Movk(rd, uint32(v>>48) & 0xFFFF, HW_48, true)
}
func (a *Assembler) Mov32(rd Reg, v uint32) {
	a.Movz(rd, (v>> 0) & 0xFFFF, HW_00, false)
	a.Movk(rd, (v>>16) & 0xFFFF, HW_16, false)
}

func (a *Assembler) Exit() {
	
	// this amount needs to match the amount in callJIT asm text header
	a.ADDImm(RSP, RSP, (80 + 16), false, false, true)
	a.Ret()

	if err := a.Error(); err != nil {
		panic(err)
	}

	cache.ClearICache(a.Buf)
}
