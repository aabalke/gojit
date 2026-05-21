package gojit

const (
	PREFIX_LOCK     = 0xF0
	PREFIX_REPNZ    = 0xF2
	PREFIX_REPZ     = 0xF3
	PREFIX_SEG_CS   = 0x2E
	PREFIX_SEG_SS   = 0x36
	PREFIX_SEG_DS   = 0x3E
	PREFIX_SEG_ES   = 0x26
	PREFIX_SEG_FS   = 0x64
	PREFIX_SEG_GS   = 0x65
	PREFIX_OPSIZE   = 0x66
	PREFIX_ADDRSIZE = 0x67

	MOD_INDIR        = 0x0
	MOD_INDIR_DISP8  = 0x1
	MOD_INDIR_DISP32 = 0x2
	MOD_REG          = 0x3

	SCALE_1 = 0x0
	SCALE_2 = 0x1
	SCALE_4 = 0x2
	SCALE_8 = 0x3

	/* overflow */
	CC_O  = 0x0
	CC_NO = 0x1
	/* unsigned comparisons */
	CC_B  = 0x2
	CC_AE = 0x3
	CC_BE = 0x6
	CC_A  = 0x7
	/* zero */
	CC_Z  = 0x4
	CC_NZ = 0x5
	/* sign */
	CC_S  = 0x8
	CC_NS = 0x9
	/* parity */
	CC_P  = 0xA
	CC_NP = 0xB
	/* unsigned comparisons */
	CC_L  = 0xC
	CC_GE = 0xD
	CC_LE = 0xE
	CC_G  = 0xF

	/* alternative mnemonics */
	CC_C  = CC_B
	CC_NC = CC_AE

	PFX_REX = 0x40
	REXW    = 0x08
	REXR    = 0x04
	REXX    = 0x02
	REXB    = 0x01

	REG_DISP32 = 5
	REG_SIB    = 4
)

type Operand interface {
	// isOperand is unexported prevents external packages from
	// implementing Operand.
	isOperand()

	Rex(asm *Assembler, reg Register)
	ModRM(asm *Assembler, reg Register)
}

type Imm int32

func U32(u uint32) int32 {
	return int32(u)
}

func (i Imm) isOperand() {}
func (i Imm) Rex(asm *Assembler, reg Register) {
	panic("Imm.Rex")
}
func (i Imm) ModRM(asm *Assembler, reg Register) {
	panic("Imm.ModRM")
}

type Register struct {
	Val  byte
	Bits byte
}

func (r Register) isOperand() {}
func (i Register) Rex(asm *Assembler, reg Register) {
	asm.rexBits(i.Bits, reg.Bits, reg.Val > 7, false, i.Val > 7)
}

func (r Register) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_REG, reg.Val&7, r.Val&7)
}

var (
	Al  = Register{0, 8}
	Ax  = Register{0, 16}
	Eax = Register{0, 32}
	Rax = Register{0, 64}
	Cl  = Register{1, 8}
	Cx  = Register{1, 16}
	Ecx = Register{1, 32}
	Rcx = Register{1, 64}
	Dl  = Register{2, 8}
	Dx  = Register{2, 16}
	Edx = Register{2, 32}
	Rdx = Register{2, 64}
	Bl  = Register{3, 8}
	Bx  = Register{3, 16}
	Ebx = Register{3, 32}
	Rbx = Register{3, 64}
	Esp = Register{4, 32}
	Rsp = Register{4, 64}
	Ebp = Register{5, 32}
	Rbp = Register{5, 64}
	Esi = Register{6, 32}
	Rsi = Register{6, 64}
    Di  = Register{7, 8}
	Edi = Register{7, 32}
	Rdi = Register{7, 64}

	R8d  = Register{8, 32}
	R8   = Register{8, 64}
	R9d  = Register{9, 32}
	R9   = Register{9, 64}
	R10d = Register{10, 32}
	R10  = Register{10, 64}
	R11d = Register{11, 32}
	R11  = Register{11, 64}
	R12d = Register{12, 32}
	R12  = Register{12, 64}
	R13d = Register{13, 32}
	R13  = Register{13, 64}
	R14d = Register{14, 32}
	R14  = Register{14, 64}
	R15d = Register{15, 32}
	R15  = Register{15, 64}
)

type Indirect struct {
	Base   Register
	Offset int32
	Bits   byte
}

func (i Indirect) short() bool {
	return int32(int8(i.Offset)) == i.Offset
}

func (i Indirect) isOperand() {}
func (i Indirect) Rex(asm *Assembler, reg Register) {
	asm.rexBits(reg.Bits, i.Bits, reg.Val > 7, false, i.Base.Val > 7)
}

func (i Indirect) ModRM(asm *Assembler, reg Register) {
	if i.Base.Val == REG_SIB {
		SIB{i.Offset, Esp, Esp, Scale1}.ModRM(asm, reg)
		return
	}
	if i.Offset == 0 {
		asm.modrm(MOD_INDIR, reg.Val&7, i.Base.Val&7)
	} else if i.short() {
		asm.modrm(MOD_INDIR_DISP8, reg.Val&7, i.Base.Val&7)
		asm.byte(byte(i.Offset))
	} else {
		asm.modrm(MOD_INDIR_DISP32, reg.Val&7, i.Base.Val&7)
		asm.int32(uint32(i.Offset))
	}
}

type PCRel struct {
	Addr uintptr
}

func (i PCRel) isOperand() {}
func (i PCRel) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, false, false)
}
func (i PCRel) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_INDIR, reg.Val&7, REG_DISP32)
	asm.rel32(i.Addr)
}

type Scale struct {
	scale byte
}

var (
	Scale1 = Scale{SCALE_1}
	Scale2 = Scale{SCALE_2}
	Scale4 = Scale{SCALE_4}
	Scale8 = Scale{SCALE_8}
)

type SIB struct {
	Offset      int32
	Base, Index Register
	Scale       Scale
}

func (s SIB) isOperand() {}
func (s SIB) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, s.Index.Val > 7, s.Base.Val > 7)
}

func (s SIB) ModRM(asm *Assembler, reg Register) {
	if s.Offset != 0 {
		if int32(int8(s.Offset)) == s.Offset {
			asm.modrm(MOD_INDIR_DISP8, reg.Val&7, REG_SIB)
			asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
			asm.byte(uint8(s.Offset))
		} else {
			asm.modrm(MOD_INDIR_DISP32, reg.Val&7, REG_SIB)
			asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
			asm.int32(uint32(s.Offset))
		}
	} else {
		asm.modrm(MOD_INDIR, reg.Val&7, REG_SIB)
		asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
	}
}
