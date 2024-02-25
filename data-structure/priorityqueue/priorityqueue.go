package priorityqueue

import (
	"cmp"
	"container/heap"
)

type PriorityQueue[T any] struct {
	h *hp[T]
}

// 昇順
func NewPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	pq := NewPriorityQueueFunc(cmp.Compare[T])
	return pq
}

func NewPriorityQueueFunc[T any](compare func(a, b T) int) *PriorityQueue[T] {
	pq := new(PriorityQueue[T])
	pq.h.s = make([]T, 0)
	pq.h.compare = compare
	return pq
}

func (pq *PriorityQueue[T]) Len() int    { return pq.h.Len() }
func (pq *PriorityQueue[T]) Top() T      { return pq.h.s[0] }
func (pq *PriorityQueue[T]) Enqueue(x T) { heap.Push(pq.h, x) }
func (pq *PriorityQueue[T]) Dequeue() T  { return heap.Pop(pq.h).(T) }

type hp[T any] struct {
	s       []T
	compare func(a, b T) int
}

func (h hp[T]) Len() int           { return len(h.s) }
func (h hp[T]) Swap(i, j int)      { h.s[i], h.s[j] = h.s[j], h.s[i] }
func (h hp[T]) Less(i, j int) bool { return h.compare(h.s[i], h.s[j]) < 0 }
func (h *hp[T]) Push(x any) {
	h.s = append(h.s, x.(T))
}
func (h *hp[T]) Pop() any {
	n := len(h.s)
	res := (h.s)[n-1]
	h.s = (h.s)[:n-1]
	return res
}
