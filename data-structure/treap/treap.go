package treap

import (
	"cmp"
	"math/rand"
	"time"
)

type Set[K any] interface {
	Len() int
	Set(k K)
	Get(k K) bool
	Delete(k K)
	Min() (K, bool)
	Max() (K, bool)

	// k以下かk未満な要素を探す
	// eqでイコールを許すかどうかを指定
	SearchLeft(k K, eq bool) (K, bool)
	// k以上かkより大きい要素を探す
	// eqでイコールを許すかどうかを指定
	SearchRight(k K, eq bool) (K, bool)
}

func NewSet[K cmp.Ordered]() Set[K] {
	return NewSetFunc(func(x, y K) int {
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
		return 0
	})
}

func NewSetFunc[K any](cmp func(x, y K) int) Set[K] {
	return &set[K]{NewOrderedMapFunc[K, struct{}](cmp)}
}

type set[K any] struct {
	m OrderedMap[K, struct{}]
}

func (st *set[K]) Len() int {
	return st.m.Len()
}

func (st *set[K]) Set(k K) {
	st.m.Set(k, struct{}{})
}

func (st *set[K]) Get(k K) bool {
	_, b := st.m.Get(k)
	return b
}

func (st *set[K]) Delete(k K) {
	st.m.Delete(k)
}

func (st *set[K]) Min() (K, bool) {
	k, _, b := st.m.Min()
	return k, b
}

func (st *set[K]) Max() (K, bool) {
	k, _, b := st.m.Max()
	return k, b
}

// k以下かk未満な要素を探す
// eqでイコールを許すかどうかを指定
func (st *set[K]) SearchLeft(k K, eq bool) (K, bool) {
	k, _, b := st.m.SearchLeft(k, eq)
	return k, b
}

// k以上かkより大きい要素を探す
// eqでイコールを許すかどうかを指定
func (st *set[K]) SearchRight(k K, eq bool) (K, bool) {
	k, _, b := st.m.SearchRight(k, eq)
	return k, b
}

type OrderedMap[K, V any] interface {
	Len() int
	Set(k K, v V)
	Get(k K) (V, bool)
	Delete(k K)
	Min() (K, V, bool)
	Max() (K, V, bool)

	// k以下かk未満な要素を探す
	// eqでイコールを許すかどうかを指定
	SearchLeft(k K, eq bool) (K, V, bool)
	// k以上かkより大きい要素を探す
	// eqでイコールを許すかどうかを指定
	SearchRight(k K, eq bool) (K, V, bool)
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return NewOrderedMapFunc[K, V](func(x, y K) int {
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
		return 0
	})
}

