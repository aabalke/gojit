#include "funcdata.h"

TEXT ·callJIT(SB), 0, $80-16 // 72 but 16 aligned
    NO_LOCAL_POINTERS
    MOVD code+0(FP), R0
    JMP (R0)
gocall:
    PCALIGN $16
    WORD $0xDEADBE00
    MOVD R15, 16(RSP) // jit return addr
    MOVD R30, 24(RSP) // LR
    MOVD R8,  32(RSP)
    MOVD R9,  40(RSP)
    MOVD R10, 48(RSP)
    MOVD R11, 56(RSP)
    MOVD R12, 64(RSP)
    CALL R13
    MOVD 16(RSP), R15 // jit return addr
    MOVD 24(RSP), R30 // LR
    MOVD 32(RSP), R8
    MOVD 40(RSP), R9
    MOVD 48(RSP), R10
    MOVD 56(RSP), R11
    MOVD 64(RSP), R12
    JMP (R15)
//cleanup:
//    PCALIGN $16
//    WORD $0xDEADBE01
//    ADD $32+16, RSP, RSP
//    RET

TEXT ·callJITImplAddr(SB), 0, $0-16
    NO_LOCAL_POINTERS
    MOVD $·callJIT(SB), R0  // address of ABI0 impl, not trampoline
    MOVD R0, ret+0(FP)
    RET
