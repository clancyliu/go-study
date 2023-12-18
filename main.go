package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "go-study/web_socket"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var GM = make(map[int64]*ws.ChatManager)

func main() {

	gin.SetMode(gin.ReleaseMode) //线上环境

	go ws.CM.Start()

	r := gin.Default()
	r.GET("/chat/:userId", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		userId, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			fmt.Println(err)
			return
		}

		chatClient := &ws.ChatClient{
			UserId: int64(userId),
			Conn:   conn,
			Send:   make(chan []byte),
		}
		ws.CM.RegisterChan <- chatClient
		go chatClient.ReadMessage()
		go chatClient.WriteMessage()
	})

	r.GET("/group/:groupId/:userId", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		groupId, err := strconv.Atoi(c.Param("groupId"))
		if err != nil {
			fmt.Println(err)
			return
		}
		userId, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			fmt.Println(err)
			return
		}

		chatClient := &ws.ChatClient{
			UserId: int64(userId),
			Conn:   conn,
			Send:   make(chan []byte),
		}
		if _, ok := GM[int64(groupId)]; !ok {
			GM[int64(groupId)] = &ws.ChatManager{
				Clients:        make(map[int64]*ws.ChatClient),
				RegisterChan:   make(chan *ws.ChatClient),
				UnRegisterChan: make(chan *ws.ChatClient),
				Broadcast:      make(chan []byte),
			}
			go GM[int64(groupId)].Start()
		}

		GM[int64(groupId)].RegisterChan <- chatClient
		go chatClient.ReadMessage()
		go chatClient.WriteMessage()
	})

	r.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8282") // listen and serve on 0.0.0.0:8080
}
