package main

import (
	"log"
	"time"
)




func main() {
	ch := make(chan struct{})
	go func ()  {
		<- time.After(time.Second)
		close(ch)
	}()

	val, ok := <- ch
	if !ok {
		log.Print("channel closed")
		return
	}

	log.Print(val)

}
