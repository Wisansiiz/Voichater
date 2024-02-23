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
		"code":     200,
		"messages": "success",
	})
}

func UserLogin(c *gin.Context) {
	var user models.UserLoginResponse
	_ = c.ShouldBind(&user)
	token, err := service.UserLogin(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":     400,
			"messages": "账号或密码错误",
			"data":     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "success",
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
		"code":     200,
		"messages": "success",
		"data":     server,
	})
}

func UserLogout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	token := parts[1]
	if err := service.UserLogout(token); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":     400,
			"messages": "退出登录发生错误",
			"data":     "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "退出登录成功",
		"data":     "",
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
			"code":     400,
			"messages": err.Error(),
			"data":     "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "创建成功",
		"data":     "",
	})
}

func JoinServer(c *gin.Context) {
	userId := c.MustGet("user_id")
	var member models.Member
	member.UserID = userId.(uint)
	_ = c.ShouldBind(&member)
	if err := service.JoinServer(&member); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":     400,
			"messages": err.Error(),
			"data":     "",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"messages": "加入成功",
		"data":     "",
	})
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
	clients       = make(map[*websocket.Conn]string)
	groupChannels = make(map[string][]*websocket.Conn)
	clientsMu     = &sync.RWMutex{}
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
		log.Println("conn.Close()")
		if err != nil {
			log.Println("conn.Close():", err)
			return
		}
	}(conn)

	channelId := c.Query("channelId")
	currentUserID := c.Query("id")
	clientsMu.Lock()
	clients[conn] = currentUserID
	clientsMu.Unlock()
	// 接收和处理消息
	for {
		var msg Msg
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("err:", err)
			break
		}
		if err = json.Unmarshal(p, &msg); err != nil {
			return
		}
		if msg.Code == "offer" {
			targetId := msg.Data["targetId"]
			offer := msg.Data["offer"]
			broadcastMessage(targetId, currentUserID, "offer", "offer", offer)
		} else if msg.Code == "answer" {
			targetId := msg.Data["targetId"]
			answer := msg.Data["answer"]
			broadcastMessage(targetId, currentUserID, "answer", "answer", answer)
		} else if msg.Code == "icecandidate" {
			targetId := msg.Data["targetId"]
			candidate := msg.Data["candidate"]
			broadcastMessage(targetId, currentUserID, "icecandidate", "candidate", candidate)
		} else if msg.Code == "join_group" {
			clientsMu.Lock()
			groupChannels[channelId] = append(groupChannels[channelId], conn)
			clientsMu.Unlock()

			if broadcastGroups(msg.Code, channelId, conn, currentUserID) {
				return
			}
		} else if msg.Code == "leave_group" {
			if broadcastGroups(msg.Code, channelId, conn, currentUserID) {
				return
			}
			break
		}
	}

	if channelId != "" {
		// 在连接关闭时，将其从房间中移除
		clientsMu.Lock()
		connections := groupChannels[channelId]
		for i, numbers := range connections {
			if numbers == conn {
				groupChannels[channelId] = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		clientsMu.Unlock()
	}

	clientsMu.Lock()
	// 在连接关闭时，将其从连接中移除
	for clientConn := range clients {
		if clientConn == conn {
			delete(clients, clientConn)
		}
	}
	clientsMu.Unlock()
}

func broadcastGroups(code string, channelId string, conn *websocket.Conn, currentUserID string) bool {
	clientsMu.RLock()
	connections := groupChannels[channelId]
	clientsMu.RUnlock()
	for _, numbers := range connections {
		if conn != numbers {
			message := gin.H{
				"code": code,
				"data": gin.H{
					"fromId": currentUserID,
				},
			}
			jsonBytes, _ := json.Marshal(message)
			if err := numbers.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
				log.Println("Error writing messages:", err)
				return true
			}
		}
	}
	return false
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
				log.Println("Error writing messages:", err)
				return
			}
		}
	}
}
