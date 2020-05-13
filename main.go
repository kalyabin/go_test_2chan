package main

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	results := make([]chan int, 2)

	results[0] = make(chan int, 1)
	results[1] = make(chan int, 1)

	for x := 0; x < n; x++ {
		go func(x int) {
			r1 := <- results[0]
			r2 := <- results[1]

			out <- r1 + r2
		}(x)
	}

	go func() {
		for {
			select {
			case x := <-in1:
				go func(x int) {
					results[0] <- f(x)
				}(x)
			case x := <-in2:
				go func(x int) {
					results[1] <- f(x)
				}(x)
			}
		}
	}()
}