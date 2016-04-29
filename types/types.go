// types
package types

type HCWORD int64

func (w HCWORD) Abs() uint64 {
	if w >= 0 {
		return uint64(w)
	} else {
		return uint64(-w)
	}
}

func MinWord(a, b HCWORD) HCWORD {
	if a < b {
		return a
	} else {
		return b
	}
}

func Swap(a, b *HCWORD) {
	*a, *b = *b, *a
}
