package diff

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

type Element interface {
	Equals(e Element) bool
	Hash() uint64
	fmt.Stringer
}

type Elements interface {
	Len() int
	At(i int) Element
	Slice(start, end int) Elements
	fmt.Stringer
}

func NewDocument(doc string) Elements {
	return NewLines(doc)
}

func NewLines(s string) Lines {
	raw := []byte(s)
	const bytesPerLine = 36
	indexOfLine := make([]int, 1, len(raw)/bytesPerLine+1)

	l := len(raw)
	if l == 0 {
		return Lines{}
	}

	indexOfLine[0] = 0
	for offset := 0; offset < l; {
		idx := bytes.IndexByte(raw[offset:], '\n')
		if idx == -1 {
			break
		}
		offset += idx + 1
		indexOfLine = append(indexOfLine, offset)
	}
	if raw[l-1] != '\n' {
		indexOfLine = append(indexOfLine, l)
	}

	return Lines{raw, indexOfLine}
}

type Lines struct {
	raw         []byte
	indexOfLine []int
}

func (l Lines) Len() int {
	return len(l.indexOfLine) - 1
}

func (l Lines) At(i int) Element {
	return NewLineRaw(l.raw[l.indexOfLine[i]:l.indexOfLine[i+1]])
}

func (l Lines) Slice(start, end int) Elements {
	indexOfLine := l.indexOfLine[start : end+1]
	return Lines{
		l.raw,
		indexOfLine,
	}
}

func (l Lines) String() string {
	var sb strings.Builder
	For(l, func(i int, e Element) {
		sb.WriteString(e.String())
	})
	return sb.String()
}

func NewLineRaw(raw []byte) Line {
	return Line(raw)
}

func NewLine(s string) Line {
	return NewLineRaw([]byte(s))
}

type Line []byte

func (l Line) Equals(e Element) bool {
	l2, ok := e.(Line)
	if !ok {
		return false
	}
	return bytes.Equal(l, l2)
}

func (l Line) runes() []rune {
	return []rune(string(l))
}

func (l Line) Len() int {
	return len(l.runes())
}

func (l Line) At(i int) Element {
	return NewRune(l.runes()[i])
}

func (l Line) String() string {
	return string(l)
}

func (l Line) Slice(start, end int) Elements {
	return Line(l[start:end])
}

func (l Line) Hash() uint64 {
	h := fnv.New64()
	h.Write(l)
	return h.Sum64()
}

func NewRune(r rune) Rune {
	return Rune(r)
}

type Rune rune

func (r Rune) Equals(e Element) bool {
	r2, ok := e.(Rune)
	if !ok {
		return false
	}
	return r == r2
}

func (r Rune) Hash() uint64 {
	return uint64(r)
}

func (r Rune) String() string {
	return string(r)
}
