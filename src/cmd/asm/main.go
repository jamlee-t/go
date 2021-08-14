// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"cmd/asm/internal/arch"
	"cmd/asm/internal/asm"
	"cmd/asm/internal/flags"
	"cmd/asm/internal/lex"

	"cmd/internal/bio"
	"cmd/internal/obj"
	"cmd/internal/objabi"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("asm: ")

	GOARCH := objabi.GOARCH

	// JAMLEE: 返回 asm.internal.arch.Arch 结构体。architecture.LinkArch.Init 方法
	architecture := arch.Set(GOARCH) // JAMLEE: GOARCH=amd64
	if architecture == nil {
		log.Fatalf("unrecognized architecture %s", GOARCH)
	}

	flags.Parse()

	// JAMLEE: Link 定义在 obj.link.Link
	ctxt := obj.Linknew(architecture.LinkArch)
	if *flags.PrintOut {
		ctxt.Debugasm = 1
	}
	ctxt.Flag_dynlink = *flags.Dynlink
	ctxt.Flag_shared = *flags.Shared || *flags.Dynlink
	ctxt.Flag_go115newobj = *flags.Go115Newobj
	ctxt.IsAsm = true

	// JAMLEE: 处理intel幽灵漏洞
	switch *flags.Spectre {
	default:
		log.Printf("unknown setting -spectre=%s", *flags.Spectre)
		os.Exit(2)
	case "":
		// nothing
	case "index":
		// known to compiler; ignore here so people can use
		// the same list with -gcflags=-spectre=LIST and -asmflags=-spectrre=LIST
	case "all", "ret":
		ctxt.Retpoline = true
	}

	// JAMLEE: Bso 设置为 os 的 Stdout
	ctxt.Bso = bufio.NewWriter(os.Stdout)
	defer ctxt.Bso.Flush()

	// JAMLEE: 架构包中的全局变量初始化，opindex，ycover之类的变量。
	architecture.Init(ctxt)

	// JAMLEE: 这个 buf 应该是写入对象，到底谁在往 buf 中写入数据呢？
	// Create object file, write header.
	buf, err := bio.Create(*flags.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer buf.Close()

	if !*flags.SymABIs {
		fmt.Fprintf(buf, "go object %s %s %s\n", objabi.GOOS, objabi.GOARCH, objabi.Version)
		fmt.Fprintf(buf, "!\n")
	}

	var ok, diag bool
	var failedFile string

	// JAMLEE: flag.Args 代表每个 .s 文件
	for _, f := range flag.Args() {
		lexer := lex.NewLexer(f)
		parser := asm.NewParser(ctxt, architecture, lexer)
		ctxt.DiagFunc = func(format string, args ...interface{}) {
			diag = true
			log.Printf(format, args...)
		}
		// JAMLEE: 解析符号ABI
		if *flags.SymABIs {
			ok = parser.ParseSymABIs(buf)
		} else {
			pList := new(obj.Plist)
			pList.Firstpc, ok = parser.Parse()
			// reports errors to parser.Errorf
			if ok {
				obj.Flushplist(ctxt, pList, nil, *flags.Importpath)
			}
		}
		if !ok {
			failedFile = f
			break
		}
	}
	if ok && !*flags.SymABIs {
		ctxt.NumberSyms(true)
		// JAMLEE: 写入 go obj 文件（不是 ar那种，而是 go obj header + data 的那种）
		obj.WriteObjFile(ctxt, buf, "")
	}
	if !ok || diag {
		if failedFile != "" {
			log.Printf("assembly of %s failed", failedFile)
		} else {
			log.Print("assembly failed")
		}
		buf.Close()
		os.Remove(*flags.OutputFile)
		os.Exit(1)
	}
}
