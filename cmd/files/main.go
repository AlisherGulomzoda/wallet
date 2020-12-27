package main

import (
	"log"
	"sync"
)



func main() {
	log.Print("Main started!")
	wg := sync.WaitGroup{}
	wg.Add(2)

	sum := 0
	go func ()  {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			sum++
		}
	}()
	go func ()  {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			sum++
		}
	}()

	wg.Wait()
	log.Print("Main finished!")
	log.Print(sum)

}
