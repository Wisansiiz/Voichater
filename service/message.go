package service

import (
	"Voichatter/dao"
	"Voichatter/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func FindHistory(list *[]models.Message, channelID string) error {
	return dao.DB.Where("channel_id = ?", channelID).Find(&list).Error
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(conn)

	// 从URL参数获取频道号
	channelID := c.Query("channelID")

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
		userId, _ := c.Get("user_id")
		dao.DB.Create(
			&models.Message{
				SenderUserID: userId.(uint),
				Content:      content,
				ChannelID:    channelID,
				SendDate:     time.Now(),
			},
		)
		// 将消息广播给房间内的所有客户端
		go broadcastMessage(channelID, messageType, p)
	}

	// 在连接关闭时，将其从房间中移除
	roomsMutex.Lock()
	connections := rooms[channelID]
	for i, number := range connections {
		if number == conn {
			rooms[channelID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	roomsMutex.Unlock()
}

func broadcastMessage(channelID string, messageType int, message []byte) {
	// 查询房间内所有连接的客户端
	roomsMutex.RLock()
	connections := rooms[channelID]
	roomsMutex.RUnlock()

	for _, conn := range connections {
		// 发送消息给所有客户端包括自己
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}
