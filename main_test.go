package main

import (
	"sync"
	"testing"
)

func TestMerge2Channels(t *testing.T) {
	n := 100

	var x1Input []int
	var x2Input []int
	var expectedResults []int

	f := func(i int) int {
		return i * 2
	}

	for x := 0; x < n; x++ {
		x1Input = append(x1Input, x + 1)
		x2Input = append(x2Input, x + 2)
		expectedResults = append(expectedResults, f(x + 1) + f(x + 2))
	}

	in1 := make(chan int)
	in2 := make(chan int)

	out := make(chan int)

	Merge2Channels(f, in1, in2, out, n)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		for x := 0; x < n; x++ {
			in1 <- x1Input[x]
			in2 <- x2Input[x]
		}

		wg.Done()
	}()

	go func() {
		for x := 0; x < n; x++ {
			result := <- out

			gotValue := false

			for k, v := range expectedResults {
				if result == v {
					gotValue = true
					expectedResults = append(expectedResults[:k], expectedResults[k+1:]...)
				}
			}

			if !gotValue {
				t.Errorf("[%d] Result %d not found", x, result)
			}
		}

		wg.Done()
	}()

	wg.Wait()
}