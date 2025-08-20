package pipeline

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
)

const RestartDelay = 5 * time.Second

// 파이프라인 에러 처리 함수 (누락된 함수 구현)
func handlePipelineError(pipeline *gst.Pipeline, streamName string, pm *PipelineManager) {
	log.Error(fmt.Sprintf("[%s] 파이프라인 에러 처리 시작", streamName))

	if err := pipeline.BlockSetState(gst.StateNull); err != nil {
		log.Error(fmt.Sprintf("[%s] 파이프라인 정지 실패: %v", streamName, err))
	}

	pm.UpdatePipelineStatus(streamName, StatusError)
}

// 종료 시그널 처리를 설정합니다.
func SetupSignalHandling(mainLoop *glib.MainLoop, pm *PipelineManager) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Info("\n종료 요청 받음...")
		pm.CleanupAllPipelines()
		mainLoop.Quit()
	}()
}
