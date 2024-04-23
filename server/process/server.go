package main

import (
	"NekoImageWorkflowMicro/proto/clientTransform"
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"time"
)

var service micro.Service

type FileUploadServiceServer struct {
}

func AddMsgToQueue(msg string) {
	if err := service.Options().Broker.Publish("greeter.new.topic", &broker.Message{
		Body: []byte(msg),
	}); err != nil {
		logger.Errorf("Failed to publish: %v", err)
	} else {
		logger.Info("Message published successfully")
	}
}

func (f *FileUploadServiceServer) HandleFilePreUpload(ctx context.Context, req *clientTransform.FilePreRequest, rsp *clientTransform.FilePreResponse) error {
	// Handle the file pre-upload request
	AddMsgToQueue(fmt.Sprintf("hello from %s in HandleFilePreUpload", req.ClientID))
	logger.Info("Received file pre-upload request")
	return nil
}

func (f *FileUploadServiceServer) HandleFilePostUpload(ctx context.Context, req *clientTransform.FilePostRequest, rsp *clientTransform.FilePostResponse) error {
	// Handle the file post-upload request
	AddMsgToQueue(fmt.Sprintf("hello from %s in HandleFilePostUpload", req.ClientID))
	logger.Info("Received file pre-upload resquest")
	return nil
}

func main() {
	// Initialize the RabbitMQ broker
	b := rabbitmq.NewBroker(
		rabbitmq.ExchangeName("my_new_exchange"),
	)
	if err := b.Init(); err != nil {
		logger.Fatal(err)
	}
	if err := b.Connect(); err != nil {
		logger.Fatal(err)
	}
	// Initialize the Consul registry
	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service = micro.NewService(
		micro.Name("new.server"),
		micro.Address("127.0.0.1:30101"),
		micro.Broker(b),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
	)
	service.Init()

	if err := service.Server().Handle(
		service.Server().NewHandler(new(FileUploadServiceServer)),
	); err != nil {
		logger.Fatal(err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

}
