#include "funcdata.h"
#include "textflag.h"

TEXT ·callJIT(SB), 0, $8-8
    NO_LOCAL_POINTERS
    MOVQ code+0(FP), AX
    JMP AX
gocall:
    CALL R10
    JMP (SP)

TEXT ·callJITImplAddr(SB), 0, $0-8
    NO_LOCAL_POINTERS
    MOVQ $·callJIT(SB), AX  // address of ABI0 impl, not trampoline
    MOVQ AX, ret+0(FP)
    RET
