package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type ChatMessage struct {
	SenderId   int64  `json:"sender_id"`
	ReceiverId int64  `json:"receiver_id"`
	GroupId    int64  `json:"group_id"`
	Content    string `json:"content"`
}

type ChatClient struct {
	UserId int64
	Conn   *websocket.Conn
	Send   chan []byte
}

type ChatManager struct {
	Clients        map[int64]*ChatClient
	RegisterChan   chan *ChatClient
	UnRegisterChan chan *ChatClient
	Broadcast      chan []byte
}

var CM = &ChatManager{
	Clients:        make(map[int64]*ChatClient),
	RegisterChan:   make(chan *ChatClient),
	UnRegisterChan: make(chan *ChatClient),
	Broadcast:      make(chan []byte),
}

func (manager *ChatManager) Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case client := <-CM.RegisterChan:
			log.Printf("新用户加入:%v", client.UserId)
			CM.Clients[client.UserId] = client
			jsonMessage, _ := json.Marshal(&ChatMessage{Content: "Successful connection to socket service"})
			client.Send <- jsonMessage
		case conn := <-CM.UnRegisterChan:
			log.Printf("用户离开:%v", conn.UserId)
			if _, ok := CM.Clients[conn.UserId]; ok {
				jsonMessage, _ := json.Marshal(&ChatMessage{Content: "A socket has disconnected"})
				conn.Send <- jsonMessage
				close(conn.Send)
				delete(CM.Clients, conn.UserId)
			}
		case message := <-CM.Broadcast:
			var m ChatMessage
			if err := json.Unmarshal(message, &m); err != nil {
				fmt.Println(err)
			}
			for id, conn := range CM.Clients {
				// 群聊消息
				if m.GroupId > 0 && id != m.SenderId {
					conn.Send <- message
					continue
				}
				// 一对一
				if id == m.ReceiverId {
					conn.Send <- message
					break
				}
			}
		}
	}
}

func (c *ChatClient) ReadMessage() {
	defer func() {
		delete(CM.Clients, c.UserId)
		c.Conn.Close()
	}()
	for {
		c.Conn.PongHandler()
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			CM.UnRegisterChan <- c
			c.Conn.Close()
			break
		}
		log.Printf("读取到客户端的信息:%s", string(message))
		CM.Broadcast <- message
	}
}

func (c *ChatClient) WriteMessage() {
	defer func() {
		delete(CM.Clients, c.UserId)
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return
				}
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}
