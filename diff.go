package diff

import "fmt"

func TrimCommonElements(a, b Elements) (Elements, Elements) {
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
