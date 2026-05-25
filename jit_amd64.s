#include "funcdata.h"
#include "textflag.h"

TEXT ·callJIT(SB), 0, $32-8
    NO_LOCAL_POINTERS
    MOVQ code+0(FP), AX
    JMP AX
gocall:
    LONG $0xDEADBE00

    MOVQ R8, 8(SP)
    MOVQ R9, 16(SP)
    MOVQ SI, 24(SP)

    CALL R10

    MOVQ 8(SP), R8
    MOVQ 16(SP), R9
    MOVQ 24(SP), SI
    JMP (SP)

TEXT ·callJITImplAddr(SB), 0, $0-8
    NO_LOCAL_POINTERS
    MOVQ $·callJIT(SB), AX  // address of ABI0 impl, not trampoline
    MOVQ AX, ret+0(FP)
    RET
