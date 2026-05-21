package gojit

type Reg uint32

const (
	R00 Reg = iota
	R01
	R02
	R03
	R04
	R05
	R06
	R07
	R08
	R09
	R10
	R11
	R12
	R13
	R14
	R15
	R16
	R17
	R18
	R19
	R20
	R21
	R22
	R23
	R24
	R25
	R26
	R27
	R28
	R29
	R30

	RSP Reg = 0x1F
	RZR Reg = 0x1F // src zeros, dst discard
)

type Cond uint32

const (
	EQ Cond = iota
	NE
	CS
	CC
	MI
	PL
	VS
	VC
	HI
	LS
	GE
	LT
	GT
	LE
	//AL
	//NV
)

const (
	// psuedo names
	Z Cond = iota
	NZ
	C
	NC
	N
	NN
	V
	NV
)
