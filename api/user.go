package api

import (
	"Voichatter/models"
	"Voichatter/service"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

func UserRegister(c *gin.Context) {
	var user models.User
	_ = c.ShouldBind(&user)
	if err := service.UserRegister(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func UserLogin(c *gin.Context) {
	var user models.UserLoginResponse
	_ = c.ShouldBind(&user)
	token, err := service.UserLogin(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "账号或密码错误",
			"data": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"token": token,
		},
	})
}

func FindUserServersList(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var user models.User
	user.UserID = userId.(uint)
	var server []models.Server
	if err := service.FindUserServersList(&user, &server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": server,
	})
}

func UserLogout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := parts[1]
	if err := service.UserLogout(token); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "退出登录发生错误",
			"data": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出登录成功",
		"data": "",
	})
}

func CreateServer(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var user models.User
	user.UserID = userId.(uint)
	var server models.Server
	_ = c.ShouldBind(&server)
	if err := service.CreateServer(&user, &server); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  err,
			"data": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": "",
	})
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
	clients   = make(map[*websocket.Conn]string)
	clientsMu = &sync.RWMutex{}
)

func HandleWebSocket(c *gin.Context) {
	type Msg struct {
		Code string         `json:"code"`
		Data map[string]any `json:"data"`
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
			return
		}
	}(conn)

	currentUserID := c.Query("id")
	clientsMu.Lock()
	clients[conn] = currentUserID
	clientsMu.Unlock()
	// 接收和处理消息
	for {
		var message Msg
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if err = json.Unmarshal(p, &message); err != nil {
			return
		}
		if message.Code == "offer" {
			targetId := message.Data["targetId"]
			offer := message.Data["offer"]
			broadcastMessage(targetId, currentUserID, "offer", "offer", offer)
		} else if message.Code == "answer" {
			targetId := message.Data["targetId"]
			answer := message.Data["answer"]
			broadcastMessage(targetId, currentUserID, "answer", "answer", answer)
		} else if message.Code == "icecandidate" {
			targetId := message.Data["targetId"]
			candidate := message.Data["candidate"]
			broadcastMessage(targetId, currentUserID, "icecandidate", "candidate", candidate)
		}
	}
	// 在连接关闭时，将其从房间中移除
	clientsMu.Lock()
	for clientConn := range clients {
		if clientConn == conn {
			delete(clients, clientConn)
		}
	}
	clientsMu.Unlock()
}

func broadcastMessage(targetId any, currentUserID string, code string, dataName string, data any) {
	for clientConn, userId := range clients {
		if userId == targetId.(string) {
			message := gin.H{
				"code": code,
				"data": gin.H{
					"fromId": currentUserID,
					dataName: data,
				},
			}
			jsonBytes, _ := json.Marshal(message)
			if err := clientConn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}
}
