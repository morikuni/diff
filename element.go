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

func NewDocument(doc string) Document {
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

	return Lines{raw, indexOfLine, 0, len(indexOfLine)}
}

type Lines struct {
	raw         []byte
	indexOfLine []int
	start       int
	end         int
}

var _ interface {
	Document
} = Lines{}

func (l Lines) Len() int {
	return l.end - l.start - 1
}

func (l Lines) At(i int) Element {
	return NewLineRaw(l.raw, l.indexOfLine[l.start+i], l.indexOfLine[l.start+i+1])
}

func (l Lines) Slice(start, end int) Document {
	x := Lines{
		l.raw,
		l.indexOfLine,
		l.start + start,
		l.start + end + 1,
	}
	return x
}

func (l Lines) Join() string {
	var sb strings.Builder
	For(l, func(i int, e Element) {
		sb.WriteString(e.String())
	})
	return sb.String()
}

func (l Lines) AbsoluteRange() (int, int) {
	return l.start, l.end
}

func NewLineRaw(raw []byte, start, end int) Line {
	return Line{raw, start, end}
}

func NewLine(s string) Line {
	return NewLineRaw([]byte(s), 0, len(s))
}

type Line struct {
	raw   []byte
	start int
	end   int
}

var _ interface {
	Document
	Element
} = Line{}

func (l Line) Equals(e Element) bool {
	l2, ok := e.(Line)
	if !ok {
		return false
	}
	return bytes.Equal(l.bytes(), l2.bytes())
}

func (l Line) bytes() []byte {
	return l.raw[l.start:l.end]
}

func (l Line) runes() []rune {
	return []rune(string(l.bytes()))
}

func (l Line) Len() int {
	return len(l.runes())
}

func (l Line) At(i int) Element {
	return NewRune(l.runes()[i])
}

func (l Line) Join() string {
	return l.String()
}

func (l Line) String() string {
	return string(l.bytes())
}

func (l Line) Slice(start, end int) Document {
	return Line{l.raw, l.start + start, l.start + end}
}

func (l Line) AbsoluteRange() (int, int) {
	return l.start, l.end
}

func (l Line) Hash() uint64 {
	h := fnv.New64()
	h.Write(l.bytes())
	return h.Sum64()
}

func NewRune(r rune) Rune {
	return Rune(r)
}

type Rune rune

var _ interface {
	Element
} = Rune(0)

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
