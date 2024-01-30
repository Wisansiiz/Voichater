package models

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

// Message 消息模型
type Message struct {
	gorm.Model
	Text   string `json:"text"`
	RoomID string `json:"room_id"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
	rooms      = make(map[string][]*websocket.Conn)
	roomsMutex = &sync.RWMutex{}
)

func FindHistory(list *[]Message, r *http.Request, db *gorm.DB) {
	db.Where("room_id = ?", r.URL.Query().Get("roomID")).Find(&list)
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
			return
		}
	}(conn)

	// 从URL参数获取房间号
	roomID := r.URL.Query().Get("roomID")

	// 将连接加入到对应的房间
	roomsMutex.Lock()
	rooms[roomID] = append(rooms[roomID], conn)
	roomsMutex.Unlock()

	// 接收和处理消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// 持久化消息到数据库
		text := string(p)
		db.Create(&Message{Text: text, RoomID: roomID})

		// 将消息广播给房间内的所有客户端
		go broadcastMessage(roomID, conn, messageType, p)
	}

	// 在连接关闭时，将其从房间中移除
	roomsMutex.Lock()
	connections := rooms[roomID]
	for i, c := range connections {
		if c == conn {
			rooms[roomID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	roomsMutex.Unlock()
}

func broadcastMessage(roomID string, sender *websocket.Conn, messageType int, message []byte) {
	// 查询房间内所有连接的客户端
	roomsMutex.RLock()
	connections := rooms[roomID]
	roomsMutex.RUnlock()

	for _, conn := range connections {
		// 发送消息给除了发送者之外的所有客户端
		//if conn != sender {
		// 发送消息给所有客户端包括自己
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Error writing message:", err)
			return
		}
		//}
	}
}
