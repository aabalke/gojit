package gojit

func callJIT(code uintptr)     {}
func callJITImplAddr() uintptr { panic("jit called from wasm") }
