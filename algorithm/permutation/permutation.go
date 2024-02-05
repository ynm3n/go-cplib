package permutation

import (
	"cmp"
	"slices"
)

// 参考にさせていただいたコード:
// https://atcoder.jp/contests/abc276/submissions/37305540
type Permutation[T cmp.Ordered] []T

func NewPermutation(n int) Permutation[int] {
	p := make(Permutation[int], n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return p
}

func NewPermutationInitialized[T cmp.Ordered](s []T) Permutation[T] {
	slices.Sort(s)
	return Permutation[T](s)
}

func (p Permutation[T]) Next() bool {
	return p.doPermutation(cmp.Less[T])
}

func (p Permutation[T]) Prev() bool {
	return p.doPermutation(func(a, b T) bool { return cmp.Less(b, a) })
}

func (p Permutation[T]) doPermutation(cmp func(T, T) bool) bool {
	n := len(p)
	idx := -1
	for i := n - 2; i >= 0; i-- {
		if cmp(p[i], p[i+1]) {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}

	for l, r := idx+1, n-1; l < r; {
		p[l], p[r] = p[r], p[l]
		l++
		r--
	}
	for i := idx + 1; i < n; i++ {
		if cmp(p[idx], p[i]) {
			p[idx], p[i] = p[i], p[idx]
			break
		}
	}
	return true
}
