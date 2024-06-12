package test

import (
	"fmt"
	"time"
)

func Goru() {
	for i := 0; i < 10; i++ {
		go wait(i)
		time.Sleep(1000 * time.Millisecond)

	}
}

func wait(i int) {
	time.Sleep(2000 * time.Millisecond)
	fmt.Println(i)
}