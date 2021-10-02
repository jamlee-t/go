// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements source, a buffered rune reader
// specialized for scanning Go code: Reading
// ASCII characters, maintaining current (line, col)
// position information, and recording of the most
// recently read source segment are highly optimized.
// This file is self-contained (go tool compile source.go
// compiles) and thus could be made into its own package.

package syntax

import (
	"io"
	"unicode/utf8"
)

// JAMLEE: source 其实是一个Reader，带有 buf。有三个 index 变量分别时 begin、read、end。
// The source buffer is accessed using three indices b (begin),
// r (read), and e (end):
//
// - If b >= 0, it points to the beginning of a segment of most
//   recently read characters (typically a Go literal).
//
// - r points to the byte immediately following the most recently
//   read character ch, which starts at r-chw.
//
// - e points to the byte immediately following the last byte that
//   was read into the buffer.
//
// The buffer content is terminated at buf[e] with the sentinel
// character utf8.RuneSelf. This makes it possible to test for
// the common case of ASCII characters with a single 'if' (see
// nextch method).
//
//                +------ content in use -------+
//                v                             v
// buf [...read...|...segment...|ch|...unread...|s|...free...]
//                ^             ^  ^            ^
//                |             |  |            |
//                b         r-chw  r            e
//
// Invariant: -1 <= b < r <= e < len(buf) && buf[e] == sentinel

type source struct {
	in   io.Reader
	errh func(line, col uint, msg string) // JAMLEE: 由外部传入一个错误处理函数

	buf       []byte // source buffer
	ioerr     error  // pending I/O error, or nil
	b, r, e   int    // buffer indices (see comment above)
	line, col uint   // source position of ch (0-based) // JAMLEE: 当前处理的列位置
	ch        rune   // most recently read character, JAMLEE: 最近读取的一个字符，rune代表一个utf-8字符
	chw       int    // width of ch ,JAMLEE: 字符宽度，utf8 最高可以 6 个 byte。但是 4 个 byte 足够表示。
}

// JAMLEE: ASCII 的字符最大值，用于判断是否是 ASCII 的字符
const sentinel = utf8.RuneSelf

// JAMLEE: source 初始化，
func (s *source) init(in io.Reader, errh func(line, col uint, msg string)) {
	s.in = in
	s.errh = errh

	if s.buf == nil {
		s.buf = make([]byte, nextSize(0))
	}
	s.buf[0] = sentinel
	s.ioerr = nil
	s.b, s.r, s.e = -1, 0, 0
	s.line, s.col = 0, 0
	s.ch = ' '
	s.chw = 0
}

// JAMLEE: 读取的内容起点是第一行第一列。默认的 base地址
// starting points for line and column numbers
const linebase = 1
const colbase = 1

// pos returns the (line, col) source position of s.ch.
func (s *source) pos() (line, col uint) {
	return linebase + s.line, colbase + s.col
}

// JAMLEE: 处理遇到错误时，用于将错误信息信息传递到传入的错误处理函数
// error reports the error msg at source position s.pos().
func (s *source) error(msg string) {
	line, col := s.pos()
	s.errh(line, col, msg)
}

// JAMLEE: 开始一个新的 「source segment」。
// start starts a new active source segment (including s.ch).
// As long as stop has not been called, the active segment's
// bytes (excluding s.ch) may be retrieved by calling segment.
func (s *source) start()          { s.b = s.r - s.chw } // JAMLEE: 根据字符长度找到字符的 begin
func (s *source) stop()           { s.b = -1 }
func (s *source) segment() []byte { return s.buf[s.b : s.r-s.chw] } // JAMLEE: 已经分析的部分，但是暂时还没有得出token结果

// JAMLEE: 回道 segment 的首部。当前分析的字符串一定是不包含换行的
// rewind rewinds the scanner's read position and character s.ch
// to the start of the currently active segment, which must not
// contain any newlines (otherwise position information will be
// incorrect). Currently, rewind is only needed for handling the
// source sequence ".."; it must not be called outside an active
// segment.
func (s *source) rewind() {
	// ok to verify precondition - rewind is rarely called
	if s.b < 0 {
		panic("no active segment")
	}
	s.col -= uint(s.r - s.b) // JAMLEE: 当前正在识别的字符串有多长
	s.r = s.b // JAMLEE: 从头开始
	s.nextch()
}

