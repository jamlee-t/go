// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by 'go generate' (with ./internal/asmgen). DO NOT EDIT.

//go:build !math_big_pure_go && (mips || mipsle)

#include "textflag.h"

// func addVV(z, x, y []Word) (c Word)
TEXT ·addVV(SB), NOSPLIT, $0
	MOVW z_len+4(FP), R1
	MOVW x_base+12(FP), R2
	MOVW y_base+24(FP), R3
	MOVW z_base+0(FP), R4
	// compute unrolled loop lengths
	AND $3, R1, R5
	SRL $2, R1
	XOR R26, R26	// clear carry
loop1:
	BEQ R5, loop1done
loop1cont:
	// unroll 1X
	MOVW 0(R2), R6
	MOVW 0(R3), R7
	ADDU R7, R6	// ADCS R7, R6, R6 (cr=R26)
	SGTU R7, R6, R23	// ...
	ADDU R26, R6	// ...
	SGTU R26, R6, R26	// ...
	ADDU R23, R26	// ...
	MOVW R6, 0(R4)
	ADDU $4, R2
	ADDU $4, R3
	ADDU $4, R4
	SUBU $1, R5
	BNE R5, loop1cont
loop1done:
loop4:
	BEQ R1, loop4done
loop4cont:
	// unroll 4X
	MOVW 0(R2), R5
	MOVW 4(R2), R6
	MOVW 8(R2), R7
	MOVW 12(R2), R8
	MOVW 0(R3), R9
	MOVW 4(R3), R10
	MOVW 8(R3), R11
	MOVW 12(R3), R12
	ADDU R9, R5	// ADCS R9, R5, R5 (cr=R26)
	SGTU R9, R5, R23	// ...
	ADDU R26, R5	// ...
	SGTU R26, R5, R26	// ...
	ADDU R23, R26	// ...
	ADDU R10, R6	// ADCS R10, R6, R6 (cr=R26)
	SGTU R10, R6, R23	// ...
	ADDU R26, R6	// ...
	SGTU R26, R6, R26	// ...
	ADDU R23, R26	// ...
	ADDU R11, R7	// ADCS R11, R7, R7 (cr=R26)
	SGTU R11, R7, R23	// ...
	ADDU R26, R7	// ...
	SGTU R26, R7, R26	// ...
	ADDU R23, R26	// ...
	ADDU R12, R8	// ADCS R12, R8, R8 (cr=R26)
	SGTU R12, R8, R23	// ...
	ADDU R26, R8	// ...
	SGTU R26, R8, R26	// ...
	ADDU R23, R26	// ...
	MOVW R5, 0(R4)
	MOVW R6, 4(R4)
	MOVW R7, 8(R4)
	MOVW R8, 12(R4)
	ADDU $16, R2
	ADDU $16, R3
	ADDU $16, R4
	SUBU $1, R1
	BNE R1, loop4cont
loop4done:
	MOVW R26, c+36(FP)
	RET

// func subVV(z, x, y []Word) (c Word)
TEXT ·subVV(SB), NOSPLIT, $0
	MOVW z_len+4(FP), R1
	MOVW x_base+12(FP), R2
	MOVW y_base+24(FP), R3
	MOVW z_base+0(FP), R4
	// compute unrolled loop lengths
	AND $3, R1, R5
	SRL $2, R1
	XOR R26, R26	// clear carry
loop1:
	BEQ R5, loop1done
loop1cont:
	// unroll 1X
	MOVW 0(R2), R6
	MOVW 0(R3), R7
	SGTU R26, R6, R23	// SBCS R7, R6, R6
	SUBU R26, R6	// ...
	SGTU R7, R6, R26	// ...
	SUBU R7, R6	// ...
	ADDU R23, R26	// ...
	MOVW R6, 0(R4)
	ADDU $4, R2
	ADDU $4, R3
	ADDU $4, R4
	SUBU $1, R5
	BNE R5, loop1cont
loop1done:
loop4:
	BEQ R1, loop4done
