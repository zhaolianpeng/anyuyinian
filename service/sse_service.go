package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// SSEManager 管理SSE连接
type SSEManager struct {
	clients    map[chan string]bool
	broadcast  chan []byte
	register   chan chan string
	unregister chan chan string
	mutex      sync.RWMutex
}

var SSEManagerInstance = &SSEManager{
	clients:    make(map[chan string]bool),
	broadcast:  make(chan []byte),
	register:   make(chan chan string),
	unregister: make(chan chan string),
}

func (manager *SSEManager) Start() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()
			log.Println("SSE客户端已连接")

		case client := <-manager.unregister:
			manager.mutex.Lock()
			delete(manager.clients, client)
			close(client)
			manager.mutex.Unlock()
			log.Println("SSE客户端已断开")

		case message := <-manager.broadcast:
			manager.mutex.RLock()
			for client := range manager.clients {
				select {
				case client <- string(message):
				default:
					// 客户端可能已断开，移除它
					delete(manager.clients, client)
					close(client)
				}
			}
			manager.mutex.RUnlock()
		}
	}
}

// SSEHandler 处理SSE连接
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 创建客户端通道
	clientChan := make(chan string)
	SSEManagerInstance.register <- clientChan

	// 确保在函数结束时清理客户端
	defer func() {
		SSEManagerInstance.unregister <- clientChan
	}()

	// 发送连接成功消息
	welcomeMsg := map[string]interface{}{
		"type":      "system",
		"message":   "SSE连接成功",
		"timestamp": time.Now().Unix(),
	}
	welcomeBytes, _ := json.Marshal(welcomeMsg)
	fmt.Fprintf(w, "data: %s\n\n", welcomeBytes)
	w.(http.Flusher).Flush()

	// 监听客户端断开
	notify := w.(http.CloseNotifier).CloseNotify()

	// 发送心跳保持连接
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			w.(http.Flusher).Flush()

		case <-ticker.C:
			// 发送心跳
			heartbeat := map[string]interface{}{
				"type":      "heartbeat",
				"timestamp": time.Now().Unix(),
			}
			heartbeatBytes, _ := json.Marshal(heartbeat)
			fmt.Fprintf(w, "data: %s\n\n", heartbeatBytes)
			w.(http.Flusher).Flush()

		case <-notify:
			log.Println("SSE客户端断开连接")
			return
		}
	}
}

// BroadcastSSEMessage 广播SSE消息
func BroadcastSSEMessage(messageType string, data interface{}) {
	message := map[string]interface{}{
		"type":      messageType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
	messageBytes, _ := json.Marshal(message)
	SSEManagerInstance.broadcast <- messageBytes
}

// BroadcastSSEOrderUpdate 广播SSE订单更新
func BroadcastSSEOrderUpdate(orderID string, status string, amount float64) {
	update := map[string]interface{}{
		"orderID": orderID,
		"status":  status,
		"amount":  amount,
	}
	BroadcastSSEMessage("orderUpdate", update)
}

// BroadcastSSENotification 广播SSE通知
func BroadcastSSENotification(title string, content string, notificationType string) {
	notification := map[string]interface{}{
		"title": title,
		"content": content,
		"type":   notificationType,
	}
	BroadcastSSEMessage("notification", notification)
}

// BroadcastSSESystemMessage 广播SSE系统消息
func BroadcastSSESystemMessage(messageType string, content string) {
	systemMsg := map[string]interface{}{
		"messageType": messageType,
		"content":     content,
	}
	BroadcastSSEMessage("systemMessage", systemMsg)
}

// GetConnectedSSEClientsCount 获取连接的SSE客户端数量
func GetConnectedSSEClientsCount() int {
	SSEManagerInstance.mutex.RLock()
	defer SSEManagerInstance.mutex.RUnlock()
	return len(SSEManagerInstance.clients)
} 