package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
