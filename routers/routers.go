package routers

import (
	"Voichatter/api"
	"Voichatter/configs"
	"Voichatter/interceptor"
	"Voichatter/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func SetupRouter() *gin.Engine {
	if configs.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middleware.Cors())

	r.GET("/yy", handleWebSocket)

	v := r.Group("api")
	{
		v.POST("/register", api.UserRegister)
		v.POST("/login", api.UserLogin)
		v.GET("/ws", api.Ws)
		authed := v.Group("/") // 需要登陆保护
		authed.Use(interceptor.ConfInterceptor())
		{
			authed.POST("/logout", api.UserLogout)
			authed.GET("/servers-list", api.FindUserServersList)
			authed.GET("/history", api.FindMessage)
			authed.POST("/create-server", api.CreateServer)
		}
	}
	return r
}

var (
	clients   = make(map[string][]*websocket.Conn)
	clientsMu = &sync.RWMutex{}
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
)

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("=================Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// 从URL参数获取频道号
	clientID := c.Query("channelID")
	//clientID := c.Request.RemoteAddr
	clientsMu.Lock()
	clients[clientID] = append(clients[clientID], conn)
	clientsMu.Unlock()

	// 接收和处理消息
	for {
		messageType, p, err := conn.ReadMessage()
		log.Println("========================mp", clientID, messageType, string(p))
		if err != nil {
			log.Println("========================", err)
			break
		}
		if messageType == websocket.BinaryMessage {
			// 处理接收到的音频二进制数据
			// 将消息广播给房间内的所有客户端
			go broadcastMessage(clientID, conn, messageType, p)
		}
	}

	// 在连接关闭时，将其从频道中移除
	defer func() {
		clientsMu.Lock()
		connections := clients[clientID]
		for i, co := range connections {
			if co == conn {
				clients[clientID] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		clientsMu.Unlock()
	}()
}

func broadcastMessage(clientID string, sender *websocket.Conn, messageType int, message []byte) {
	// 查询房间内所有连接的客户端
	clientsMu.RLock()
	connections := clients[clientID]
	clientsMu.RUnlock()

	for _, conn := range connections {
		// 发送消息给除了发送者之外的所有客户端
		if conn != sender {
			// 发送消息给所有客户端包括自己
			//if err := conn.WriteMessage(messageType, message); err != nil {
			//	log.Println("====================Error writing message:", err)
			//	return
			//}
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Println("error:", err)
				conn.Close()
			}
		}
	}
}
