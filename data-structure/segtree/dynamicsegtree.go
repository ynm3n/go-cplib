package segtree

import "fmt"

// ポインタを使う 必要な部分だけ作るやつ
// 参考にさせていただいた記事:
// https://kazuma8128.hatenablog.com/entry/2018/11/29/093827
// https://lorent-kyopro.hatenablog.com/entry/2021/03/12/025644
type DynamicSegmentTree[T any] struct {
	l, r int
	root *node[T]
	e    func() T
	op   func(a, b T) T
}

func NewDynamicSegmentTree[T any](l, r int, e func() T, op func(a, b T) T) *DynamicSegmentTree[T] {
	return &DynamicSegmentTree[T]{
		l:    l,
		r:    r,
		root: nil,
		e:    e,
		op:   op,
	}
}

func NewDynamicSegmentTreeWith[T any](s []T, e func() T, op func(a, b T) T) *DynamicSegmentTree[T] {
	sg := NewDynamicSegmentTree(0, len(s), e, op)
	for i, val := range s {
		sg.Set(i, val)
	}
	return sg
}

func (sg *DynamicSegmentTree[T]) Set(i int, val T) {
	sg.checkInRange(i)
	if sg.root == nil {
		sg.root = newNode(i, val)
		return
	}

	p := (*node[T])(nil)
	now := sg.root
	for l, r := sg.l, sg.r; ; {
		if i == now.i {
			now.val = val
			break
		}
		if m := (l + r) / 2; i < m {
			if i > now.i {
				i, now.i = now.i, i
				val, now.val = now.val, val
			}
			p, now = now, now.l
			r = m
		} else {
			if i < now.i {
				i, now.i = now.i, i
				val, now.val = now.val, val
			}
			p, now = now, now.r
			l = m
		}
		if now == nil {
			now = newNode(i, val)
			now.p = p
			if now.i < p.i {
				p.l = now
			} else {
				p.r = now
			}
			now = p
			break
		}
	}

	for upd := now; upd != nil; upd = upd.p {
		sg.update(upd)
	}
}

func (sg *DynamicSegmentTree[T]) Get(i int) T {
	sg.checkInRange(i)
	for now := sg.root; now != nil; {
		if i == now.i {
			return now.val
		} else if i < now.i {
			now = now.l
		} else {
			now = now.r
		}
	}
	return sg.e()
}

func (sg *DynamicSegmentTree[T]) Product(l, r int) T {
	sg.checkInRangeLR(l, r)
	if sg.root == nil || l == r {
		return sg.e()
	}
	return sg.product(sg.root, l, r, sg.l, sg.r)
}

func (sg *DynamicSegmentTree[T]) product(n *node[T], argL, argR, l, r int) T {
	if n == nil || r <= argL || argR <= l {
		return sg.e()
	}
	if argL <= l && r <= argR {
		return n.subVal
	}
	res := sg.product(n.l, argL, argR, l, n.i)
	if argL <= n.i && n.i < argR {
		res = sg.op(res, n.val)
	}
	res = sg.op(res, sg.product(n.r, argL, argR, n.i+1, r))
	return res
}

func (sg *DynamicSegmentTree[T]) ProductAll() T {
	return sg.root.subVal
}

func (sg *DynamicSegmentTree[T]) update(n *node[T]) {
	n.subVal = sg.op(sg.subtreeVal(n.l), n.val)
	n.subVal = sg.op(n.subVal, sg.subtreeVal(n.r))
}

func (sg *DynamicSegmentTree[T]) subtreeVal(n *node[T]) T {
	if n != nil {
		return n.subVal
	}
	return sg.e()
}

func (sg *DynamicSegmentTree[T]) checkInRange(i int) {
	if i < sg.l || sg.r <= i {
		panic(fmt.Errorf("DynamicSegmentTree: index out of range: l=%d, r=%d, i=%d", sg.l, sg.r, i))
	}
}

func (sg *DynamicSegmentTree[T]) checkInRangeLR(argL, argR int) {
	if argL < sg.l || sg.r < argR {
		panic(fmt.Errorf("DynamicSegmentTree: index out of range: l=%d, r=%d, argL=%d, argR=%d", sg.l, sg.r, argL, argR))
	}
}

type node[T any] struct {
	i           int
	val, subVal T
	p, l, r     *node[T]
}

func newNode[T any](i int, val T) *node[T] {
	return &node[T]{
		i:      i,
		val:    val,
		subVal: val,
	}
}
