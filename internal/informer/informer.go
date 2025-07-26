package informer

import (
	"ac/internal/info"
)

type Informer struct {
	Infos []info.Info
}

func NewInformer() *Informer {
	return &Informer{}
}

func (i *Informer) Get() []info.Info {
	return i.Infos
}

func (i *Informer) Add(info info.Info) {
	i.Infos = append(i.Infos, info)
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
