package schemas




import (
	"sync"
)


type Manager struct {
	Mu      sync.Mutex
	Workers map[string]any
}


type WorkerMetrics struct {
    Count          int     `json:"count"`
    Speed          float64 `json:"speed"`
    Uncrawable     int     `json:"uncrawable"`
    TotalCrawlTime int64   `json:"totalCrawlTime"`
}

//all mem-bounded , working fine