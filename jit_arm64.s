#include "funcdata.h"

TEXT ·callJIT(SB), 0, $144-16 // 72 but 16 aligned
    NO_LOCAL_POINTERS
    MOVD code+0(FP), R0
    JMP (R0)
gocall:
    PCALIGN $16
    WORD $0xDEADBE00

    MOVD R25, 16(RSP) // jit return addr
    MOVD R30, 24(RSP) // LR

    MOVD R8,  32(RSP)
    MOVD R9,  40(RSP)
    MOVD R10, 48(RSP)
    MOVD R11, 56(RSP)
    MOVD R12, 64(RSP)
    MOVD R13, 72(RSP)
    MOVD R14, 80(RSP)
    MOVD R15, 88(RSP)
    MOVD R16, 96(RSP)
    MOVD R17, 104(RSP)
    MOVD R19, 112(RSP)
    MOVD R20, 120(RSP)
    MOVD R21, 128(RSP)
    MOVD R22, 136(RSP)

    CALL R23

    MOVD 16(RSP), R25 // jit return addr
    MOVD 24(RSP), R30 // LR

    MOVD  32(RSP), R8
    MOVD  40(RSP), R9
    MOVD  48(RSP), R10
    MOVD  56(RSP), R11
    MOVD  64(RSP), R12
    MOVD  72(RSP), R13
    MOVD  80(RSP), R14
    MOVD  88(RSP), R15
    MOVD  96(RSP), R16
    MOVD 104(RSP), R17
    MOVD 112(RSP), R19
    MOVD 120(RSP), R20
    MOVD 128(RSP), R21
    MOVD 136(RSP), R22

    JMP (R25)
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
