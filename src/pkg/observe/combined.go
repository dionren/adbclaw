package observe

import (
	"sync"

	"github.com/llm-net/adbclaw/pkg/adb"
)

// ObserveResult holds the combined output of screenshot + UI tree.
type ObserveResult struct {
	Screenshot *ScreenshotResult `json:"screenshot,omitempty"`
	UI         *UITree           `json:"ui,omitempty"`
	Errors     []string          `json:"errors,omitempty"`
}

// Observe captures both screenshot and UI tree in parallel.
// Partial failure is tolerated — one failing doesn't block the other.
func Observe(cmd adb.Commander) *ObserveResult {
	result := &ObserveResult{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	// Screenshot goroutine
	go func() {
		defer wg.Done()
		ss, err := ScreenshotAsBase64(cmd)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, "screenshot: "+err.Error())
		} else {
			result.Screenshot = ss
		}
	}()

	// UI tree goroutine
	go func() {
		defer wg.Done()
		tree, err := DumpUITree(cmd)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, "uitree: "+err.Error())
		} else {
			result.UI = tree
		}
	}()

	wg.Wait()

	// nil out empty errors slice for cleaner JSON
	if len(result.Errors) == 0 {
		result.Errors = nil
	}

	return result
}
