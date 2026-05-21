package arm64

import "github.com/aabalke/gojit/arm64/cache"

func (a *Assembler) Mov64(rd Reg, v uint64) {
	a.Movz(rd, uint32(v>> 0) & 0xFFFF, HW_00, true)
	a.Movk(rd, uint32(v>>16) & 0xFFFF, HW_16, true)
	a.Movk(rd, uint32(v>>32) & 0xFFFF, HW_32, true)
	a.Movk(rd, uint32(v>>48) & 0xFFFF, HW_48, true)
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
