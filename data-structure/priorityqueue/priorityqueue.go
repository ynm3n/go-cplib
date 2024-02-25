package priorityqueue

import (
	"cmp"
	"container/heap"
)

type PriorityQueue[T any] interface {
	Len() int
	Top() T
	Enqueue(x T)
	Dequeue() T
}

// 昇順
func NewPriorityQueue[T cmp.Ordered]() PriorityQueue[T] {
	pq := NewPriorityQueueFunc(cmp.Less[T])
	return pq
}

// lessは並んでる時にtrue、並んでない時にfalseを返すような関数
// 参考: https://pkg.go.dev/cmp@go1.22.0#Less
func NewPriorityQueueFunc[T any](less func(a, b T) bool) PriorityQueue[T] {
	pq := new(priorityQueue[T])
	pq.s = make([]T, 0)
	pq.less = less
	return pq
}

type priorityQueue[T any] struct {
	s    []T
	less func(a, b T) bool
}

func (pq *priorityQueue[T]) Top() T      { return pq.s[0] }
func (pq *priorityQueue[T]) Enqueue(x T) { heap.Push(pq, x) }
func (pq *priorityQueue[T]) Dequeue() T  { return heap.Pop(pq).(T) }

func (pq priorityQueue[T]) Len() int           { return len(pq.s) }
func (pq priorityQueue[T]) Swap(i, j int)      { pq.s[i], pq.s[j] = pq.s[j], pq.s[i] }
func (pq priorityQueue[T]) Less(i, j int) bool { return pq.less(pq.s[i], pq.s[j]) }
func (pq *priorityQueue[T]) Push(x any) {
	pq.s = append(pq.s, x.(T))
}
func (pq *priorityQueue[T]) Pop() any {
	n := len(pq.s)
	res := (pq.s)[n-1]
	pq.s = (pq.s)[:n-1]
	return res
}
