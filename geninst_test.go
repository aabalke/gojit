package gojit

//import (
//	"testing"
//	"unsafe"
//)
//
//var (
//    in = [3]uint64{
//        0xFFFF_FFFF,
//        uint64(int8(100)),
//        uint64(0xFFFF),
//    }
//
//    out = [3]uint64{
//    }
//
//    corr = [3]uint64{
//        0xFFFF_FFFF_FFFF_FFFF,
//        uint64(int64(int8(100))),
//        uint64(0xFFFF_FFFF),
//    }
//)

//func TestIndirect(t *testing.T) {
//
//    for i := range int32(len(in)) {
//
//        pagesize := 128
//        asm, err := New(pagesize)
//        if err != nil {
//            panic(err)
//        }
//
//        asm.MovAbs(uint64(uintptr(unsafe.Pointer(&in))), Rbx)
//        asm.MovAbs(uint64(uintptr(unsafe.Pointer(&out))), Rcx)
//        asm.Mov(Indirect{Rbx, i * 8, 64}, Rax)
//
//        switch i {
//        case 0:
//            asm.Movsxd(Eax, Rax)
//        case 1:
//            asm.Movsx(Al, Rax)
//        case 2:
//            asm.Movsx(Ax, Eax)
//        }
//
//        asm.Mov(Rax, Indirect{Rcx, i * 8, 64})
//
//        testExit(asm)
//
//        if out[i] != corr[i] {
//            t.Errorf("Bad Movsxd Idx %04d CORR %16X RECI %16X\n", i, corr[i], out[i])
//        }
//    }
//}
//