loop4cont:
	// unroll 4X
	MOVW 0(R2), R5
	MOVW 4(R2), R6
	MOVW 8(R2), R7
	MOVW 12(R2), R8
	MOVW 0(R3), R9
	MOVW 4(R3), R10
	MOVW 8(R3), R11
	MOVW 12(R3), R12
	SGTU R26, R5, R23	// SBCS R9, R5, R5
	SUBU R26, R5	// ...
	SGTU R9, R5, R26	// ...
	SUBU R9, R5	// ...
	ADDU R23, R26	// ...
	SGTU R26, R6, R23	// SBCS R10, R6, R6
	SUBU R26, R6	// ...
	SGTU R10, R6, R26	// ...
	SUBU R10, R6	// ...
	ADDU R23, R26	// ...
	SGTU R26, R7, R23	// SBCS R11, R7, R7
	SUBU R26, R7	// ...
	SGTU R11, R7, R26	// ...
	SUBU R11, R7	// ...
	ADDU R23, R26	// ...
	SGTU R26, R8, R23	// SBCS R12, R8, R8
	SUBU R26, R8	// ...
	SGTU R12, R8, R26	// ...
	SUBU R12, R8	// ...
	ADDU R23, R26	// ...
	MOVW R5, 0(R4)
	MOVW R6, 4(R4)
	MOVW R7, 8(R4)
	MOVW R8, 12(R4)
	ADDU $16, R2
	ADDU $16, R3
	ADDU $16, R4
	SUBU $1, R1
	BNE R1, loop4cont
loop4done:
	MOVW R26, c+36(FP)
	RET

// func lshVU(z, x []Word, s uint) (c Word)
TEXT ·lshVU(SB), NOSPLIT, $0
	MOVW z_len+4(FP), R1
	BEQ R1, ret0
	MOVW s+24(FP), R2
	MOVW x_base+12(FP), R3
	MOVW z_base+0(FP), R4
	// run loop backward
	SLL $2, R1, R5
	ADDU R5, R3
	SLL $2, R1, R5
	ADDU R5, R4
	// shift first word into carry
	MOVW -4(R3), R5
	MOVW $32, R6
	SUBU R2, R6
	SRL R6, R5, R7
	SLL R2, R5
	MOVW R7, c+28(FP)
	// shift remaining words
	SUBU $1, R1
	// compute unrolled loop lengths
	AND $3, R1, R7
	SRL $2, R1
loop1:
	BEQ R7, loop1done
loop1cont:
	// unroll 1X
	MOVW -8(R3), R8
	SRL R6, R8, R9
	OR R5, R9
	SLL R2, R8, R5
	MOVW R9, -4(R4)
	ADDU $-4, R3
	ADDU $-4, R4
	SUBU $1, R7
	BNE R7, loop1cont
loop1done:
loop4:
	BEQ R1, loop4done
loop4cont:
	// unroll 4X
	MOVW -8(R3), R7
	MOVW -12(R3), R8
	MOVW -16(R3), R9
	MOVW -20(R3), R10
	SRL R6, R7, R11
	OR R5, R11
	SLL R2, R7, R5
	SRL R6, R8, R7
	OR R5, R7
	SLL R2, R8, R5
	SRL R6, R9, R8
	OR R5, R8
	SLL R2, R9, R5
	SRL R6, R10, R9
	OR R5, R9
	SLL R2, R10, R5
	MOVW R11, -4(R4)
	MOVW R7, -8(R4)
	MOVW R8, -12(R4)
	MOVW R9, -16(R4)
	ADDU $-16, R3
	ADDU $-16, R4
	SUBU $1, R1
	BNE R1, loop4cont
loop4done:
	// store final shifted bits
	MOVW R5, -4(R4)
	RET
ret0:
	MOVW R0, c+28(FP)
	RET

// func rshVU(z, x []Word, s uint) (c Word)
TEXT ·rshVU(SB), NOSPLIT, $0
	MOVW z_len+4(FP), R1
	BEQ R1, ret0
	MOVW s+24(FP), R2
	MOVW x_base+12(FP), R3
	MOVW z_base+0(FP), R4
	// shift first word into carry
	MOVW 0(R3), R5
	MOVW $32, R6
	SUBU R2, R6
	SLL R6, R5, R7
	SRL R2, R5
	MOVW R7, c+28(FP)
	// shift remaining words
	SUBU $1, R1
	// compute unrolled loop lengths
	AND $3, R1, R7
	SRL $2, R1
