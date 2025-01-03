package main

import (
	"fmt"
	"unsafe"
)

type Hello struct {
	len   int
	cap   int
	array *[5]int
}

//func main() {
//	s := [5]int{1, 2, 3}
//	h := Hello{
//		len:   3,
//		cap:   5,
//		array: &s,
//	}
//	appendHello(h)
//
//	fmt.Println(s)
//}

func appendSlice(s []int) {
	fmt.Printf("%p\n", &s)

	h := append(s, 4)                     // 可能触发底层数组重新分配
	fmt.Println("Inside appendSlice:", s) // 输出: 新的 slice

	fmt.Printf("%p\n", &h)
	fmt.Printf("%p\n", &s)

	s = h

	fmt.Println("Inside appendSlice:", s)
}

func main() {

	var h int
	fmt.Printf("%p\n", &h)

	s1 := make([]int, 4, 6) //[0 0 0 ]
	s2 := s1[1:3]           //[0 0]
	s1[1] = 10

	//fmt.Println(s1)
	//fmt.Println(len(s2), cap(s2))

	fmt.Println("s1=", s1, unsafe.Pointer(&s1[0]))
	fmt.Println("s2=", s2, unsafe.Pointer(&s2[0]))

	s2 = append(s2, 2)

	fmt.Println("s1=", s1, unsafe.Pointer(&s1[0]))
	fmt.Println("s2=", s2, unsafe.Pointer(&s2[0]))
}

func appendHello(h Hello) {
	h.len++
	h.array[3] = 4

	fmt.Println(h.array)
}

func changeValue(s []int) {
	s[0] = 100

	printAddr(s)
}

func changeSize(s []int) {
	s = append(s, 10)

	fmt.Println(s, len(s), cap(s))
	printAddr(s)
}

func printAddr(s []int) {
	if len(s) < 2 {
		return
	}

	fmt.Printf("s value: %v, s addr: %p, s[0]: %p, s[1]: %p\n", s, &s, &s[0], &s[1])
}
