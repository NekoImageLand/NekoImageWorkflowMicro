package scraper

type LocalScraperInstance struct {
	ScraperInstance
}

func (c *LocalScraperInstance) PrepareData() error {
	// Preparation steps for API client
	return nil
}

func (c *LocalScraperInstance) ProcessData() error {
	// Simulate fetching and processing data from API
	return nil
}
