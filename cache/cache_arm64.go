package cache

// arm requires explicit instruction cache clearing

//#include <stdint.h>
//
//static void clearcache(void* start, void* end) {
//    __builtin___clear_cache((char*)start, (char*)end);
//}
import "C"

import (
	"runtime"
	"unsafe"
)

func ClearICache(code []byte) {
	if len(code) == 0 {
		return
	}

	var pinner runtime.Pinner
	pinner.Pin(&code[0])
	defer pinner.Unpin()

	start := unsafe.Pointer(&code[0])
	end := unsafe.Add(start, len(code))
	C.clearcache(start, end)
}
