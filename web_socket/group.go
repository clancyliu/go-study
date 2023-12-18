package ws

//
//type GroupMessage struct {
//	GroupId  int64  `json:"group_id"`
//	SenderId int64  `json:"sender_id"`
//	Content  string `json:"content"`
//}
//
//type GroupClient struct {
//	UserId  int64
//	GroupId int64
//	Conn    *websocket.Conn
//	Send    chan []byte
//}
//
//type GroupManager struct {
//	Clients        map[int64]*GroupClient
//	RegisterChan   chan *GroupClient
//	UnRegisterChan chan *GroupClient
//	Broadcast      chan []byte
//}

//
//func (manager *GroupManager) Start() {
//	for {
//		log.Println("<---管道通信--->")
//		select {
//		case client := <-manager.RegisterChan:
//			log.Printf("新用户加入:%v", client.UserId)
//			manager.Clients[client.UserId] = client
//			jsonMessage, _ := json.Marshal(&GroupMessage{Content: "Successful connection to socket service"})
//			client.Send <- jsonMessage
//		case conn := <-manager.UnRegisterChan:
//			log.Printf("用户离开:%v", conn.UserId)
//			if _, ok := manager.Clients[conn.UserId]; ok {
//				jsonMessage, _ := json.Marshal(&ChatMessage{Content: "A socket has disconnected"})
//				conn.Send <- jsonMessage
//				close(conn.Send)
//				delete(CM.Clients, conn.UserId)
//			}
//		case message := <-manager.Broadcast:
//			var m ChatMessage
//			json.Unmarshal(message, &m)
//			for userId, conn := range manager.Clients {
//				if userId != m.SenderId {
//					conn.Send <- message
//				}
//			}
//		}
//	}
//}

//
//func (c *GroupClient) ReadMessage() {
//	defer func() {
//		delete(GM[c.GroupId].Clients, c.UserId)
//		c.Conn.Close()
//	}()
//	for {
//		c.Conn.PongHandler()
//		_, message, err := c.Conn.ReadMessage()
//		if err != nil {
//			GM[c.GroupId].UnRegisterChan <- c
//			c.Conn.Close()
//			break
//		}
//		log.Printf("读取到客户端的信息:%s", string(message))
//		GM[c.GroupId].Broadcast <- message
//	}
//}
//
//func (c *GroupClient) WriteMessage() {
//	defer func() {
//		delete(GM[c.GroupId].Clients, c.UserId)
//		c.Conn.Close()
//	}()
//	for {
//		select {
//		case message, ok := <-c.Send:
//			if !ok {
//				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
//					return
//				}
//				return
//			}
//			log.Printf("发送到到客户端的信息:%s", string(message))
//			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
//				return
//			}
//		}
//	}
//}
