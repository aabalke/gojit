package arm64

import (
	"bytes"
	"unsafe"
)

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

func callJIT(code uintptr)
func callJITImplAddr() uintptr

func getTaggedLabelAddr(tagIdx uint8) uintptr {
	impl := callJITImplAddr()
	bts := unsafe.Slice((*uint8)(unsafe.Pointer(impl)), 0x100)
	tagBytes := []uint8{tagIdx, 0xBE, 0xAD, 0xDE}
	offset := bytes.Index(bts, tagBytes)
	offset += 4 // past offset
	return impl + uintptr(offset)
}
