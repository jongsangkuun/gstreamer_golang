package service

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-gst/go-glib/glib"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	pipe "gitlab.hds-robotcenter.com/gstreamer-convert/internal/pipeline"
)

func GstServiceStart(env common.Env) (*glib.MainLoop, *pipe.PipelineManager, error) {
	err := pipe.InitGstreamer(env.GstDebugLevel)
	if err != nil {
		return nil, nil, err
	}

	err = common.SetupXDGRuntimeDir()
	if err != nil {
		return nil, nil, err
	}

	// 출력 디렉토리 설정
	outputDir := env.HlsOutput
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	// 메인 루프 생성
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	// 파이프라인 매니저 생성
	pipelineManager := pipe.NewPipelineManager()
	// 스트림 처리를 위한 WaitGroup
	var wg sync.WaitGroup
	pipe.CreatePipelines(pipelineManager, &wg, outputDir)

	// 모든 파이프라인이 생성될 때까지 대기
	wg.Wait()
	log.Info(fmt.Sprintf("모든 스트림 처리 중... 총 %d개 스트림", len(pipelineManager.Pipelines)))
	return mainLoop, pipelineManager, nil
}

func GstServiceStop(mainLoop *glib.MainLoop, pipelineManager *pipe.PipelineManager) {
	// 종료 시그널 처리
	pipe.SetupSignalHandling(mainLoop, pipelineManager)
	// 메인 루프 실행
	log.Info("실행 중... Ctrl+C로 종료할 수 있습니다")

	mainLoop.Run()
	log.Info("프로그램 종료")

}