func NewOrderedMapFunc[K, V any](cmp func(x, y K) int) OrderedMap[K, V] {
	tr := &treap[K, V]{
		0,
		nil,
		cmp,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return tr
}

// キーの重複を許さない仕様
type treap[K, V any] struct {
	len  int
	root *treapNode[K, V]
	cmp  func(x, y K) int
	rnd  *rand.Rand
}

type treapNode[K, V any] struct {
	k                      K
	v                      V
	priority               uint64 // 優先度が高いノードを根に近い位置に置く
	parent, childL, childR *treapNode[K, V]
}

func (tr *treap[K, V]) newTreapNode(k K, v V, parent, childL, childR *treapNode[K, V]) *treapNode[K, V] {
	nd := &treapNode[K, V]{
		k,
		v,
		tr.rnd.Uint64(),
		parent, childL, childR,
	}
	return nd
}

// find returns (parent, child)
func (tr *treap[K, V]) find(k K) (*treapNode[K, V], *treapNode[K, V]) {
	var p *treapNode[K, V]
	c := tr.root
	for c != nil {
		switch cm := tr.cmp(k, c.k); {
		case cm < 0:
			p = c
			c = c.childL
		case cm > 0:
			p = c
			c = c.childR
		case cm == 0:
			return p, c
		}
	}
	return p, c
}

func (tr *treap[K, V]) rotate(p, c *treapNode[K, V]) {
	if p == tr.root {
		tr.root = c
	}
	if p.childL == c {
		p.childL = c.childR
		if c.childR != nil {
			c.childR.parent = p
		}
		c.childR = p
	} else {
		p.childR = c.childL
		if c.childL != nil {
			c.childL.parent = p
		}
		c.childL = p
	}
	c.parent = p.parent
	if p.parent != nil {
		if p.parent.childL == p {
			p.parent.childL = c
		} else {
			p.parent.childR = c
		}
	}
	p.parent = c
}

func (tr *treap[K, V]) prev(nd *treapNode[K, V]) (*treapNode[K, V], bool) {
	if nd == nil {
		return nil, false
	}
	if nd.childL != nil {
		nd = nd.childL
		for nd.childR != nil {
			nd = nd.childR
		}
		return nd, true
	}
	for nd.parent != nil && nd == nd.parent.childL {
		nd = nd.parent
	}
	if nd.parent != nil {
		return nd.parent, true
	}
	return nil, false
}

func (tr *treap[K, V]) next(nd *treapNode[K, V]) (*treapNode[K, V], bool) {
	if nd == nil {
		return nil, false
	}
	if nd.childR != nil {
		nd = nd.childR
		for nd.childL != nil {
			nd = nd.childL
		}
		return nd, true
	}
	for nd.parent != nil && nd == nd.parent.childR {
		nd = nd.parent
	}
	if nd.parent != nil {
		return nd.parent, true
	}
	return nil, false
}

func (tr *treap[K, V]) Len() int {
	return tr.len
}

func (tr *treap[K, V]) Set(k K, v V) {
	if tr.Len() == 0 {
		tr.root = tr.newTreapNode(k, v, nil, nil, nil)
		tr.len++
		return
	}
	par, now := tr.find(k)
	if now != nil {
		now.v = v
		return
	}
	now = tr.newTreapNode(k, v, par, nil, nil)
	tr.len++
	if cm := tr.cmp(k, par.k); cm < 0 {
		par.childL = now
	} else {
		par.childR = now
	}
	for par != nil && par.priority < now.priority {
		tr.rotate(par, now)
		par = now.parent
	}
	if par == nil {
		tr.root = now
	}
}

func (tr *treap[K, V]) Get(k K) (V, bool) {
	if _, now := tr.find(k); now != nil {
		return now.v, true
	}
	return *new(V), false
}

func (tr *treap[K, V]) Delete(k K) {
	_, now := tr.find(k)
	if now == nil {
		return
	}
	for now.childL != nil || now.childR != nil {
		switch {
		case now.childL != nil && now.childR != nil:
			if now.childL.priority > now.childR.priority {
				tr.rotate(now, now.childL)
			} else {
				tr.rotate(now, now.childR)
			}
		case now.childL != nil:
			tr.rotate(now, now.childL)
		case now.childR != nil:
			tr.rotate(now, now.childR)
		}
	}
	if now.parent != nil {
		if now == now.parent.childL {
			now.parent.childL = nil
		} else {
			now.parent.childR = nil
		}
	}
	if now == tr.root {
		tr.root = nil
	}
	*now = treapNode[K, V]{}
	tr.len--
}

func (tr *treap[K, V]) Min() (K, V, bool) {
	if tr.Len() == 0 {
		return *new(K), *new(V), false
	}
	now := tr.root
	for now.childL != nil {
		now = now.childL
	}
	return now.k, now.v, true
}

func (tr *treap[K, V]) Max() (K, V, bool) {
	if tr.Len() == 0 {
		return *new(K), *new(V), false
	}
	now := tr.root
	for now.childR != nil {
		now = now.childR
	}
	return now.k, now.v, true
}

// k以下かk未満な要素を探す
// eqでイコールを許すかどうかを指定
func (tr *treap[K, V]) SearchLeft(k K, eq bool) (K, V, bool) {
	par, now := tr.find(k)
	if now != nil {
		if eq {
			return k, now.v, true
		} else if prev, b := tr.prev(now); b {
			return prev.k, prev.v, true
		}
	} else if par != nil {
		if tr.cmp(k, par.k) > 0 {
			return par.k, par.v, true
		} else if prev, b := tr.prev(par); b {
			return prev.k, prev.v, true
		}
	}
	return *new(K), *new(V), false
}

// k以上かkより大きい要素を探す
// eqでイコールを許すかどうかを指定
func (tr *treap[K, V]) SearchRight(k K, eq bool) (K, V, bool) {
	par, now := tr.find(k)
	if now != nil {
		if eq {
			return k, now.v, true
		} else if next, b := tr.next(now); b {
			return next.k, next.v, true
		}
	} else if par != nil {
		if tr.cmp(k, par.k) < 0 {
			return par.k, par.v, true
		} else if next, b := tr.next(par); b {
			return next.k, next.v, true
		}
	}
	return *new(K), *new(V), false
}