loop1:
	BEQ R7, loop1done
loop1cont:
	// unroll 1X
	MOVW 4(R3), R8
	SLL R6, R8, R9
	OR R5, R9
	SRL R2, R8, R5
	MOVW R9, 0(R4)
	ADDU $4, R3
	ADDU $4, R4
	SUBU $1, R7
	BNE R7, loop1cont
loop1done:
loop4:
	BEQ R1, loop4done
loop4cont:
	// unroll 4X
	MOVW 4(R3), R7
	MOVW 8(R3), R8
	MOVW 12(R3), R9
	MOVW 16(R3), R10
	SLL R6, R7, R11
	OR R5, R11
	SRL R2, R7, R5
	SLL R6, R8, R7
	OR R5, R7
	SRL R2, R8, R5
	SLL R6, R9, R8
	OR R5, R8
	SRL R2, R9, R5
	SLL R6, R10, R9
	OR R5, R9
	SRL R2, R10, R5
	MOVW R11, 0(R4)
	MOVW R7, 4(R4)
	MOVW R8, 8(R4)
	MOVW R9, 12(R4)
	ADDU $16, R3
	ADDU $16, R4
	SUBU $1, R1
	BNE R1, loop4cont
loop4done:
	// store final shifted bits
	MOVW R5, 0(R4)
	RET
ret0:
	MOVW R0, c+28(FP)
	RET

// func mulAddVWW(z, x []Word, m, a Word) (c Word)
TEXT ·mulAddVWW(SB), NOSPLIT, $0
	MOVW m+24(FP), R1
	MOVW a+28(FP), R2
	MOVW z_len+4(FP), R3
	MOVW x_base+12(FP), R4
	MOVW z_base+0(FP), R5
	// compute unrolled loop lengths
	AND $3, R3, R6
	SRL $2, R3
loop1:
	BEQ R6, loop1done
loop1cont:
	// unroll 1X
	MOVW 0(R4), R7
	// synthetic carry, one column at a time
	MULU R1, R7
	MOVW LO, R8
	MOVW HI, R9
	ADDU R2, R8, R7	// ADDS R2, R8, R7 (cr=R26)
	SGTU R2, R7, R26	// ...
	ADDU R26, R9, R2	// ADC $0, R9, R2
	MOVW R7, 0(R5)
	ADDU $4, R4
	ADDU $4, R5
	SUBU $1, R6
	BNE R6, loop1cont
loop1done:
loop4:
	BEQ R3, loop4done
loop4cont:
	// unroll 4X
	MOVW 0(R4), R6
	MOVW 4(R4), R7
	MOVW 8(R4), R8
	MOVW 12(R4), R9
	// synthetic carry, one column at a time
	MULU R1, R6
	MOVW LO, R10
	MOVW HI, R11
	ADDU R2, R10, R6	// ADDS R2, R10, R6 (cr=R26)
	SGTU R2, R6, R26	// ...
	ADDU R26, R11, R2	// ADC $0, R11, R2
	MULU R1, R7
	MOVW LO, R10
	MOVW HI, R11
	ADDU R2, R10, R7	// ADDS R2, R10, R7 (cr=R26)
	SGTU R2, R7, R26	// ...
	ADDU R26, R11, R2	// ADC $0, R11, R2
	MULU R1, R8
	MOVW LO, R10
	MOVW HI, R11
	ADDU R2, R10, R8	// ADDS R2, R10, R8 (cr=R26)
	SGTU R2, R8, R26	// ...
	ADDU R26, R11, R2	// ADC $0, R11, R2
	MULU R1, R9
	MOVW LO, R10
	MOVW HI, R11
	ADDU R2, R10, R9	// ADDS R2, R10, R9 (cr=R26)
	SGTU R2, R9, R26	// ...
	ADDU R26, R11, R2	// ADC $0, R11, R2
	MOVW R6, 0(R5)
	MOVW R7, 4(R5)
	MOVW R8, 8(R5)
	MOVW R9, 12(R5)
	ADDU $16, R4
	ADDU $16, R5
	SUBU $1, R3
	BNE R3, loop4cont
