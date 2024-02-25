package sqdecomp

import (
	"math"
	"slices"
)

type RangeSqrtDecomposition[S, F any] interface {
	Len() int
	Set(i int, x S)
	Get(i int) S
	Product(l, r int) S
	Apply(i int, f F)
	ApplyRange(l, r int, f F)
}

func NewRangeSqrtDecomposition[S, F any](
	n int,
	e func() S,
	product func(x, y S) S,
	id func() F,
	mapping func(f F, x S) S,
	mappingBlock func(f F, x S) S,
	composition func(f, g F) F,
) RangeSqrtDecomposition[S, F] {
	data := make([]S, n)
	for i := range data {
		data[i] = e()
	}
	return NewRangeSqrtDecompositionWith(data, e, product, id, mapping, mappingBlock, composition)
}

func NewRangeSqrtDecompositionWith[S, F any](
	data []S,
	e func() S,
	product func(x, y S) S,
	id func() F,
	mapping func(f F, x S) S,
	mappingBlock func(f F, x S) S,
	composition func(f, g F) F,
) RangeSqrtDecomposition[S, F] {
	data = slices.Clone(data)
	n := len(data)
	block := int(math.Round(math.Sqrt(float64(n))))
	if mappingBlock == nil {
		mappingBlock = mapping
	}

	var cache []S
	if product != nil {
		cache = make([]S, (n+block-1)/block)
		for i := range cache {
			cache[i] = e()
		}
	}
	lazy := make([]F, (n+block-1)/block)
	for i := range lazy {
		lazy[i] = id()
	}

	sd := &rangeSqrtDecomposition[S, F]{
		n:     n,
		block: block,

		data:      data,
		cache:     cache,
		isFresh:   make([]bool, len(cache)),
		lazy:      lazy,
		isWaiting: make([]bool, len(lazy)),

		e:       e,
		product: product,

		id:           id,
		mapping:      mapping,
		mappingBlock: mappingBlock,
		composition:  composition,
	}
	return sd
}

// 区間を処理するために平方分割するやつ
// 参考にさせていただいた記事:
// https://kujira16.hateblo.jp/entry/2016/12/15/000000
// https://betrue12.hateblo.jp/entry/2020/09/22/194541
type rangeSqrtDecomposition[S, F any] struct {
	n     int // data len
	block int // block len

	data      []S
	cache     []S
	isFresh   []bool // cacheが新鮮かどうか
	lazy      []F
	isWaiting []bool // 未評価lazyが待機しているかどうか

	e       func() S // https://ja.wikipedia.org/wiki/単位元
	product func(x, y S) S

	id           func() F // https://ja.wikipedia.org/wiki/恒等写像
	mapping      func(f F, x S) S
	mappingBlock func(f F, x S) S
	composition  func(f, g F) F
}

func (sd *rangeSqrtDecomposition[S, F]) Len() int {
	return sd.n
}

func (sd *rangeSqrtDecomposition[S, F]) Set(i int, x S) {
	b := sd.nowBlock(i)
	sd.evalLazy(b)
	sd.data[i] = x
	sd.markAsStale(b)
}

func (sd *rangeSqrtDecomposition[S, F]) Get(i int) S {
	sd.evalLazy(sd.nowBlock(i))
	return sd.data[i]
}

func (sd *rangeSqrtDecomposition[S, F]) Product(l, r int) S {
	lb, rb := sd.nowBlock(l), sd.nowBlock(r-1)
	res := sd.e()
	if lb == rb {
		sd.evalLazy(lb)
		for i := l; i < r; i++ {
			res = sd.product(res, sd.data[i])
		}
		return res
	}

	lCeil := (l + sd.block - 1) / sd.block
	rFloor := r / sd.block
	if lMax := lCeil * sd.block; l < lMax {
		sd.evalLazy(lb)
		for i := l; i < lMax; i++ {
			res = sd.product(res, sd.data[i])
		}
	}
	for b := lCeil; b < rFloor; b++ {
		if !sd.isFresh[b] && sd.product != nil {
			il := b * sd.block
			ir := min(il+sd.block, sd.n)
			sd.cache[b] = sd.e()
			for i := il; i < ir; i++ {
				sd.cache[b] = sd.product(sd.cache[b], sd.data[i])
			}
			sd.isFresh[b] = true
		}
		if sd.isWaiting[b] {
			res = sd.product(res, sd.mappingBlock(sd.lazy[b], sd.cache[b]))
		} else {
			res = sd.product(res, sd.cache[b])
		}
	}
	if rMin := rFloor * sd.block; rMin < r {
		sd.evalLazy(rb)
		for i := rMin; i < r; i++ {
			res = sd.product(res, sd.data[i])
		}
	}
	return res
}

func (sd *rangeSqrtDecomposition[S, F]) Apply(i int, f F) {
	b := sd.nowBlock(i)
	sd.evalLazy(b)
	sd.data[i] = sd.mapping(f, sd.data[i])
	sd.markAsStale(b)
}

func (sd *rangeSqrtDecomposition[S, F]) ApplyRange(l, r int, f F) {
	lb, rb := sd.nowBlock(l), sd.nowBlock(r-1)
	if lb == rb {
		sd.evalLazy(lb)
		for i := l; i < r; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.markAsStale(lb)
		return
	}

	lCeil := (l + sd.block - 1) / sd.block
	rFloor := r / sd.block
	if lMax := lCeil * sd.block; l < lMax {
		sd.evalLazy(lb)
		for i := l; i < lMax; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.markAsStale(lb)
	}
	for b := lCeil; b < rFloor; b++ {
		sd.lazy[b] = sd.composition(sd.lazy[b], f)
		sd.isWaiting[b] = true
	}
	if rMin := rFloor * sd.block; rMin < r {
		sd.evalLazy(rb)
		for i := rMin; i < r; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.markAsStale(rb)
	}
}

func (sd *rangeSqrtDecomposition[S, F]) nowBlock(i int) int {
	return i / sd.block
}

func (sd *rangeSqrtDecomposition[S, F]) evalLazy(b int) {
	if !sd.isWaiting[b] {
		return
	}
	l := b * sd.block
	r := min(l+sd.block, sd.n)
	f := sd.lazy[b]
	for i := l; i < r; i++ {
		sd.data[i] = sd.mapping(f, sd.data[i])
	}
	sd.isWaiting[b] = false
	sd.lazy[b] = sd.id()
	sd.markAsStale(b)
}

func (sd *rangeSqrtDecomposition[S, F]) markAsStale(b int) {
	if sd.product == nil {
		return
	}
	sd.isFresh[b] = false
}
