package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"time"
)

func main() {
	// Set up RabbitMQ broker
	b := rabbitmq.NewBroker(
		rabbitmq.ExchangeName("my_new_exchange"),
	)
	if err := b.Init(); err != nil {
		logger.Fatalf("Broker init error: %v", err)
	}
	if err := b.Connect(); err != nil {
		logger.Fatalf("Broker connect error: %v", err)
	}
	// Set up Consul registry
	reg := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service := micro.NewService(
		micro.Name("new.server.another"),
		micro.Address("127.0.0.1:30102"),
		micro.Broker(b),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
	)
	service.Init()

	_, err := b.Subscribe("greeter.new.topic", func(p broker.Event) error {
		logger.Infof("Received message: %s", string(p.Message().Body))
		fmt.Println("Received message:", string(p.Message().Body))
		return p.Ack()
	}, broker.Queue("greeter_queue"), broker.DisableAutoAck())

	if err != nil {
		logger.Fatalf("Failed to subscribe: %v", err)
	}

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("anotherserver is now running and subscribed to greeter.topic")
}
