package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		log.Println("hello")
	})
	err := r.Run(":9080")
	if err != nil {
		log.Fatalln("run server error")
		return
	}
}
