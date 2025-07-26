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
