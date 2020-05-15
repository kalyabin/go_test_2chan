package main

import (
	"sync"
	"testing"
	"time"
)

func TestMerge2Channels(t *testing.T) {
	n := 100

	gotResultCount := 0
	var x1Input []int
	var x2Input []int
	var expectedResults []int

	f := func(i int) int {
		time.Sleep(time.Millisecond * 1)
		return i * i
	}

	for x := 0; x < n; x++ {
		x1 := x + 123
		x2 := x + 456

		x1Input = append(x1Input, x1)
		x2Input = append(x2Input, x2)

		expectedResults = append(expectedResults, f(x1) + f(x2))
	}

	t.Logf("Wait Result count: %d\n", len(expectedResults))

	in1 := make(chan int, 1)
	in2 := make(chan int, 1)

	out := make(chan int)

	Merge2Channels(f, in1, in2, out, n)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		for x := 0; x < n; x++ {
			in1 <- x1Input[x]
			in2 <- x2Input[x]
		}

		in1 <- 10
		in2 <- 20

		wg.Done()
	}()

	go func() {
		for x := 0; x < n; x++ {
			result := <- out
			gotValue := false
			gotResultCount++

			for k, v := range expectedResults {
				if result == v {
					gotValue = true
					expectedResults = append(expectedResults[:k], expectedResults[k+1:]...)
					t.Logf("Got expected result: %d", result)
				}
			}

			if !gotValue {
				t.Errorf("[%d] Result %d not found", n, result)
			}
		}

		wg.Done()
	}()

	wg.Wait()

	t.Logf("Got result count: %d", gotResultCount)
	t.Logf("Got expected results: %d", len(expectedResults))

	close(in1)
	close(in2)
	close(out)
}
