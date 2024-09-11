package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	go processUsers()
	time.Sleep(3 * time.Second)
}

func processUsers() {
	i := 0
	ch := make(chan string, 30)
	go func() {
	}()
	for true {
		go processUser(i, ch)
		time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
		i++
	}
}

func processUser(userID int, ch chan string) {
	result := fmt.Sprintf("I'm processing user %d", userID)
	time.Sleep(200 * time.Millisecond)
	ch <- result
	fmt.Println("Se cargo el mensaje", result)
}
