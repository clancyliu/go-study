package main

import (
	"log"
	"os"
)

func main() {
	var logger = log.New(os.Stdout, "", log.Lmsgprefix)
	logger.Println("hello logger")

	file, err := os.OpenFile("sys.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Faild to open error logger file:", err)
		return
	}
	log.SetOutput(file)
	log.Println("hello log")

}
