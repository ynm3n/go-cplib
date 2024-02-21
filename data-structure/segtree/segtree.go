package segtree

type SegmentTree[T any] interface {
	Len() int
	Set(i int, val T)
	Get(i int) T
	Product(l, r int) T
	ProductAll() T
}

func NewSegmentTree[T any](n int, e func() T, op func(a, b T) T) SegmentTree[T] {
	data := make([]T, n*2)
	for i := range data {
		data[i] = e()
	}
	sg := &segmentTree[T]{
		n:    n,
		data: data,
		e:    e,
		op:   op,
	}
	for i := sg.n - 1; i >= 1; i-- {
		sg.update(i)
	}
	return sg
}

func NewSegmentTreeWith[T any](s []T, e func() T, op func(a, b T) T) SegmentTree[T] {
	data := make([]T, len(s)*2)
	for i := 0; i < len(s); i++ {
		data[i] = e()
	}
	copy(data[len(s):], s)
	sg := &segmentTree[T]{
		n:    len(s),
		data: data,
		e:    e,
		op:   op,
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
	n    int            // 初期化時に渡したスライスの要素数
	data []T            // セグメントツリーの本体
	e    func() T       // 単位元
	op   func(a, b T) T // モノイド積(?)を計算する関数
}

func (sg *segmentTree[T]) Len() int {
	return sg.n
}

func (sg *segmentTree[T]) Set(i int, val T) {
	now := sg.n + i
	sg.data[now] = val
	for now > 1 {
		now /= 2
		sg.update(now)
	}
}

func (sg *segmentTree[T]) Get(i int) T {
	return sg.data[sg.n+i]
}

func (sg *segmentTree[T]) Product(l, r int) T {
	l += sg.n
	r += sg.n
	valL, valR := sg.e(), sg.e()
	for l < r {
		if l%2 == 1 {
			valL = sg.op(valL, sg.data[l])
			l++
		}
		if r%2 == 1 {
			r--
			valR = sg.op(sg.data[r], valR)
		}
		l /= 2
		r /= 2
	}
	return sg.op(valL, valR)
}

func (sg *segmentTree[T]) ProductAll() T {
	return sg.Product(0, sg.n)
}

func (sg *segmentTree[T]) update(now int) {
	child1, child2 := now*2, now*2+1
	sg.data[now] = sg.op(sg.data[child1], sg.data[child2])
}
