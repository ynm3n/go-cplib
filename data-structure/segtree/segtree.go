package segtree

type SegmentTree[T any] interface {
	Len() int
	Set(i int, val T)
	Get(i int) T
	Product(l, r int) T
	ProductAll() T
}

func NewSegmentTree[T any](s []T, op func(a, b T) T, unit T) SegmentTree[T] {
	sg := &segmentTree[T]{
		n:         len(s),
		data:      make([]T, len(s)*2),
		operation: op,
		unit:      unit,
	}
	for i := 0; i < sg.n; i++ {
		sg.data[i] = sg.unit
	}
	for i, v := range s {
		sg.data[i+sg.n] = v
	}
	for i := sg.n - 1; i >= 1; i-- {
		sg.update(i)
	}
	return sg
}

// 参考にさせていただいた記事:
// https://maspypy.com/segment-tree-%e3%81%ae%e3%81%8a%e5%8b%89%e5%bc%b71
// https://github.com/ktateish/go-competitive/blob/master/ac_segtree.go2
// https://github.com/monkukui/ac-library-go/blob/master/segtree/segtree.go
// https://atcoder.github.io/ac-library/master/document_ja/segtree.html
type segmentTree[T any] struct {
	n         int            // 初期化時に渡したスライスの要素数
	data      []T            // セグメントツリーの本体
	operation func(a, b T) T // モノイド積(?)を計算する関数
	unit      T              // 単位元
}

func (sg *segmentTree[T]) Len() int {
	return sg.n
}

func (sg *segmentTree[T]) Set(idx int, v T) {
	now := sg.n + idx
	sg.data[now] = v
	for now > 1 {
		now /= 2
		sg.update(now)
	}
}

func (sg *segmentTree[T]) Get(idx int) T {
	return sg.data[sg.n+idx]
}

func (sg *segmentTree[T]) Product(l, r int) T {
	l += sg.n
	r += sg.n
	valL, valR := sg.unit, sg.unit
	for l < r {
		if l%2 == 1 {
			valL = sg.operation(valL, sg.data[l])
			l++
		}
		if r%2 == 1 {
			r--
			valR = sg.operation(sg.data[r], valR)
		}
		l /= 2
		r /= 2
	}
	return sg.operation(valL, valR)
}

func (sg *segmentTree[T]) ProductAll() T {
	return sg.Product(0, sg.n)
}

func (sg *segmentTree[T]) update(now int) {
	child1, child2 := now*2, now*2+1
	sg.data[now] = sg.operation(sg.data[child1], sg.data[child2])
}
