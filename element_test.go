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
	assert.Equal(t, NewLine("aaa\n"), ls1.At(0))
	assert.Equal(t, NewLine("bbbb\n"), ls1.At(1))
	assert.Equal(t, s1, ls1.String())

	ls2 := NewLines(s2)
	assert.Equal(t, 2, ls2.Len())
	assert.Equal(t, NewLine("aaa\n"), ls1.At(0))
	assert.Equal(t, NewLine("bbbb\n"), ls1.At(1))
	assert.Equal(t, s2, ls2.String())

	assert.Equal(t, NewLines("").String(), ls1.Slice(0, 0).String())
	assert.Equal(t, NewLines("aaa\n").String(), ls1.Slice(0, 1).String())
	assert.Equal(t, NewLines("bbbb\n").String(), ls1.Slice(1, 2).String())
	assert.Equal(t, ls1, ls1.Slice(0, ls1.Len()))
}
