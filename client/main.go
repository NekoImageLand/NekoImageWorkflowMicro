package main

import (
	"NekoImageWorkflowMicro/client/impl"
	"NekoImageWorkflowMicro/client/scraper"
	"NekoImageWorkflowMicro/client/transfer"
	"NekoImageWorkflowMicro/common"
	FinalLog "NekoImageWorkflowMicro/log"
	FinalLogAdapter "NekoImageWorkflowMicro/log/adapter"
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

func createCustomZapLogger() (*zap.Logger, error) {
	_logger := FinalLogAdapter.NewCustomLogger()
	return _logger, nil
}

func hookEtcdLogger(reg registry.Registry) error {
	newLogger, _ := createCustomZapLogger()
	val := reflect.ValueOf(reg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	clientField := val.FieldByName("client")
	if !clientField.IsValid() {
		return fmt.Errorf("client field not found")
	}
	clientVal := clientField.Elem()
	lgField := clientVal.FieldByName("lg")
	if !lgField.IsValid() {
		return fmt.Errorf("lg field not found")
	}
	lgFieldPtr := (*unsafe.Pointer)(unsafe.Pointer(lgField.UnsafeAddr()))
	*lgFieldPtr = unsafe.Pointer(newLogger)
	return nil
}

func main() {
	ctx := context.Background()
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&FinalLog.CustomFormatter{})
	// go micro logger modify
	logger.DefaultLogger = &FinalLogAdapter.LogrusAdapter{}
	client := impl.ClientInstance{
		ClientInfo:        &common.ClientConfig{},
		Scrapers:          new([]scraper.ScraperInstance),
		PreUploadBridge:   transfer.GetPreUploadTransBridgeInstance(),
		UploadTransBridge: transfer.GetUploadTransBridgeInstance(),
	}
	if err := client.OnInit(); err != nil {
		logrus.Error("OnInit error:", err)
	}
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("bais"),
		registry.Logger(&FinalLogAdapter.LogrusAdapter{}),
	)
	if err := hookEtcdLogger(etcdRegistry); err != nil {
		logrus.Fatal(err)
	}
	service := micro.NewService(
		micro.Name(client.ClientInfo.ClientName),
		micro.Address(client.ClientInfo.ClientRegisterAddress),
		micro.Registry(etcdRegistry),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Logger(&FinalLogAdapter.LogrusAdapter{}),
	)
	service.Init()
	if err := client.OnStart(); err != nil {
		logrus.Error("OnStart error:", err)
	}
	defer func() {
		if err := client.OnStop(); err != nil {
			logrus.Error("OnStop error:", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Warning("Received shutdown signal")
		os.Exit(0)
	}()
	// TODO: run scrapers in goroutine
	go func() {
		for _, scraperInstance := range *client.Scrapers {
			go scraperInstance.PrepareData()
			go scraperInstance.ProcessData()
		}
	}()
	for {
		if err := client.PreUpload(ctx, service.Client()); err != nil {
			logrus.Error("PreUpload error:", err)
		}
		if err := client.PostUpload(ctx, service.Client()); err != nil {
			logrus.Error("PostUpload error:", err)
		}
		time.Sleep(time.Second * time.Duration(client.ClientInfo.PostUploadPeriod))
	}
}
