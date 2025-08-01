// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Export guts for testing.

package runtime

import (
	"internal/runtime/sys"
	"unsafe"
)

var (
	OsYield                 = osyield
	TimeBeginPeriodRetValue = &timeBeginPeriodRetValue
)

func NumberOfProcessors() int32 {
	var info systeminfo
	stdcall(_GetSystemInfo, uintptr(unsafe.Pointer(&info)))
	return int32(info.dwnumberofprocessors)
}

type ContextStub struct {
	context
}

func (c ContextStub) GetPC() uintptr {
	return c.ip()
}

func NewContextStub() *ContextStub {
	var ctx context
	ctx.set_ip(sys.GetCallerPC())
	ctx.set_sp(sys.GetCallerSP())
	ctx.set_fp(getcallerfp())
	return &ContextStub{ctx}
}
