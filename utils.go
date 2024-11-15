package main

import (
	"slices"
	"sync"
)

type stack[T any] []T

func (s *stack[T]) push(item T) {
	*s = append(*s, item)
}
func (s *stack[T]) pop() T {
	n := len(*s) - 1
	item := (*s)[n]
	*s = (*s)[:n]
	return item
}
func (s *stack[T]) peek() T {
	n := len(*s) - 1
	return (*s)[n]
}

func (s *stack[T]) size() int {
	return len(*s)
}

func chanToSlice[T any](ch <-chan T) []T {
	var wg sync.WaitGroup

	result := make([]T, 0)
	wg.Add(1)
	go func(cha <-chan T) {
		for v := range cha {
			result = append(result, v)
		}
		wg.Done()
	}(ch)
	wg.Wait()
	return result
}
func chanToSortedSlice[T any](ch <-chan T, sortFunc func(a, b T) int) []T {
	slice := chanToSlice(ch)
	slices.SortFunc(slice, sortFunc)
	return slice
}
