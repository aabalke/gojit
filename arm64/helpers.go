package arm64

func (a *Assembler) Mov64(rd Reg, v uint64) {
	a.Movz(rd, uint32(v>> 0) & 0xFFFF, HW_00, true)
	a.Movk(rd, uint32(v>>16) & 0xFFFF, HW_16, true)
	a.Movk(rd, uint32(v>>32) & 0xFFFF, HW_32, true)
	a.Movk(rd, uint32(v>>48) & 0xFFFF, HW_48, true)
}
