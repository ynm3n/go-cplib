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
	n         int
	blockSize int

	data      []S
	result    []S
	lazyApply []F
	isLazy    []bool

	e           func() S // https://ja.wikipedia.org/wiki/単位元
	productFunc func(x, y S) S

	id               func() F // https://ja.wikipedia.org/wiki/恒等写像
	mappingFunc      func(f F, x S) S
	mappingBlockFunc func(f F, x S) S
	compositionFunc  func(f, g F) F
}

func NewRangeSqrtDecomposition[S, F any](
	n int,
	e func() S,
	productFunc func(x, y S) S,
	id func() F,
	mappingFunc func(f F, x S) S,
	mappingBlockFunc func(f F, x S) S,
	compositionFunc func(f, g F) F,
) *RangeSqrtDecomposition[S, F] {
	data := make([]S, n)
	for i := range data {
		data[i] = e()
	}
	return NewRangeSqrtDecompositionWith(data, e, productFunc, id, mappingFunc, mappingBlockFunc, compositionFunc)
}

func NewRangeSqrtDecompositionWith[S, F any](
	data []S,
	e func() S,
	productFunc func(x, y S) S,
	id func() F,
	mappingFunc func(f F, x S) S,
	mappingBlockFunc func(f F, x S) S,
	compositionFunc func(f, g F) F,
) *RangeSqrtDecomposition[S, F] {
	data = slices.Clone(data)
	n := len(data)
	blockSize := int(math.Round(math.Sqrt(float64(n))))
	result := make([]S, (n+blockSize-1)/blockSize)
	lazyApply := make([]F, (n+blockSize-1)/blockSize)
	for i := 0; i < (n+blockSize-1)/blockSize; i++ {
		result[i] = e()
		lazyApply[i] = id()
	}
	if mappingBlockFunc == nil {
		mappingBlockFunc = mappingFunc
	}

	sd := &RangeSqrtDecomposition[S, F]{
		n:         n,
		blockSize: blockSize,

		data:      data,
		result:    result,
		lazyApply: lazyApply,
		isLazy:    make([]bool, (n+blockSize-1)/blockSize),

		e:           e,
		productFunc: productFunc,

		id:               id,
		mappingFunc:      mappingFunc,
		mappingBlockFunc: mappingBlockFunc,
		compositionFunc:  compositionFunc,
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
	sd.calcResult(b)
}

func (sd *RangeSqrtDecomposition[S, F]) Product(l, r int) S {
	lb, rb := sd.nowBlock(l), sd.nowBlock(r-1)
	lMax, rMin := (lb+1)*sd.blockSize, rb*sd.blockSize
	sd.evalLazy(lb)
	sd.evalLazy(rb)

	res := sd.e()
	for i, lEnd := l, min(lMax, r); i < lEnd; i++ {
		res = sd.productFunc(res, sd.data[i])
	}
	for b := lb + 1; b < rb; b++ {
		if sd.isLazy[b] {
			res = sd.productFunc(res, sd.mappingBlockFunc(sd.lazyApply[b], sd.result[b]))
		} else {
			res = sd.productFunc(res, sd.result[b])
		}
	}
	if lb != rb {
		for i := rMin; i < r; i++ {
			res = sd.productFunc(res, sd.data[i])
		}
	}
	return res
}

func (sd *RangeSqrtDecomposition[S, F]) Apply(i int, f F) {
	b := sd.nowBlock(i)
	sd.evalLazy(b)
	sd.data[i] = sd.mappingFunc(f, sd.data[i])
	sd.calcResult(b)
}

func (sd *RangeSqrtDecomposition[S, F]) ApplyRange(l, r int, f F) {
	lb, rb := sd.nowBlock(l), sd.nowBlock(r-1)
	lMax, rMin := (lb+1)*sd.blockSize, rb*sd.blockSize
	sd.evalLazy(lb)
	sd.evalLazy(rb)

	for i, lEnd := l, min(lMax, r); i < lEnd; i++ {
		sd.data[i] = sd.mappingFunc(f, sd.data[i])
	}
	sd.calcResult(lb)
	for b := lb + 1; b < rb; b++ {
		sd.lazyApply[b] = sd.compositionFunc(sd.lazyApply[b], f)
		sd.isLazy[b] = true
	}
	if lb != rb {
		for i := rMin; i < r; i++ {
			sd.data[i] = sd.mappingFunc(f, sd.data[i])
		}
		sd.calcResult(rb)
	}
}

func (sd *RangeSqrtDecomposition[S, F]) nowBlock(i int) int {
	return i / sd.blockSize
}

func (sd *RangeSqrtDecomposition[S, F]) evalLazy(b int) {
	if !sd.isLazy[b] {
		return
	}
	l := b * sd.blockSize
	r := min(l+sd.blockSize, sd.n)
	f := sd.lazyApply[b]
	for i := l; i < r; i++ {
		sd.data[i] = sd.mappingFunc(f, sd.data[i])
	}
	sd.isLazy[b] = false
	sd.lazyApply[b] = sd.id()
	sd.calcResult(b)
}

func (sd *RangeSqrtDecomposition[S, F]) calcResult(b int) {
	l := b * sd.blockSize
	r := min(l+sd.blockSize, sd.n)
	sd.result[b] = sd.e()
	for i := l; i < r; i++ {
		sd.result[b] = sd.productFunc(sd.result[b], sd.data[i])
	}
}
