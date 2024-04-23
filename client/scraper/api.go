package scraper

type APIScraperInstance struct {
	ScraperInstance
}

func (c *APIScraperInstance) PrepareData() error {
	// Preparation steps for API client
	return nil
}

func (c *APIScraperInstance) ProcessData() error {
	// Simulate fetching and processing data from API
	return nil
}
