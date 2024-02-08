package service

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"online-voice-channel/models"
	"sync"
	"time"
)

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

func FindHistory(list *[]models.Message, r *http.Request, db *gorm.DB) error {
	return db.Where("channel_id = ?", r.URL.Query().Get("channelID")).Find(&list).Error
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

	// 从URL参数获取频道号
	channelID := r.URL.Query().Get("channelID")

	// 将连接加入到对应的房间
	roomsMutex.Lock()
	rooms[channelID] = append(rooms[channelID], conn)
	roomsMutex.Unlock()

	// 接收和处理消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// 持久化消息到数据库
		content := string(p)
		db.Create(&models.Message{Content: content, ChannelID: channelID, SendDate: time.Now()})

		// 将消息广播给房间内的所有客户端
		go broadcastMessage(channelID, conn, messageType, p)
	}

	// 在连接关闭时，将其从房间中移除
	roomsMutex.Lock()
	connections := rooms[channelID]
	for i, c := range connections {
		if c == conn {
			rooms[channelID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	roomsMutex.Unlock()
}

func broadcastMessage(channelID string, sender *websocket.Conn, messageType int, message []byte) {
	// 查询房间内所有连接的客户端
	roomsMutex.RLock()
	connections := rooms[channelID]
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
