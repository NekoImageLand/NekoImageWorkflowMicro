package impl

import (
	"NekoImageWorkflowMicro/client/scraper"
	"NekoImageWorkflowMicro/client/storage"
	"NekoImageWorkflowMicro/client/transfer"
	"NekoImageWorkflowMicro/common"
	"NekoImageWorkflowMicro/proto/clientTransform"
	"context"
	"github.com/sirupsen/logrus"
	microclient "go-micro.dev/v4/client"
)

type ClientImpl interface {
	// OnInit load client self config and before data
	OnInit() error
	// OnStart start client hook
	OnStart() error
	// PreUpload report pre upload data
	PreUpload(ctx context.Context, cli microclient.Client) error
	// PostUpload report post upload data
	PostUpload(ctx context.Context, cli microclient.Client) error
	// OnStop stop client hook
	OnStop() error
}

type ClientInstance struct {
	ClientImpl
	ClientInfo        *common.ClientConfig
	Scrapers          *[]scraper.ScraperInstance
	PreUploadBridge   *transfer.PreUploadTransBridgeInstance
	UploadTransBridge *transfer.UploadTransBridgeInstance
}

// OnInit load client self config and before data, then init Scrapers
func (ci *ClientInstance) OnInit() error {
	// init
	logrus.Debug("ClientInstance OnInit start")
	storage.LoadConfig(ci.ClientInfo)
	return nil
}

// OnStart currently do nothing
func (ci *ClientInstance) OnStart() error {
	logrus.Debug("ClientInstance OnStart start")
	return nil
}

// PreUpload report pre upload data
func (ci *ClientInstance) PreUpload(ctx context.Context, cli microclient.Client) error {
	logrus.Debug("ClientInstance PreUpload start")
	// TODO:
	preReq := &clientTransform.FilePreRequest{
		ClientID:   ci.ClientInfo.ClientID,
		ClientType: clientTransform.ClientType_LOCAL,
		FileUUID:   []string{"sample-uuid-from-local"},
	}
	preRsp := &clientTransform.FilePreResponse{}
	return cli.Call(ctx, cli.NewRequest("new.server", "FileUploadServiceServer.HandleFilePreUpload", preReq), preRsp)
}

// PostUpload report post upload data
func (ci *ClientInstance) PostUpload(ctx context.Context, cli microclient.Client) error {
	logrus.Debug("ClientInstance PostUpload start")
	// TODO:
	fileData := []*clientTransform.FileData{{
		FileUUID:    "sample-uuid-from-local",
		FileContent: []byte("local file content here"),
	}}
	postReq := &clientTransform.FilePostRequest{
		ClientID:   ci.ClientInfo.ClientID,
		ClientType: clientTransform.ClientType_LOCAL,
		LocalData:  fileData,
	}
	postRsp := &clientTransform.FilePostResponse{}
	return cli.Call(ctx, cli.NewRequest("new.server", "FileUploadServiceServer.HandleFilePostUpload", postReq), postRsp)
}

// OnStop write PreUploadBridge data and UploadTransBridge data to disk
func (ci *ClientInstance) OnStop() error {
	logrus.Debug("ClientInstance OnStop start")
	return nil
}
