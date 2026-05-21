package arm64

import (
	"encoding/binary"
)

func (a *Assembler) Cinc(rd, rn Reg, cond Cond, sf bool) {

	if rn == 0x1F {
		panic("cinc with rn == 0x1F")
	}

	a.Csinc(rd, rn, rn, cond ^ 1, sf)
}

func (a *Assembler) Cset(rd Reg, cond Cond, sf bool) {
	a.Csinc(rd, 0x1F, 0x1F, cond ^ 1, sf)
}

func (a *Assembler) Csel(rd, rn, rm Reg, cond Cond, sf bool) {
	a._condSelect(rd, rn, rm, cond, false, false, sf)
}

func (a *Assembler) Csinc(rd, rn, rm Reg, cond Cond, sf bool) {
	a._condSelect(rd, rn, rm, cond, false, true, sf)
}

func (a *Assembler) Csinv(rd, rn, rm Reg, cond Cond, sf bool) {
	a._condSelect(rd, rn, rm, cond, true, false, sf)
}

func (a *Assembler) Csneg(rd, rn, rm Reg, cond Cond, sf bool) {
	a._condSelect(rd, rn, rm, cond, true, true, sf)
}

func (a *Assembler) _condSelect(rd, rn, rm Reg, cond Cond, op, op2, sf bool) {

	v := uint32(0b00011010100) << 21

	if op2 {
		v |= 1 << 10
	}

	if op {
		v |= 1 << 30
	}

	v |= uint32(rm)   << 16
	v |= uint32(cond) << 12
	v |= uint32(rn)   << 5
	v |= uint32(rd)

	if sf {
		v |= 1 << 31
	}

	a.addInst(v)
}

func (a *Assembler) Smaddl(rd, rm, rn, ra Reg) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b01, false, false, true)
}

func (a *Assembler) Smsubl(rd, rm, rn, ra Reg) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b01, false, true, true)
}

func (a *Assembler) Umaddl(rd, rm, rn, ra Reg) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b01, true, false, true)
}

func (a *Assembler) Umsubl(rd, rm, rn, ra Reg) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b01, true, true, true)
}

func (a *Assembler) Mul(rd, rm, rn Reg, sf bool) {
	a.MAdd(rd, rm, rn, RZR, sf)
}

func (a *Assembler) MAdd(rd, rm, rn, ra Reg, sf bool) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b00, false, false, sf)
}

func (a *Assembler) MSub(rd, rm, rn, ra Reg, sf bool) {
	a._dataProcessSrc3(rd, rm, rn, ra, 0b00, true, true, sf)
}

func (a *Assembler) _dataProcessSrc3(rd, rm, rn, ra Reg, op uint32, unsigned, sub, sf bool) {

	if op >= 1 << 2 {
		panic("data process src3 with op >= 1 << 2")
	}

	v := uint32(0b11011) << 24
	v |= uint32(rm) << 16
	v |= uint32(ra) << 10
	v |= uint32(rn) << 5
	v |= uint32(op) << 21
	v |= uint32(rd)

	if sf {
		v |= 1 << 31
	}

	if unsigned {
		v |= 1 << 23
	}

	if sub {
		v |= 1 << 15
	}

	a.addInst(v)
}

func (a *Assembler) Clz(rd, rn Reg, sf bool) {

	v := uint32(0b0101_1010_1100_0000_0001_00) << 10
	v |= uint32(rn) << 5
	v |= uint32(rd)

	if sf {
		v |= 1 << 31
	}

	a.addInst(v)
}

func (a *Assembler) TstReg(rn, rm Reg, imm, shift uint32, sf bool) {
	a.AndReg(RZR, rn, rm, imm, shift, true, sf)
}

func (a *Assembler) AndReg(rd, rn, rm Reg, imm, shift uint32, s, sf bool) {
	if s {
		a._LogicalReg(rd, rn, rm, imm, 0b11, shift, false, sf)
		return
	}
	a._LogicalReg(rd, rn, rm, imm, 0b00, shift, false, sf)
}

