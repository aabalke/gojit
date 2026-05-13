package arm64




func (a *Assembler) AndReg(rd, rn, rm Reg, imm, shift uint32, s, sf bool) {
	if s {
		a._LogicalReg(rd, rn, rm, imm, 0b00, shift, false, sf)
		return
	}
	a._LogicalReg(rd, rn, rm, imm, 0b11, shift, false, sf)
}

func (a *Assembler) BicReg(rd, rn, rm Reg, imm, shift uint32, s, sf bool) {
	if s {
		a._LogicalReg(rd, rn, rm, imm, 0b00, shift, true, sf)
		return
	}
	a._LogicalReg(rd, rn, rm, imm, 0b11, shift, true, sf)
}

func (a *Assembler) OrrReg(rd, rn, rm Reg, imm, shift uint32, sf bool) {
	a._LogicalReg(rd, rn, rm, imm, 0b01, shift, false, sf)
}

func (a *Assembler) OrnReg(rd, rn, rm Reg, imm, shift uint32, sf bool) {
	a._LogicalReg(rd, rn, rm, imm, 0b01, shift, true, sf)
}

func (a *Assembler) EorReg(rd, rn, rm Reg, imm, shift uint32, sf bool) {
	a._LogicalReg(rd, rn, rm, imm, 0b10, shift, false, sf)
}

func (a *Assembler) EonReg(rd, rn, rm Reg, imm, shift uint32, sf bool) {
	a._LogicalReg(rd, rn, rm, imm, 0b10, shift, true, sf)
}

func (a *Assembler) _LogicalReg(rd, rn, rm Reg, imm, opc, shift uint32, n, sf bool) {

	if shift >= 0b100 {
		panic("logical reg bad shift")
	}

	if imm >= 1 << 6 {
		panic("logical reg imm >= 1 << 6")
	}

	v := uint32(0b1010) << 24

	if sf {
		v |= 1 << 31
	}

	if n {
		v |= 1 << 21
	}

	v |= opc << 29
	v |= shift << 22
	v |= imm << 10

	v |= uint32(rm) << 16
	v |= uint32(rn) << 5
	v |= uint32(rd)

	a.addInst(v)
}

func (a *Assembler) ADDReg(rd, rn, rm Reg, imm, shift uint32, set, sh, sf bool) {
	a._ADDSUBReg(rd, rn, rm, imm, shift, set, false, sf)
}
func (a *Assembler) SUBReg(rd, rn, rm Reg, imm, shift uint32, set, sh, sf bool) {
	a._ADDSUBReg(rd, rn, rm, imm, shift, set, true, sf)
}

func (a *Assembler) _ADDSUBReg(rd, rn, rm Reg, imm, shift uint32, set, sub, sf bool) {

	if shift > 0b10 {
		panic("reg add/sub with invalid shift")
	}

	if !sf && imm >= 1 << 4 {
		panic("reg add/sub word with invalid imm")
	}

	v := uint32(0b01011) << 24

	if sf {
		v |= 1 << 31
	}

	if sub {
		v |= 1 << 30
	}

	if set {
		v |= 1 << 29
	}

	v |= shift << 22
	v |= imm << 10
	v |= uint32(rm) << 16
	v |= uint32(rn) << 5
	v |= uint32(rd)

	a.addInst(v)
}

func (a *Assembler) ADDImm(rd, rn Reg, imm uint32, set, sh, sf bool) {
	a._ADDSUBImm(rd, rn, imm, false, set, sh, sf)
}
func (a *Assembler) SUBImm(rd, rn Reg, imm uint32, set, sh, sf bool) {
	a._ADDSUBImm(rd, rn, imm, true, set, sh, sf)
}

func (a *Assembler) _ADDSUBImm(rd, rn Reg, imm uint32, sub, set, sh, sf bool) {

	v := uint32(0b100010) << 23

	if sf {
		v |= 1 << 31
	}

	if set {
		v |= 1 << 29
	}

	if sh {
		v |= 1 << 22
	}

	if sub {
		v |= 1 << 30
	}

	if imm >= 1 << 12 {
		panic("add imm >= 1 << 12")
	}

	v |= imm << 10
	v |= uint32(rn) << 5
	v |= uint32(rd)

	a.addInst(v)
}

func (a *Assembler) STRReg(rt, rn, rm Reg, option uint32, scale, sf bool){
	a._LDRSTRReg(rt, rn, rm, option, false, scale, sf)
}

func (a *Assembler) LDRReg(rt, rn, rm Reg, option uint32, scale, sf bool){
	a._LDRSTRReg(rt, rn, rm, option, true, scale, sf)
}

// rt, [rn + imm]
func (a *Assembler) _LDRSTRReg(rt, rn, rm Reg, option uint32, ldr, scale, sf bool) {

	if option & 0b10 != 0 || option >= 1<<4 {
		panic("ldr/str invalid option value")
	}

	v := uint32(0b10_1110) << 26

	if sf {
		v |= 1 << 30
	}

	if scale {
		v |= 1 << 12
	}

	v |= 1 << 21

	if ldr {
		v |= 1 << 22
	}

	v |= uint32(rt)
	v |= uint32(rn) << 5
	v |= uint32(rm) << 16
	v |= option << 13

	a.addInst(v)
}

func (a *Assembler) STRImm(rt, rn Reg, imm uint32, pre, unsigned, sf bool){
	a._LDRSTRImm(rt, rn, imm, false, pre, unsigned, sf)
}

func (a *Assembler) LDRImm(rt, rn Reg, imm uint32, pre, unsigned, sf bool){
	a._LDRSTRImm(rt, rn, imm, true, pre, unsigned, sf)
}

// rt, [rn + imm]
func (a *Assembler) _LDRSTRImm(rt, rn Reg, imm uint32, ldr, pre, unsigned, sf bool) {

	v := uint32(0b10_1110) << 26

	if sf {
		v |= 1 << 30
	}

	if unsigned {
		v |= 1 << 24

		if imm >= 1 << 12 {
			panic("unsigned ldr/str imm. imm >= 1 << 12")
		}

		v |= imm << 10
	} else {
		if imm >= 1 << 19 {
			panic("signed ldr/str imm. imm >= 1 << 9")
		}

		v |= imm << 12

		if pre {
			v |= 0b11 << 10
		} else {
			v |= 0b01 << 10
		}
	}

	if ldr {
		v |= 1 << 22
	}

	v |= uint32(rt)
	v |= uint32(rn) << 5

	a.addInst(v)
}

const (
	MOVN = 0b00
	MOVZ = 0b10
	MOVK = 0b11

	HW_00 = 0b00
	HW_16 = 0b01
	HW_32 = 0b10
	HW_48 = 0b11
)

func (a *Assembler) Movz(rd Reg, imm, hw uint32, sf bool) {
	a._mov(rd, imm, MOVZ, hw, sf)
}

func (a *Assembler) Movk(rd Reg, imm, hw uint32, sf bool) {
	a._mov(rd, imm, MOVK, hw, sf)
}

func (a *Assembler) Movn(rd Reg, imm, hw uint32, sf bool) {
	a._mov(rd, imm, MOVN, hw, sf)
}

func (a *Assembler) _mov(rd Reg, imm, opc, hw uint32, sf bool) {

	if !sf && hw > 1 {
		panic("imm mov word inst with hw > 1")
	}

	v := uint32(0b1_0010_1) << 23

	if sf {
		v |= 1 << 31
	}

	v |= uint32(rd)
	v |= imm << 5
	v |= hw  << 21
	v |= opc << 29

	a.addInst(v)
}

func (a *Assembler) Ret() {
	a.addInst(0xD65F03C0)
}
