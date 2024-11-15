package main

import (
	"testing"
)

func TestStackPushPop(t *testing.T) {
	s := make(stack[int], 0)
	n := 5
	for i := 0; i <= n; i++ {
		s.push(i)
		if s.size() != i+1 {
			t.Fatalf("stack size expected %d, got %d", i+1, len(s))
		}
	}
	for j := n; j >= 0; j-- {
		actual := s.pop()
		if s.size() != j {
			t.Fatalf("stack size expected %d, got %d", j, len(s))
		}
		if actual != j {
			t.Fatalf("pop expected %d, got %d", j, actual)
		}
	}
}

func TestStackPushPeek(t *testing.T) {
	s := make(stack[int], 0)
	n := 5
	for i := 0; i <= n; i++ {
		s.push(i)
		if s.size() != i+1 {
			t.Fatalf("stack size expected %d, got %d", i+1, len(s))
		}
	}
	for j := n; j >= 0; j-- {
		actual := s.peek()
		if s.size() != n+1 {
			t.Fatalf("stack size expected %d, got %d", j, len(s))
		}
		if actual != n {
			t.Fatalf("peek expected %d, got %d", j, actual)
		}
	}
}

func TestChanToSlice(t *testing.T) {
	expectedLen := 10
	ch := make(chan int, 1)
	go func() {
		defer close(ch)
		for i := 0; i < expectedLen; i++ {
			ch <- i
		}
	}()
	result := chanToSlice(ch)
	if len(result) != expectedLen {
		t.Fatalf("Expected %d elements, got %d", expectedLen, len(result))
	}
}
func TestChanToSortedSlice(t *testing.T) {
	vals := []int{10, 5, 6, 3, 1, 2, 4, 9, 8, 7}
	expectedLen := len(vals)
	ch := make(chan int, 1)
	go func() {
		defer close(ch)
		for i := 0; i < expectedLen; i++ {
			ch <- vals[i]
		}
	}()
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := chanToSortedSlice(ch, func(a, b int) int {
		return a - b
	})
	if len(result) != expectedLen {
		t.Fatalf("Expected %d elements, got %d", expectedLen, len(result))
	}
	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {

			t.Fatalf("Expected %v, got %v", expected, result)
		}
	}
}
