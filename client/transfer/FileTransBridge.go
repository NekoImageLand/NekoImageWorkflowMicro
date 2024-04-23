package transfer

import (
	"NekoImageWorkflowMicro/client/model"
	"sync"
)

type FileTransBridge[T any] interface {
	Length() int
	Insert(number int, val T) error
	Pop(number int) error
}

type BaseFileTransBridgeInstance[T any] struct {
	Channel chan T
	FileTransBridge[T]
}

func (c *BaseFileTransBridgeInstance[T]) Length() int {
	return len(c.Channel)
}

func (c *BaseFileTransBridgeInstance[T]) Insert(number int, val T) error {
	for i := 0; i < number; i++ {
		c.Channel <- val
	}
	return nil
}

func (c *BaseFileTransBridgeInstance[T]) Pop(number int) error {
	for i := 0; i < number; i++ {
		<-c.Channel
	}
	return nil
}

type PreUploadTransBridgeInstance struct {
	BaseFileTransBridgeInstance[model.PreUploadFileData]
}

type UploadTransBridgeInstance struct {
	BaseFileTransBridgeInstance[model.UploadFileData]
}

var preUploadInstance *PreUploadTransBridgeInstance
var preUploadOnce sync.Once
var uploadInstance *UploadTransBridgeInstance
var uploadOnce sync.Once

func GetPreUploadTransBridgeInstance() *PreUploadTransBridgeInstance {
	preUploadOnce.Do(func() {
		preUploadInstance = &PreUploadTransBridgeInstance{
			BaseFileTransBridgeInstance: BaseFileTransBridgeInstance[model.PreUploadFileData]{
				Channel: make(chan model.PreUploadFileData),
			},
		}
	})
	return preUploadInstance
}

func GetUploadTransBridgeInstance() *UploadTransBridgeInstance {
	uploadOnce.Do(func() {
		uploadInstance = &UploadTransBridgeInstance{
			BaseFileTransBridgeInstance: BaseFileTransBridgeInstance[model.UploadFileData]{
				Channel: make(chan model.UploadFileData),
			},
		}
	})
	return uploadInstance
}
