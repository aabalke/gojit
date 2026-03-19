package gojit

import (
    "runtime"
    "testing"
    "unsafe"
)

// note: variables called within go funcs have to be global

func testExit(asm *Assembler) {

    ExitAssembler(asm)

    if asm.err != nil {
        panic(asm.err)
    }

	callJIT(uintptr(unsafe.Pointer(&asm.Buf[0])))

	asm.Release()
}

var called = false

func TestCall(t *testing.T) {

    pagesize := 512


	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(func() {
        called = true
	})

    testExit(asm)

    if !called {
        t.Errorf("Failed Test Call: called variable not set\n")
    }
}

var i = 1 << 16

//go:nosplit
func recursive() {
    if i > 0 {
        i--
        recursive()
    }
}

func TestCallRecursion(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(recursive)

    testExit(asm)
}

func TestCallGc(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(func ()  {
        runtime.GC()
	})

    testExit(asm)
}

var v uint64

func TestIndirect(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(func ()  {
        v = 0xBEEF
	})

	asm.Mov(Imm(0xDEAD), Rax)
	asm.MovAbs(uint64(uintptr(unsafe.Pointer(&v))), Rbx)
	asm.Mov(Rax, Indirect{Rbx, 0, 64})

	asm.InternalCallFunc(func ()  {
        v = 0xABBA
	})

    testExit(asm)

    if v != 0xABBA {
        t.Errorf("Failed Test Call: v variable not set to 0xABBA\n")
    }
}


func TestCallArguments(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

    var o uint64

    asm.Mov(Imm(1), Rax)
    asm.Mov(Imm(1), Rbx)
    asm.Mov(Imm(1), Rcx)
    asm.Mov(Imm(1), Rdi)
    asm.Mov(Imm(1), Rsi)
    asm.Mov(Imm(1), R8)
    asm.Mov(Imm(1), R9)

	asm.InternalCallFunc(func(a, b, c, d, e, f, g uint64) uint64  {
        return a+b+c+d+e+f+g
	})

    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&o))), Rbx)
    asm.Mov(Rax, Indirect{Rbx, 0, 64})

    testExit(asm)

    if o != 7 {
        t.Errorf("Failed Test Call Arguments: o variable not set to 0x7. Value: %X\n", o)
    }
}

func TestCallResults(t *testing.T) {

    pagesize := 512

	asm, err := New(pagesize)
	if err != nil {
		panic(err)
	}

	asm.InternalCallFunc(func() (
        uint64, uint64, uint64, uint64,
        uint64, uint64, uint64, uint64)  {
        return 1, 2, 3, 4, 5, 6, 7, 8
	})

    var ta, tb, tc, td, te, tf, tg, th uint64
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&ta))), R11)
    asm.Mov(Rax, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tb))), R11)
    asm.Mov(Rbx, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tc))), R11)
    asm.Mov(Rcx, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&td))), R11)
    asm.Mov(Rdi, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&te))), R11)
    asm.Mov(Rsi, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tf))), R11)
    asm.Mov(R8, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&tg))), R11)
    asm.Mov(R9, Indirect{R11, 0, 64})
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(&th))), R11)
    asm.Mov(R10, Indirect{R11, 0, 64})

    testExit(asm)

    if (
        ta != 1 ||
        tb != 2 ||
        tc != 3 ||
        td != 4 ||
        te != 5 ||
        tf != 6 ||
        tg != 7 ||
        th != 8) {
        t.Error("BAD")
    }
}

//func TestStress(t *testing.T) {
//
//    // this initially would be fine to run 1 << 6 times but not 1 << 7. 
//    // not sure why
//
//    println("starting stress test")
//
//	asm, err := New(1024 * 0x30)
//	if err != nil {
//		panic(err)
//	}
//
//    f := func ()  {
//        runtime.GC()
//    }
//
//    for range 1024 {
//        asm.InternalCallFunc(f)
//    }
//
//    testExit(asm)
//}


var (
    CpuPointer *Cpu
    CPU = R9
)

type Cpu struct {
    mem *Mem
}

type Mem struct {
}

func (m *Mem) Read(addr uint32, arm9 bool) uint8 {
    
	if arm9 {


		switch addr >> 24 {
        case 0x8, 0x9, 0xA:
            return 0
            return m.ReadGbaSlot(addr, arm9)
		}
	}

    panic("BAD")
}

func (m *Mem) ReadGbaSlot(addr uint32, arm9 bool) uint8 {

    if arm9 {
        return 0
    }

    panic("BAD")
}


func (m *Mem) Read8(addr uint32, arm9 bool) uint32 {
    return uint32(m.Read(addr, arm9))
}


func CallFunc(asm *Assembler, f any) {
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(CpuPointer))), CPU)
    asm.InternalCallFunc(f)
    asm.MovAbs(uint64(uintptr(unsafe.Pointer(CpuPointer))), CPU)
}

//go:nosplit
func Read(addr uint32) uint32 {
    return CpuPointer.mem.Read8(addr, true)
}

func TestSpecific(t *testing.T) {

    cpu := &Cpu{
        mem: &Mem{},
    }

    CpuPointer = cpu

    println("starting stress test specific")

    for range 0x100_0000 {

        asm, err := New(1024)
        if err != nil {
            panic(err)
        }

        asm.Mov(Imm(0x800_0000), Rax)
        CallFunc(asm, Read)

        testExit(asm)
    }


}
