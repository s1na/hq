package api

import "time"

type Discovery struct {
	Name            string   `json:"name"`
	Address         string   `json:"address"`
	GithubWorkflows []string `json:"github_workflows"`
}

type ListingEntry struct {
	Name     string            `json:"name"`
	NTests   int               `json:"ntests"`
	Passes   int               `json:"passes"`
	Fails    int               `json:"fails"`
	Timeout  bool              `json:"timeout"`
	Clients  []string          `json:"clients"`
	Versions map[string]string `json:"versions"`
	Start    time.Time         `json:"start"`
	FileName string            `json:"fileName"`
	Size     int64             `json:"size"`
	SimLog   string            `json:"simLog"`
}

type TestSuiteResult struct {
	ID             int                  `json:"id"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	ClientVersions map[string]string    `json:"clientVersions"`
	TestCases      map[string]TestCase  `json:"testCases"`
	TestDetailsLog string               `json:"testDetailsLog"`
	SimLog         string               `json:"simLog"`
}

type TestCase struct {
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Start         time.Time             `json:"start"`
	End           time.Time             `json:"end"`
	SummaryResult SummaryResult         `json:"summaryResult"`
	ClientInfo    map[string]ClientInfo `json:"clientInfo"`
}

type SummaryResult struct {
	Pass    bool   `json:"pass"`
	Details string `json:"details"`
	Log     struct {
		Begin int64 `json:"begin"`
		End   int64 `json:"end"`
	} `json:"log"`
}

type ClientInfo struct {
	ID            string `json:"id"`
	IP            string `json:"ip"`
	Name          string `json:"name"`
	InstantiatedAt time.Time `json:"instantiatedAt"`
	LogFile       string `json:"logFile"`
}
