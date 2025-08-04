package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，微信云托管会自动处理
	},
}

type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

var WsManager = &WebSocketManager{
	clients:    make(map[*websocket.Conn]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
}

func (manager *WebSocketManager) Start() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()
			log.Println("WebSocket客户端已连接")

		case client := <-manager.unregister:
			manager.mutex.Lock()
			delete(manager.clients, client)
			manager.mutex.Unlock()
			log.Println("WebSocket客户端已断开")

		case message := <-manager.broadcast:
			manager.mutex.RLock()
			for client := range manager.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("发送消息失败: %v", err)
					client.Close()
					delete(manager.clients, client)
				}
			}
			manager.mutex.RUnlock()
		}
	}
}

// WebSocketHandler 处理微信云托管的WebSocket连接
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	WsManager.register <- conn

	// 发送欢迎消息
	welcomeMsg := map[string]interface{}{
		"type":      "system",
		"message":   "WebSocket连接成功",
		"timestamp": time.Now().Unix(),
	}
	welcomeBytes, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, welcomeBytes)

	// 处理消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			WsManager.unregister <- conn
			conn.Close()
			break
		}

		// 处理接收到的消息
		handleWebSocketMessage(conn, message)
	}
}

func handleWebSocketMessage(conn *websocket.Conn, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("解析消息失败: %v", err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "connect":
		// 处理连接消息
		response := map[string]interface{}{
			"type":      "connect_ack",
			"message":   "连接已确认",
			"timestamp": time.Now().Unix(),
		}
		responseBytes, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responseBytes)

	case "ping":
		// 处理心跳消息
		response := map[string]interface{}{
			"type":      "pong",
			"timestamp": time.Now().Unix(),
		}
		responseBytes, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, responseBytes)

	case "order_update":
		// 处理订单更新消息
		handleOrderUpdateMessage(conn, msg)

	case "notification":
		// 处理通知消息
		handleNotificationMessage(conn, msg)

	default:
		log.Printf("未知消息类型: %s", msgType)
	}
}

func handleOrderUpdateMessage(conn *websocket.Conn, msg map[string]interface{}) {
	// 处理订单更新逻辑
	response := map[string]interface{}{
		"type":      "order_update_ack",
		"message":   "订单更新已确认",
		"timestamp": time.Now().Unix(),
	}
	responseBytes, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseBytes)
}

func handleNotificationMessage(conn *websocket.Conn, msg map[string]interface{}) {
	// 处理通知消息逻辑
	response := map[string]interface{}{
		"type":      "notification_ack",
		"message":   "通知已确认",
		"timestamp": time.Now().Unix(),
	}
	responseBytes, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, responseBytes)
}

// BroadcastMessage 广播消息给所有客户端
func BroadcastMessage(messageType string, data interface{}) {
	message := map[string]interface{}{
		"type":      messageType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
	messageBytes, _ := json.Marshal(message)
	WsManager.broadcast <- messageBytes
}

// BroadcastOrderUpdate 广播订单更新
func BroadcastOrderUpdate(orderID string, status string, amount float64) {
	message := map[string]interface{}{
		"type": "order_update",
		"data": map[string]interface{}{
			"orderId": orderID,
			"status":  status,
			"amount":  amount,
		},
		"timestamp": time.Now().Unix(),
	}
	messageBytes, _ := json.Marshal(message)
	WsManager.broadcast <- messageBytes
}

// BroadcastNotification 广播通知消息
func BroadcastNotification(title string, content string, notificationType string) {
	message := map[string]interface{}{
		"type": "notification",
		"data": map[string]interface{}{
			"title":   title,
			"content": content,
			"type":    notificationType,
		},
		"timestamp": time.Now().Unix(),
	}
	messageBytes, _ := json.Marshal(message)
	WsManager.broadcast <- messageBytes
}

// BroadcastSystemMessage 广播系统消息
func BroadcastSystemMessage(messageType string, content string) {
	message := map[string]interface{}{
		"type": "system_message",
		"data": map[string]interface{}{
			"type":    messageType,
			"content": content,
		},
		"timestamp": time.Now().Unix(),
	}
	messageBytes, _ := json.Marshal(message)
	WsManager.broadcast <- messageBytes
}

// GetConnectedClientsCount 获取连接的客户端数量
func GetConnectedClientsCount() int {
	WsManager.mutex.RLock()
	defer WsManager.mutex.RUnlock()
	return len(WsManager.clients)
}
