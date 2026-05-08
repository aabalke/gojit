//go:build !amd64

package gojit

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

func callJIT(code uintptr) {
    panic("calling jit from unsupported platform")
}
func callJITImplAddr() uintptr {
    panic("calling jit from unsupported platform")
}
