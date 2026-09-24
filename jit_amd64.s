#include "funcdata.h"
#include "textflag.h"


TEXT ·CallJit(SB), 0, $48-8
    NO_LOCAL_POINTERS
    // assembler adds
    // add stack: return PC 
    // add stack: RBP on entry (PUSHQ BP)
    // locals 
    // call includes PUSHQ BP, MOVQ SP, BP, SUBQ (framesize), SP
    // (framesize = locals + args)

    MOVQ code+0(FP), AX
    JMP AX

call:
    LONG $0x636E7566 // func

    MOVQ R8,  8(SP) 
    MOVQ R9,  16(SP)
    MOVQ R10, 24(SP)
    MOVQ R11, 32(SP)
    MOVQ SI,  40(SP)

    PCALIGN $8

    CALL R12

    MOVQ 8(SP),  R8
    MOVQ 16(SP), R9
    MOVQ 24(SP), R10
    MOVQ 32(SP), R11
    MOVQ 40(SP), SI
    JMP (SP)

exit:
    LONG $0x74697865 // exit
    // assembler adds
    //ADDQ framesize, SP
    //POPQ BP
    RET

TEXT ·GetCallJitPtr(SB), 0, $0-8
    NO_LOCAL_POINTERS
    MOVQ $·CallJit(SB), AX  // address of ABI0 impl, not trampoline
    MOVQ AX, ret+0(FP)
    RET
