package informer

import (
	"ac/internal/info"
	"ac/internal/websocket"
	"encoding/json"
)

type Informer struct {
	Infos []info.Info
	hub   *websocket.Hub
}

func NewInformer(hub *websocket.Hub) *Informer {
	return &Informer{
		hub: hub,
	}
}

func (i *Informer) Get() []info.Info {
	return i.Infos
}

func (i *Informer) Add(info info.Info) {
	i.Infos = append(i.Infos, info)

	// 广播新添加的信息
	jsonData, _ := json.Marshal(map[string]any{
		"type": "informer",
		"data": info,
	})
	i.hub.Broadcast(jsonData)
}

func (i *Informer) Format(infos []info.Info) string {
	var formatted string
	for _, info := range infos {
		formatted += info.Format() + "\n"
	}
	return formatted
}

func (i *Informer) Clear() {
	i.Infos = []info.Info{}
}
