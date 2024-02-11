package sqdecomp

import (
	"math"
	"slices"
)

// 区間を処理するために平方分割するやつ
// 参考にさせていただいた記事:
// https://kujira16.hateblo.jp/entry/2016/12/15/000000
// https://betrue12.hateblo.jp/entry/2020/09/22/194541
type RangeSqrtDecomposition[S, F any] struct {
	n     int // data len
	block int // block len

	data      []S
	result    []S
	modified  []bool
	lazyApply []F
	isLazy    []bool

	e       func() S // https://ja.wikipedia.org/wiki/単位元
	product func(x, y S) S

	id           func() F // https://ja.wikipedia.org/wiki/恒等写像
	mapping      func(f F, x S) S
	mappingBlock func(f F, x S) S
	composition  func(f, g F) F
}

func NewRangeSqrtDecomposition[S, F any](
	n int,
	e func() S,
	product func(x, y S) S,
	id func() F,
	mapping func(f F, x S) S,
	mappingBlock func(f F, x S) S,
	composition func(f, g F) F,
) *RangeSqrtDecomposition[S, F] {
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
) *RangeSqrtDecomposition[S, F] {
	data = slices.Clone(data)
	n := len(data)
	block := int(math.Round(math.Sqrt(float64(n))))
	if mappingBlock == nil {
		mappingBlock = mapping
	}

	var result []S
	if product != nil {
		result = make([]S, (n+block-1)/block)
		for i := 0; i < (n+block-1)/block; i++ {
			result[i] = e()
		}
	}
	lazyApply := make([]F, (n+block-1)/block)
	for i := 0; i < (n+block-1)/block; i++ {
		lazyApply[i] = id()
	}

	sd := &RangeSqrtDecomposition[S, F]{
		n:     n,
		block: block,

		data:      data,
		result:    result,
		modified:  make([]bool, len(result)),
		lazyApply: lazyApply,
		isLazy:    make([]bool, len(lazyApply)),

		e:       e,
		product: product,

		id:           id,
		mapping:      mapping,
		mappingBlock: mappingBlock,
		composition:  composition,
	}
	return sd
}

func (sd *RangeSqrtDecomposition[S, F]) Len() int {
	return sd.n
}

func (sd *RangeSqrtDecomposition[S, F]) Get(i int) S {
	sd.evalLazy(sd.nowBlock(i))
	return sd.data[i]
}

func (sd *RangeSqrtDecomposition[S, F]) Set(i int, x S) {
	b := sd.nowBlock(i)
	sd.evalLazy(b)
	sd.data[i] = x
	sd.flagModified(b)
}

func (sd *RangeSqrtDecomposition[S, F]) Product(l, r int) S {
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
		if sd.modified[b] {
			sd.calcResult(b)
			sd.modified[b] = false
		}
		if sd.isLazy[b] {
			res = sd.product(res, sd.mappingBlock(sd.lazyApply[b], sd.result[b]))
		} else {
			res = sd.product(res, sd.result[b])
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

func (sd *RangeSqrtDecomposition[S, F]) Apply(i int, f F) {
	b := sd.nowBlock(i)
	sd.evalLazy(b)
	sd.data[i] = sd.mapping(f, sd.data[i])
	sd.flagModified(b)
}

func (sd *RangeSqrtDecomposition[S, F]) ApplyRange(l, r int, f F) {
	lb, rb := sd.nowBlock(l), sd.nowBlock(r-1)
	if lb == rb {
		sd.evalLazy(lb)
		for i := l; i < r; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.flagModified(lb)
		return
	}

	lCeil := (l + sd.block - 1) / sd.block
	rFloor := r / sd.block
	if lMax := lCeil * sd.block; l < lMax {
		sd.evalLazy(lb)
		for i := l; i < lMax; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.flagModified(lb)
	}
	for b := lCeil; b < rFloor; b++ {
		sd.lazyApply[b] = sd.composition(sd.lazyApply[b], f)
		sd.isLazy[b] = true
	}
	if rMin := rFloor * sd.block; rMin < r {
		sd.evalLazy(rb)
		for i := rMin; i < r; i++ {
			sd.data[i] = sd.mapping(f, sd.data[i])
		}
		sd.flagModified(rb)
	}
}

func (sd *RangeSqrtDecomposition[S, F]) nowBlock(i int) int {
	return i / sd.block
}

func (sd *RangeSqrtDecomposition[S, F]) evalLazy(b int) {
	if !sd.isLazy[b] {
		return
	}
	l := b * sd.block
	r := min(l+sd.block, sd.n)
	f := sd.lazyApply[b]
	for i := l; i < r; i++ {
		sd.data[i] = sd.mapping(f, sd.data[i])
	}
	sd.isLazy[b] = false
	sd.lazyApply[b] = sd.id()
	sd.flagModified(b)
}

func (sd *RangeSqrtDecomposition[S, F]) calcResult(b int) {
	if sd.product == nil {
		return
	}
	l := b * sd.block
	r := min(l+sd.block, sd.n)
	sd.result[b] = sd.e()
	for i := l; i < r; i++ {
		sd.result[b] = sd.product(sd.result[b], sd.data[i])
	}
}

func (sd *RangeSqrtDecomposition[S, F]) flagModified(b int) {
	if sd.product == nil {
		return
	}
	sd.modified[b] = true
}
