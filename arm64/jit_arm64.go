package arm64

// PageSize is the size of a memory page. The len argument to Alloc
// should be an integer multiple of the page size.
const PageSize = 4096

func callJIT(code uintptr)
