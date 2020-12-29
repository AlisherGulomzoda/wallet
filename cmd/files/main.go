package main

import "log"




func main() {
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}

	ch := make(chan int)
	defer close(ch)

	parts := 10
	size := len(data) / parts
	channels := make([] <- chan int, parts)

	for i := 0; i < parts; i++ {
		ch := make(chan int)
		channels[i] = ch
		go func (ch chan <- int, data []int)  {
			sum := 0
			for _, v := range data {
				sum += v
			}
			ch <- sum
		}(ch, data[i * size : (i +1) * size])
	}

	total := 0
	for value := range channels {
		total += value
	}

	log.Print(total)
}
