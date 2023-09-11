package main

import (
	"fmt"
	_ "net/http/pprof"
)

func main() {
	//go func() {
	//	for {
	//		log.Println("hello")
	//		time.Sleep(2 * time.Second)
	//	}
	//}()
	//
	//err := http.ListenAndServe("0.0.0.0:6060", nil)
	//if err != nil {
	//	return
	//}

	//fmt.Println(findN(2))

	fmt.Println("hello")

}

//
//func findN(n int) int {
//	if n < 2 {
//		return 1
//	}
//	pre1, pre2 := 1, 1
//	for i := 2; i <= n; i++ {
//		pre2, pre1 = pre2+pre1, pre2
//	}
//	return pre2
//}
