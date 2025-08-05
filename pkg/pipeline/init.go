package pipeline

import (
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"os"
)

const (
	runtimeDirMode   = 0700
	runtimeDirPrefix = "/tmp/runtime-"
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

// SetupXDGRuntimeDir는 XDG_RUNTIME_DIR 환경 변수가 설정되어 있지 않은 경우
// 임시 디렉토리를 생성하고 해당 환경 변수를 설정합니다.
// 이는 일부 리눅스 애플리케이션에서 필요한 XDG 규격을 준수하기 위함입니다.
func SetupXDGRuntimeDir() error {
	// 이미 환경 변수가 설정되어 있으면 아무 작업도 필요 없음
	if os.Getenv("XDG_RUNTIME_DIR") != "" {
		return nil
	}

	// 임시 디렉토리 생성
	tmpDir := runtimeDirPrefix + os.Getenv("USER")
	if err := os.MkdirAll(tmpDir, runtimeDirMode); err != nil {
		return fmt.Errorf("XDG_RUNTIME_DIR 디렉토리 생성 실패: %w", err)
	}

	// 환경 변수 설정
	if err := os.Setenv("XDG_RUNTIME_DIR", tmpDir); err != nil {
		return fmt.Errorf("XDG_RUNTIME_DIR 환경 변수 설정 실패: %w", err)
	}

	return nil
}
