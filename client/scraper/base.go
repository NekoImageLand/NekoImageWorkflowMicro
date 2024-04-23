package scraper

import "NekoImageWorkflowMicro/common"
import "NekoImageWorkflowMicro/client/transfer"

type Scraper interface {
	// PrepareData prepare raw data, designed to be run in a goroutine once
	PrepareData() error
	// ProcessData make raw data to data which client can directly post, designed to be run in a goroutine once
	ProcessData() error
}

type ScraperInstance struct {
	Scraper
	common.ScraperType
	PreUploadBridge transfer.PreUploadTransBridgeInstance
}
