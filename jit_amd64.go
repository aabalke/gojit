package gojit

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

// asm stubs
func callJIT(code *byte)
func callJITImplAddr() uintptr
