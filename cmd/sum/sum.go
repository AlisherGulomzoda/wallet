package sum

import "sync"

func Regular() int64 {
	sum := int64(0)
	for i := 0; i < 4_000_000; i++ {
		sum++
	}
	return sum
}

func Concurently() int64 {
	wg := sync.WaitGroup{}
	wg.Add(2)

	mu := sync.Mutex{}
	sum := int64(0)
	
	go func() {
		defer wg.Done()
		val := int64(0)
		for i := 0; i < 2_000_000; i++ {
			val++
		}
		mu.Lock()
		sum += val
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		val := int64(0)
		for i := 0; i < 2_000_000; i++ {
			val++
		}
		mu.Lock()
		sum += val
		mu.Unlock()
	}()

	wg.Wait()
	return sum
}