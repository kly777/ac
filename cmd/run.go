package main

import(
	"ac/internal/role/manager"
	"ac/internal/taskManager"
	"ac/internal/role/developer"
	"ac/internal/informer"
)


func Run(){
	taskManager := taskManager.NewTaskManager()
	informer:= informer.NewInformer()

	manager := manager.NewManager(taskManager, informer)
}