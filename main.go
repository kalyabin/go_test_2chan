package main

type fxResult struct {
	num int
	x int
}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	processed := 0
	processFx := func(item fxResult, resultChan chan fxResult) {
		item.x = f(item.x)
		processed++
		resultChan <- item
	}

	go func() {
		results := make([][]int, n)
		resultChan := make(chan fxResult)

		processed := 0
		ch1 := 0
		ch2 := 0

		for {
			select {
			case x := <- in1:
				if ch1 < n {
					item := fxResult{ch1, x}
					ch1++
					go processFx(item, resultChan)
				}
			case x := <- in2:
				if ch2 < n {
					item := fxResult{ch2, x}
					ch2++
					go processFx(item, resultChan)
				}
			case item := <- resultChan:
				results[item.num] = append(results[item.num], item.x)

				if len(results[item.num]) == 2 {
					result := results[item.num][0] + results[item.num][1]
					out <- result
					processed++
				}

				if processed == n {
					return
				}
			}
		}
	}()
}