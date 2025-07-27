package websocket

import (
	"encoding/json"
	"log"
)

// 消息类型枚举
const (
	MessageTypeAIResponse   = "ai_response"
	MessageTypeConnection   = "connection_status"
	MessageTypeRealTimeData = "real_time_data"
	MessageTypeNotification = "notification"
)

// 统一消息结构
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// 广播器接口
type Broadcaster interface {
	Broadcast(msgType string, data interface{})
}

// HubBroadcaster 实现
type HubBroadcaster struct {
	hub *Hub
}

func NewHubBroadcaster(hub *Hub) *HubBroadcaster {
	return &HubBroadcaster{hub: hub}
}

func (b *HubBroadcaster) Broadcast(msgType string, data interface{}) {
	msg := Message{
		Type: msgType,
		Data: data,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("JSON编码失败: %v", err)
		return
	}

	b.hub.Broadcast(jsonData)
}