func (a *Assembler) BicReg(rd, rn, rm Reg, imm, shift uint32, s, sf bool) {
	if s {
		a._LogicalReg(rd, rn, rm, imm, 0b11, shift, true, sf)
		return
	}
	a._LogicalReg(rd, rn, rm, imm, 0b00, shift, true, sf)
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

const (
	SHIFT_LSL = iota
	SHIFT_LSR
	SHIFT_ASR
	SHIFT_ROR
)

func (a *Assembler) _LogicalReg(rd, rn, rm Reg, imm, opc, shift uint32, n, sf bool) {

	if shift >= 0b100 {
		panic("logical reg bad shift")
	}

	if imm >= 1 << 6 {
		panic("logical reg imm >= 1 << 6")
	}

	if opc >= 1 << 2 {
		panic("logical reg opc >= 1 << 2")
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

func (a *Assembler) TstImm(rd Reg, imm Immediate, sf bool) {
	a.AndImm(rd, rd, imm, true, sf)
}

func (a *Assembler) AndImm(rd, rn Reg, imm Immediate, set, sf bool) {
	if set {
		a._LogicalImm(rd, rn, imm, 0b11, sf)
	}

	a._LogicalImm(rd, rn, imm, 0b00, sf)
}

func (a *Assembler) OrrImm(rd, rn Reg, imm Immediate, sf bool) {
	a._LogicalImm(rd, rn, imm, 0b01, sf)
}

func (a *Assembler) EorImm(rd, rn Reg, imm Immediate, sf bool) {
	a._LogicalImm(rd, rn, imm, 0b10, sf)
}

func (a *Assembler) _LogicalImm(rd, rn Reg, imm Immediate, op uint32, sf bool) {

	if op >= 1 << 2 {
		panic("logical imm with op >= 1 << 2")
	}

	v := uint32(0b100100) << 23

	if sf {
		v |= 1 << 31
	}

	v |= op    << 29
	v |= imm.n << 22
	v |= imm.immr << 16
	v |= imm.imms << 10
	v |= uint32(rn) << 5
	v |= uint32(rd)

	a.addInst(v)
}

func (a *Assembler) NegReg(rd, rm Reg, imm, shift uint32, set, sh, sf bool) {
	a.SUBReg(rd, RZR, rm, imm, shift, set, sh, sf)
}

func (a *Assembler) CmnReg(rn, rm Reg, imm, shift uint32, sh, sf bool) {
	a.ADDReg(RZR, rn, rm, imm, shift, true, sh, sf)
}

func (a *Assembler) CmpReg(rn, rm Reg, imm, shift uint32, sh, sf bool) {
	a.SUBReg(RZR, rn, rm, imm, shift, true, sh, sf)
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

func (a *Assembler) CmnImm(rn Reg, imm, shift uint32, sh, sf bool) {
	a.ADDImm(RZR, rn, imm, true, sh, sf)
}

func (a *Assembler) CmpImm(rn Reg, imm, shift uint32, sh, sf bool) {
	a.SUBImm(RZR, rn, imm, true, sh, sf)
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

const (
	SIZE_BYTE = 0b00
	SIZE_HALF = 0b01
	SIZE_WORD = 0b10
	SIZE_DWRD = 0b11
)

func (a *Assembler) StrReg(rt, rn, rm Reg, option, size uint32, scale bool) {
	a._ldrStrReg(rt, rn, rm, option, size, false, scale)
}

func (a *Assembler) LdrReg(rt, rn, rm Reg, option, size uint32, scale bool) {
	a._ldrStrReg(rt, rn, rm, option, size, true, scale)
}

// rt, [rn + imm]
func (a *Assembler) _ldrStrReg(rt, rn, rm Reg, option, size uint32, ldr, scale bool) {

	if option & 0b10 != 0 || option >= 1<<4 {
		panic("ldr/str invalid option value")
	}

	v := uint32(0b1110) << 26

	v |= size << 30

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

func (a *Assembler) StrImm(rt, rn Reg, imm, size uint32, pre, unsigned bool) {
	a._ldrStrImm(rt, rn, imm, size, false, pre, unsigned)
}

func (a *Assembler) LdrImm(rt, rn Reg, imm, size uint32, pre, unsigned bool) {
	a._ldrStrImm(rt, rn, imm, size, true, pre, unsigned)
}

// rt, [rn + imm]
func (a *Assembler) _ldrStrImm(rt, rn Reg, imm, size uint32, ldr, pre, unsigned bool) {

	v := uint32(0b1110) << 26

	v |= size << 30

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

	if imm >= 1 << 16 {
		panic("imm mov imm >= 16")
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

func (a *Assembler) Sxtb(rd, rn Reg, sf bool) {

	i := Immediate{0, 0x7, 0x0, sf}

	if sf { i.n = 1 }

	a.Bfm(rd, rn, i, false, sf)
}

func (a *Assembler) Sxth(rd, rn Reg, sf bool) {

	i := Immediate{0, 0xF, 0x0, sf}

	if sf { i.n = 1 }

	a.Bfm(rd, rn, i, false, sf)
}

func (a *Assembler) Sxtw(rd, rn Reg) {
	a.Bfm(rd, rn, Immediate{1, 0x1F, 0x0, true}, false, true)
}

func (a *Assembler) AsrImm(rd, rn Reg, shift uint32, sf bool) {

	var width, N, imms uint32
	if sf {
		width = 64
		N = 1
		imms = 0x3F
	} else {
		width = 32
		N = 0
		imms = 0x1F
	}

    if shift >= width {
        panic("invalid ASR shift")
    }

	a.Bfm(rd, rn, Immediate{N, imms, shift, sf}, false, sf)
}

func (a *Assembler) LsrImm(rd, rn Reg, shift uint32, sf bool) {

	var width, N uint32
	if sf {
		width = 64
		N = 1
	} else {
		width = 32
		N = 0
	}

    if shift >= width {
        panic("invalid LSR shift")
    }

    immr := shift
    imms := width - 1
	a.Bfm(rd, rn, Immediate{N, imms, immr, sf}, true, sf)
}

func (a *Assembler) LslImm(rd, rn Reg, shift uint32, sf bool) {

	var width, N uint32
	if sf {
		width = 64
		N = 1
	} else {
		width = 32
	}

	immr := (width - shift) % width
	imms := (width - 1) - shift

	a.Bfm(rd, rn, Immediate{N, imms, immr, sf}, true, sf)
}

func (a *Assembler) Bfm(rd, rn Reg, imm Immediate, unsigned, sf bool) {

	v := uint32(0b100110) << 23
	v |= imm.n      << 22
	v |= imm.immr   << 16
	v |= imm.imms   << 10
	v |= uint32(rn) << 5
	v |= uint32(rd)

	if unsigned {
		v |= 1 << 30
	}

	if sf {
		v |= 1 << 31
	}

	a.addInst(v)
}

func (a *Assembler) Bfi(rd, rn Reg, width, lsb uint32, sf bool) {

	v := uint32(0b001100110) << 23

	size := uint32(32)
	if sf {
		v |= 1 << 31
		v |= 1 << 22
		size = 64
	}

	if rn >= 0x1F {
		panic("bfi invalid rn")
	}

	if width == 0 {
		panic("bfi width cannot be 0")
	}
	if lsb >= size {
		panic("bfi lsb out of range")
	}
	if width > size-lsb {
		panic("bfi width out of range")
	}

	immr := (size - lsb) & (size - 1)
    imms := width - 1

	v |= (immr & 0x3F) << 16
	v |= (imms & 0x3F) << 10
	v |= uint32(rn) << 5
	v |= uint32(rd)

	a.addInst(v)
}

func (a *Assembler) MovReg(rd, rm Reg, sf bool) {

	v := uint32(0b10101) << 25
	v |= 0x1F << 5
	v |= uint32(rd)
	v |= uint32(rm) << 16

	if sf {
		v |= 1 << 31
	}

	a.addInst(v)
}


func (a *Assembler) RorImm(rd, rs Reg, shift uint32, sf bool) {
	if shift >= 1 << 5 {
		panic("ror imm with shift >= 1 << 5")
	}
	a.ExtrReg(rd, rs, rs, shift, sf)
}

func (a *Assembler) ExtrReg(rd, rn, rm Reg, imm uint32, sf bool) {

	if imm >= 1 << 6 {
		panic("extr reg imm >= 1 << 6")
	}

	v := uint32(0b100111) << 23

	if sf {
		v |= 1 << 31
		v |= 1 << 22
	}

	v |= uint32(rd)
	v |= uint32(rn) << 5
	v |= uint32(rm) << 16
	v |= imm << 10

	a.addInst(v)
}



const (
	UDIV = 0b000010
	SDIV = 0b000011
	LSLV = 0b001000
	LSRV = 0b001001
	ASRV = 0b001010
	RORV = 0b001011
	SMAX = 0b011000
	UMAX = 0b011001
	SMIN = 0b011010
	UMIN = 0b011011
)

func (a *Assembler) UdivReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, UDIV, sf)
}

func (a *Assembler) SdivReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, SDIV, sf)
}

func (a *Assembler) LslReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, LSLV, sf)
}

func (a *Assembler) LsrReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, LSRV, sf)
}

