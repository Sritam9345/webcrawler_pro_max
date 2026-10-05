package schemas




import (
	"sync"
)


type Manager struct {
	Mu      sync.Mutex
	Workers map[string]any
}