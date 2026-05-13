TEXT ·callJIT(SB), 0, $0-8
    PCALIGN $16
    MOVD code+0(FP), R0
    JMP (R0)
