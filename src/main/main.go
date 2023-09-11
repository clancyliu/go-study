package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Age      *int      `form:"age"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	route := gin.Default()
	route.GET("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBind(&person) == nil {
		fmt.Println(fmt.Sprintf("age: %d", *person.Age))
	}

	c.String(200, "Success")
}