// JAMLEE: 读取下一个字符，到当前 ch 和 chw 中
func (s *source) nextch() {
redo:
	// JAMLEE: 当前的列位置和行位置一起定位用。
	s.col += uint(s.chw)
	if s.ch == '\n' { // JAMLEE: 当前字符是不是换行，是换行的话就重置 col
		s.line++
		s.col = 0
	}

	// JAMLEE: 如果是 ascii 码。这里只能表示1个index，rune把其强制转为rune类型了。s.ch 表示当前读取的utf8字符。
	// fast common case: at least one ASCII character
	if s.ch = rune(s.buf[s.r]); s.ch < sentinel {
		s.r++
		s.chw = 1
		if s.ch == 0 {
			s.error("invalid NUL character")
			goto redo
		}
		return
	}

	// slower general case: add more bytes to buffer if we don't have a full rune
	for s.e-s.r < utf8.UTFMax && !utf8.FullRune(s.buf[s.r:s.e]) && s.ioerr == nil {
		s.fill()
	}

	// EOF
	if s.r == s.e {
		if s.ioerr != io.EOF {
			// ensure we never start with a '/' (e.g., rooted path) in the error message
			s.error("I/O error: " + s.ioerr.Error())
			s.ioerr = nil
		}
		s.ch = -1
		s.chw = 0
		return
	}

	// JAMLEE: 解码 utf8 字符
	s.ch, s.chw = utf8.DecodeRune(s.buf[s.r:s.e])
	s.r += s.chw

	if s.ch == utf8.RuneError && s.chw == 1 {
		s.error("invalid UTF-8 encoding")
		goto redo
	}

	// JAMLEE: 如果是 BOM 则重新再读一次
	// BOM's are only allowed as the first character in a file
	const BOM = 0xfeff
	if s.ch == BOM {
		if s.line > 0 || s.col > 0 {
			s.error("invalid BOM in the middle of the file")
		}
		goto redo
	}
}

// JAMLEE: 读取更多的内容到 buf 中
// fill reads more source bytes into s.buf.
// It returns with at least one more byte in the buffer, or with s.ioerr != nil.
func (s *source) fill() {
	// determine content to preserve
	b := s.r
	if s.b >= 0 {
		b = s.b
		s.b = 0 // after buffer has grown or content has been moved down
	}
	content := s.buf[b:s.e] // JAMLEE: 当前还未解析的content（包括解析中未开始解析）

	// grow buffer or move content down
	if len(content)*2 > len(s.buf) { // JAMLEE: 如果未解析的超过了1半, double 当前 buf 的长度, content 拷贝到顶部
		s.buf = make([]byte, nextSize(len(s.buf)))
		copy(s.buf, content)
	} else if b > 0 { // JAMLEE: 如果没有超过一半 content 的值拷贝到顶部
		copy(s.buf, content)
	}
	s.r -= b // JAMLEE: r - b 意味着当前 r 的相对位置。加上0，也是拷贝后的buf从 0 开始了
	s.e -= b // JAMLEE: e - b 意味着当前 e 的相对位置。加上0，也是拷贝后的buf从 0 开始了

	// read more data: try a limited number of times
	for i := 0; i < 10; i++ {
		var n int
		n, s.ioerr = s.in.Read(s.buf[s.e : len(s.buf)-1]) // -1 to leave space for sentinel
		// JAMLEE: 读取时处理错误，这里 i 尝试时为什么呢
		if n < 0 {
			panic("negative read") // incorrect underlying io.Reader implementation
		}
		if n > 0 || s.ioerr != nil {
			s.e += n
			s.buf[s.e] = sentinel
			return
		}
		// JAMLEE: n == 0 时，重试
		// n == 0
	}

	s.buf[s.e] = sentinel
	s.ioerr = io.ErrNoProgress
}

// JAMLEE: 为 buf 返回 1 个合适的大小
// nextSize returns the next bigger size for a buffer of a given size.
func nextSize(size int) int {
	const min = 4 << 10 // 4K: minimum buffer size
	const max = 1 << 20 // 1M: maximum buffer size which is still doubled
	if size < min {
		return min
	}
	if size <= max {
		return size << 1
	}
	return size + max
}
