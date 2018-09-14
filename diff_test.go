package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	sa := `
    David Axelrod
    Electric Prunes
    Gil Scott Heron
    The Slits
    Faust
    The Sonics
    The Sonics
`

	sb := `
    The Slits
    Gil Scott Heron
    David Axelrod
    Electric Prunes
    Faust
    The Sonics
    The Sonics
`

	a := NewDocument(sa)
	b := NewDocument(sb)
	edits := Diff(a, b)
	for _, e := range edits {
		t.Log(e)
	}
}

func TestUniqueElements(t *testing.T) {
	lines := NewDocument(`111
22222
333

111
1111
11
`)
	uniques := UniqueElements(lines)
	want := map[uint64][]int{
		NewLine("22222\n").Hash(): []int{1},
		NewLine("333\n").Hash():   []int{2},
		NewLine("\n").Hash():      []int{3},
		NewLine("1111\n").Hash():  []int{5},
		NewLine("11\n").Hash():    []int{6},
	}
	assert.Equal(t, want, uniques)
}

func TestTrimCommonElements(t *testing.T) {
	l1 := NewDocument(`111
22222
333

111
1111

11
`)
	l2 := NewDocument(`111
22222
3335

111
111155

11
`)
	tl1, tl2 := TrimCommonElements(l1, l2)
	el1 := NewDocument(`333

111
1111
`)
	el2 := NewDocument(`3335

111
111155
`)
	assert.Equal(t, el1.Join(), tl1.Join())
	assert.Equal(t, el2.Join(), tl2.Join())
}
