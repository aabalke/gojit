package gojit

// these functions are useful for auto setting up framesizes for call jit

func exit(asm *Assembler) {

    // framesize of 16 is used since "TEXT ·callJIT(SB), 0, $8-8"

	// addq fs, sp
	fs := byte(8 + 8)
	asm.byte(0x48)
	asm.byte(0x83)
	asm.byte(0xc4)
	asm.byte(fs)
	asm.byte(0x48)

    // ret
	asm.Ret()
}
