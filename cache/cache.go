package cache

// arm requires explicit instruction cache clearing

//#include <stdint.h>
//
//static void clearcache(void* start, void* end) {
//    __builtin___clear_cache((char*)start, (char*)end);
//}
import "C"

import "unsafe"

func ClearICache(code []byte) {
    if len(code) == 0 {
        return
    }

	start := unsafe.Pointer(&code[0])

	end := unsafe.Add(start, len(code))
    C.clearcache(start, end)
}
