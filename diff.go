package diff

import (
	"fmt"
	"sort"
)

type Document interface {
	Len() int
	At(i int) Element
	Slice(start, end int) Document
	Join() string
	AbsoluteRange() (from, to int)
}

func Diff(a, b Document) []Edit {
	a, b = TrimCommonElements(a, b)
	ua := UniqueElements(a)
	ub := UniqueElements(b)

	m := make(map[uint64]struct{})
	bs := make([]int, 0, len(ub))
	btoa := make(map[int]int, len(ua))
	For(a, func(i int, e Element) {
		h := e.Hash()
		_, ok := m[h]
		if ok {
			return
		}
		m[h] = struct{}{}

		idxsb, ok := ub[h]
		if !ok {
			return
		}
		for _, idxb := range idxsb {
			if e.Equals(b.At(idxb)) {
				bs = append(bs, idxb)
				btoa[idxb] = i
			}
		}
	})

	lisb := append(LongestIncreasingSubsequence(bs), b.Len())
	btoa[b.Len()] = a.Len()

	var edits []Edit
	preva, prevb := 0, 0
	for _, idxb := range lisb {
		idxa := btoa[idxb]
		sa := a.Slice(preva, idxa)
		sb := b.Slice(prevb, idxb)
		e := NewEdit(sa, sb)
		if e.Type() != Empty {
			edits = append(edits, e)
		}
		preva = idxa + 1
		prevb = idxb + 1
	}

	return edits
}

func NewEdit(a, b Document) Edit {
	return Edit{a, b}
}

type Edit struct {
	A Document
	B Document
}

func (e Edit) Type() EditType {
	la, lb := e.A.Len(), e.B.Len()
	switch {
	case la == 0 && lb == 0:
		return Empty
	case la == 0:
		return Insert
	case lb == 0:
		return Delete
	default:
		return Replace
	}
}

func (e Edit) String() string {
	sa, ea := e.A.AbsoluteRange()
	sb, eb := e.B.AbsoluteRange()
	switch e.Type() {
	case Empty:
		return s
	case Insert:
		return s + fmt.Sprintf("+ %q", e.B.Join())
	case Delete:
		return s + fmt.Sprintf("- %q", e.A.Join())
	case Replace:
		return s + fmt.Sprintf("- %q\n+ %q", e.A.Join(), e.B.Join())
	default:
		panic(fmt.Sprintf("invalid edit type: %d", e.Type()))
	}
}

type EditType int

func (et EditType) String() string {
	switch et {
	case Empty:
		return "empty"
	case Insert:
		return "insert"
	case Delete:
		return "delete"
	case Replace:
		return "replace"
	default:
		panic(fmt.Sprintf("invalid edit type: %d", et))
	}
}

const (
	Empty EditType = iota
	Insert
	Delete
	Replace
)

func UniqueElements(es Document) map[uint64][]int {
	m := make(map[uint64][]int)
	For(es, func(i int, e Element) {
		h := e.Hash()
		idxs, ok := m[h]
		if !ok {
			m[h] = []int{i}
			return
		}
		for j, idx := range idxs {
			if es.At(idx).Equals(es.At(i)) {
				// remove idx from idxs
				l := len(idxs)
				idxs[j] = idxs[l-1]
				idxs = idxs[:l-1]
				if len(idxs) == 0 {
					delete(m, h)
				} else {
					m[h] = idxs
				}
				return
			}
		}
		m[h] = append(m[h], i)
	})
	return m
}

func For(es Document, f func(i int, e Element)) {
	l := es.Len()
	for i := 0; i < l; i++ {
		f(i, es.At(i))
	}
}

func TrimCommonElements(a, b Document) (Document, Document) {
	sa, ea := 0, a.Len()
	sb, eb := 0, b.Len()
	for sa < ea && sb < eb && a.At(sa).Equals(b.At(sb)) {
		sa++
		sb++
	}

	for sa < ea && sb < eb && a.At(ea-1).Equals(b.At(eb-1)) {
		ea--
		eb--
	}

	return a.Slice(sa, ea), b.Slice(sb, eb)
}

func LongestIncreasingSubsequence(vs []int) []int {
	var stacks [][]int

	backPointer := make(map[int]int, len(vs))
	for _, v := range vs {
		i := sort.Search(len(stacks), func(i int) bool { return stacks[i][len(stacks[i])-1] >= v })
		if i < len(stacks) {
			stacks[i] = append(stacks[i], v)
		} else {
			stacks = append(stacks, []int{v})
		}
		if i != 0 {
			backPointer[v] = stacks[i-1][len(stacks[i-1])-1]
		}
	}

	stack := stacks[len(stacks)-1]
	v := stack[len(stack)-1]
	var lis []int
	for {
		lis = append(lis, v)
		x, ok := backPointer[v]
		if !ok {
			break
		}
		v = x
	}
	reverse(lis)

	return lis
}

func reverse(vs []int) {
	l := len(vs)
	for i := 0; i < l/2; i++ {
		vs[i], vs[l-i-1] = vs[l-i-1], vs[i]
	}
}
