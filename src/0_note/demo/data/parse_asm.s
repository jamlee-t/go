// func checking()(before uintptr, after uintptr)
TEXT ·checking(SB),$4112-16
        LEAQ x-0(SP), DI         //
        MOVQ DI, before+0(FP)    // 将原伪寄存器SP偏移量存入返回值before

        MOVQ    SP, BP           // 存储物理SP偏移量到BP寄存器
        ADDQ    $4096, SP        // 将物理SP偏移增加4K

        LEAQ x-0(SP), SI
        MOVQ    BP, SP           // 恢复物理SP，因为修改物理SP后，伪寄存器FP/SP随之改变，
                                 // 为了正确访问FP，先恢复物理SP
        MOVQ SI, cpu+8(FP)       // 将偏移后的伪寄存器SP偏移量存入返回值after
        RET
                                 // 从输出的after-before来看，正好相差4K

GLOBL ·Id(SB),$8
DATA ·Id+0(SB)/1,$0x37
