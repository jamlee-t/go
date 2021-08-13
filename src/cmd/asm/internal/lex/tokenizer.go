// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"io"
	"os"
	"strings"
	"text/scanner"
	"unicode"

	"cmd/asm/internal/flags"
	"cmd/internal/objabi"
	"cmd/internal/src"
)

// JAMLEE: 实现接口 TokenReader
// A Tokenizer is a simple wrapping of text/scanner.Scanner, configured
// for our purposes and made a TokenReader. It forms the lowest level,
// turning text from readers into tokens.
type Tokenizer struct {
	tok  ScanToken
	s    *scanner.Scanner
	base *src.PosBase
	line int
	file *os.File // If non-nil, file descriptor to close.
}

// JAMLEE: name 是文件名，r 和 file 都是文件的 fd。使用 os.Open 返回的
func NewTokenizer(name string, r io.Reader, file *os.File) *Tokenizer {
	var s scanner.Scanner
	// JAMLEE: Scanner 的用法是 Init 一个标准的 reader 进去
	s.Init(r)
	// JAMLEE: 设置空白字符定义
	// Newline is like a semicolon; other space characters are fine.
	s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	// Don't skip comments: we need to count newlines.
	s.Mode = scanner.ScanChars |
		scanner.ScanFloats |
		scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanStrings |
		scanner.ScanComments
	s.Position.Filename = name
	s.IsIdentRune = isIdentRune

	// JAMLEE: 返回一个 Tokenizer, 里面封装了 scanner
	return &Tokenizer{
		s:    &s,
		base: src.NewFileBase(name, objabi.AbsFile(objabi.WorkingDir(), name, *flags.TrimPath)),
		line: 1,
		file: file,
	}
}

// JAMLEE: 是否是 identifier character。定义给 scanner 用的
// We want center dot (·) and division slash (∕) to work as identifier characters.
func isIdentRune(ch rune, i int) bool {
	if unicode.IsLetter(ch) {
		return true
	}
	switch ch {
	case '_': // Underscore; traditional.
		return true
	case '\u00B7': // Represents the period in runtime.exit. U+00B7 '·' middle dot
		return true
	case '\u2215': // Represents the slash in runtime/debug.setGCPercent. U+2215 '∕' division slash
		return true
	}
	// Digits are OK only after the first character.
	return i > 0 && unicode.IsDigit(ch)
}

// JAMLEE: 扫描后返回当前扫描的字符串。类似 s.TokenText()
func (t *Tokenizer) Text() string {
	switch t.tok {
	case LSH:
		return "<<"
	case RSH:
		return ">>"
	case ARR:
		return "->"
	case ROT:
		return "@>"
	}
	return t.s.TokenText()
}

func (t *Tokenizer) File() string {
	return t.base.Filename()
}

func (t *Tokenizer) Base() *src.PosBase {
	return t.base
}

func (t *Tokenizer) SetBase(base *src.PosBase) {
	t.base = base
}

func (t *Tokenizer) Line() int {
	return t.line
}

func (t *Tokenizer) Col() int {
	return t.s.Pos().Column
}

// JAMLEE: 等于 Scanner.Scan, ScanToken 表示当前扫描返回的 token 类型
func (t *Tokenizer) Next() ScanToken {
	s := t.s
	// JAMLEE: 这里跳过注释内容
	for {
		t.tok = ScanToken(s.Scan())
		if t.tok != scanner.Comment {
			break
		}
		length := strings.Count(s.TokenText(), "\n")
		t.line += length
		// TODO: If we ever have //go: comments in assembly, will need to keep them here.
		// For now, just discard all comments.
	}
	switch t.tok {
	case '\n':
		t.line++
	case '-':
		if s.Peek() == '>' {
			s.Next()
			t.tok = ARR
			return ARR
		}
	case '@':
		if s.Peek() == '>' {
			s.Next()
			t.tok = ROT
			return ROT
		}
	case '<':
		if s.Peek() == '<' {
			s.Next()
			t.tok = LSH
			return LSH
		}
	case '>':
		if s.Peek() == '>' {
			s.Next()
			t.tok = RSH
			return RSH
		}
	}
	return t.tok
}

func (t *Tokenizer) Close() {
	if t.file != nil {
		t.file.Close()
	}
}
