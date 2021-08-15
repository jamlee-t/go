// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amd64

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj/x86"
)

var leaptr = x86.ALEAQ

// JAMLEE: gc.Arch 。类似的 asm 里面也有这个扩展自己的 Arch 结构。将机器代码转换的能力放到中。这里是 x86.Linkamd64 中关联的span6方法。
func Init(arch *gc.Arch) {
	arch.LinkArch = &x86.Linkamd64
	arch.REGSP = x86.REGSP
	arch.MAXWIDTH = 1 << 50

	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop
	arch.Ginsnopdefer = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
}
