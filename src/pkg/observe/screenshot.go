package observe

import (
	"encoding/base64"
	"fmt"

	"github.com/llm-net/adbclaw/pkg/adb"
)

// ScreenshotResult holds the raw PNG data and optional base64 encoding.
type ScreenshotResult struct {
	Format string `json:"format"`
	Base64 string `json:"base64,omitempty"`
	Path   string `json:"path,omitempty"`
	Size   int    `json:"size_bytes"`
}

// TakeScreenshot captures the device screen via "adb exec-out screencap -p".
// Returns raw PNG bytes.
func TakeScreenshot(cmd adb.Commander) ([]byte, error) {
	data, err := cmd.ExecOut("screencap", "-p")
	if err != nil {
		return nil, fmt.Errorf("screencap failed: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("screencap returned empty data")
	}
	// Validate PNG header
	if len(data) < 8 || string(data[1:4]) != "PNG" {
		return nil, fmt.Errorf("screencap returned invalid PNG data (%d bytes)", len(data))
	}
	return data, nil
}

// ScreenshotAsBase64 captures a screenshot and returns a ScreenshotResult with base64 encoding.
func ScreenshotAsBase64(cmd adb.Commander) (*ScreenshotResult, error) {
	data, err := TakeScreenshot(cmd)
	if err != nil {
		return nil, err
	}
	return &ScreenshotResult{
		Format: "png",
		Base64: base64.StdEncoding.EncodeToString(data),
		Size:   len(data),
	}, nil
}
