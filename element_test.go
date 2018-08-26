package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLines(t *testing.T) {
	s1 := `aaa
bbbb
`
	s2 := `aaa
bbbb`

	ls1 := NewLines(s1)
	assert.Equal(t, 2, ls1.Len())
	assert.Equal(t, NewLine("aaa\n").String(), ls1.At(0).String())
	assert.Equal(t, NewLine("bbbb\n").String(), ls1.At(1).String())
	assert.Equal(t, s1, ls1.Join())

	ls2 := NewLines(s2)
	assert.Equal(t, 2, ls2.Len())
	assert.Equal(t, NewLine("aaa\n").String(), ls1.At(0).String())
	assert.Equal(t, NewLine("bbbb\n").String(), ls1.At(1).String())
	assert.Equal(t, s2, ls2.Join())

	assert.Equal(t, NewLines("").Join(), ls1.Slice(0, 0).Join())
	assert.Equal(t, NewLines("aaa\n").Join(), ls1.Slice(0, 1).Join())
	assert.Equal(t, NewLines("bbbb\n").Join(), ls1.Slice(1, 2).Join())
	assert.Equal(t, ls1, ls1.Slice(0, ls1.Len()))
}

func TestLine(t *testing.T) {
	s := `aaa
bbbb
aaa
`
	ls := NewLines(s)
	l1 := ls.At(0)
	l2 := ls.At(1)
	l3 := ls.At(2)

	assert.Equal(t, "aaa\n", l1.String())
	assert.Equal(t, "bbbb\n", l2.String())
	assert.Equal(t, "aaa\n", l3.String())

	assert.True(t, l1.Equals(l3))
	assert.True(t, l3.Equals(l1))
	assert.False(t, l1.Equals(l2))
}
