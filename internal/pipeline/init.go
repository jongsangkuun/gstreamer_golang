package pipeline

import (
	"fmt"
	"os"

	"github.com/go-gst/go-gst/gst"
)

// GStreamer를 초기화합니다.
func InitGstreamer(gstDebugLevel string) error {
	// GStreamer 초기화
	gst.Init(nil)

	// GStreamer 디버그 레벨 설정
	if gstDebugLevel != "" {
		if err := os.Setenv("GST_DEBUG", gstDebugLevel); err != nil {
			return fmt.Errorf("GStreamer 디버그 레벨 설정 실패: %w", err)
		}
	}

	return nil
}