func (a *Assembler) AsrReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, ASRV, sf)
}

func (a *Assembler) RorReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, RORV, sf)
}

func (a *Assembler) SmaxReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, SMAX, sf)
}

func (a *Assembler) UmaxReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, SMIN, sf)
}

func (a *Assembler) SminReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, UMAX, sf)
}

func (a *Assembler) UminReg(rd, rn, rm Reg, sf bool) {
	a._dataProcessSrc2(rd, rn, rm, UMIN, sf)
}

func (a *Assembler) _dataProcessSrc2(rd, rn, rm Reg, op uint32, sf bool) {

	v := uint32(0b1101011) << 22
	v |= op << 10
	v |= uint32(rd)
	v |= uint32(rn) << 5
	v |= uint32(rm) << 16

	if sf {
		v |= 1 << 31
	}

	a.addInst(v)
}

func (a *Assembler) AdcReg(rd, rn, rm Reg, set, sf bool) {
	a._ADDSUBCARRY(rd, rn, rm, false, set, sf)
}

func (a *Assembler) SbcReg(rd, rn, rm Reg, set, sf bool) {
	a._ADDSUBCARRY(rd, rn, rm, true, set, sf)
}

func (a *Assembler) _ADDSUBCARRY(rd, rn, rm Reg, sub, set, sf bool) {

	v := uint32(0b1101)  << 25

	if sf {
		v |= 1<<31
	}

	if sub {
		v |= 1<<30
	}

	if set {
		v |= 1<<29
	}

	v |= uint32(rd)
	v |= uint32(rn) << 5
	v |= uint32(rm) << 16

	a.addInst(v)
}

//// move sys reg => general purpose reg
func (a *Assembler) Mrs(rt Reg) {
	v := uint32(0xd51b4200)
	v |= 1 << 21
	v |= uint32(rt)
	a.addInst(v)
}

// move general purpose reg => sys reg
func (a *Assembler) Msr(rt Reg) {
	v := uint32(0xd51b4200)
	v |= uint32(rt)
	a.addInst(v)
}


func (a *Assembler) B() func() {

	branchOff := a.Off
	a.addInst(0xDEADBEEF) // temporary, will crash if func not called

	return func() {
		immBytes := (a.Off - branchOff) >> 2
		v := uint32(0b101) << 26 | uint32(immBytes)
		binary.LittleEndian.PutUint32(a.Buf[branchOff:], v)
	}
}

func (a *Assembler) BCond(cond Cond) func() {

	branchOff := a.Off
	a.addInst(0xDEADBEEF) // temporary, will crash if func not called

	return func() {
		immBytes := ((a.Off - branchOff) >> 2) << 5
		v := uint32(0b010101) << 26 | uint32(immBytes) | uint32(cond)
		binary.LittleEndian.PutUint32(a.Buf[branchOff:], v)
	}
}

func (a *Assembler) Custom(v uint32) {
	a.addInst(v)
}