loop4done:
	MOVW R2, c+32(FP)
	RET

// func addMulVVWW(z, x, y []Word, m, a Word) (c Word)
TEXT ·addMulVVWW(SB), NOSPLIT, $0
	MOVW m+36(FP), R1
	MOVW a+40(FP), R2
	MOVW z_len+4(FP), R3
	MOVW x_base+12(FP), R4
	MOVW y_base+24(FP), R5
	MOVW z_base+0(FP), R6
	// compute unrolled loop lengths
	AND $3, R3, R7
	SRL $2, R3
loop1:
	BEQ R7, loop1done
loop1cont:
	// unroll 1X
	MOVW 0(R4), R8
	MOVW 0(R5), R9
	// synthetic carry, one column at a time
	MULU R1, R9
	MOVW LO, R10
	MOVW HI, R11
	ADDU R8, R10	// ADDS R8, R10, R10 (cr=R26)
	SGTU R8, R10, R26	// ...
	ADDU R26, R11	// ADC $0, R11, R11
	ADDU R2, R10, R9	// ADDS R2, R10, R9 (cr=R26)
	SGTU R2, R9, R26	// ...
	ADDU R26, R11, R2	// ADC $0, R11, R2
	MOVW R9, 0(R6)
	ADDU $4, R4
	ADDU $4, R5
	ADDU $4, R6
	SUBU $1, R7
	BNE R7, loop1cont
loop1done:
loop4:
	BEQ R3, loop4done
loop4cont:
	// unroll 4X
	MOVW 0(R4), R7
	MOVW 4(R4), R8
	MOVW 8(R4), R9
	MOVW 12(R4), R10
	MOVW 0(R5), R11
	MOVW 4(R5), R12
	MOVW 8(R5), R13
	MOVW 12(R5), R14
	// synthetic carry, one column at a time
	MULU R1, R11
	MOVW LO, R15
	MOVW HI, R16
	ADDU R7, R15	// ADDS R7, R15, R15 (cr=R26)
	SGTU R7, R15, R26	// ...
	ADDU R26, R16	// ADC $0, R16, R16
	ADDU R2, R15, R11	// ADDS R2, R15, R11 (cr=R26)
	SGTU R2, R11, R26	// ...
	ADDU R26, R16, R2	// ADC $0, R16, R2
	MULU R1, R12
	MOVW LO, R15
	MOVW HI, R16
	ADDU R8, R15	// ADDS R8, R15, R15 (cr=R26)
	SGTU R8, R15, R26	// ...
	ADDU R26, R16	// ADC $0, R16, R16
	ADDU R2, R15, R12	// ADDS R2, R15, R12 (cr=R26)
	SGTU R2, R12, R26	// ...
	ADDU R26, R16, R2	// ADC $0, R16, R2
	MULU R1, R13
	MOVW LO, R15
	MOVW HI, R16
	ADDU R9, R15	// ADDS R9, R15, R15 (cr=R26)
	SGTU R9, R15, R26	// ...
	ADDU R26, R16	// ADC $0, R16, R16
	ADDU R2, R15, R13	// ADDS R2, R15, R13 (cr=R26)
	SGTU R2, R13, R26	// ...
	ADDU R26, R16, R2	// ADC $0, R16, R2
	MULU R1, R14
	MOVW LO, R15
	MOVW HI, R16
	ADDU R10, R15	// ADDS R10, R15, R15 (cr=R26)
	SGTU R10, R15, R26	// ...
	ADDU R26, R16	// ADC $0, R16, R16
	ADDU R2, R15, R14	// ADDS R2, R15, R14 (cr=R26)
	SGTU R2, R14, R26	// ...
	ADDU R26, R16, R2	// ADC $0, R16, R2
	MOVW R11, 0(R6)
	MOVW R12, 4(R6)
	MOVW R13, 8(R6)
	MOVW R14, 12(R6)
	ADDU $16, R4
	ADDU $16, R5
	ADDU $16, R6
	SUBU $1, R3
	BNE R3, loop4cont
loop4done:
	MOVW R2, c+44(FP)
	RET
