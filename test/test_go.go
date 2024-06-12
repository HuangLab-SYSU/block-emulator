package test

import "fmt"

func F() {
	a := 1
	fmt.Printf("a=%v",a)
	go g()

}

func g() {
	for i:=0;i<1000;i++{
		fmt.Println(i)
	}
}