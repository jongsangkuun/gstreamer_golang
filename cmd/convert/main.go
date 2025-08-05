package main

import (
	"fmt"
	"github.com/go-gst/go-glib/glib"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/pipeline"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/validation"
	"os"
	"sync"
)

func main() {
	env, err := common.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Init()
	err = pipeline.InitGstreamer(env.GstDebugLevel)
	if err != nil {
		log.Fatal(err)
	}

	err = pipeline.SetupXDGRuntimeDir()
	if err != nil {
		log.Fatal(err)
	}

	// 하드웨어 상태 확인
	hwStatus := validation.DetectHardware()

	// 출력 디렉토리 설정
	outputDir := env.HlsOutput
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	// 메인 루프 생성
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	// 파이프라인 매니저 생성
	pipelineManager := pipeline.NewPipelineManager()

	// 스트림 처리를 위한 WaitGroup
	var wg sync.WaitGroup
	pipeline.CreatePipelines(pipelineManager, &wg, hwStatus, outputDir)

	// 모든 파이프라인이 생성될 때까지 대기
	wg.Wait()
	log.Info(fmt.Sprintf("모든 스트림 처리 중... 총 %d개 스트림", len(pipelineManager.Pipelines)))

	// 종료 시그널 처리
	pipeline.SetupSignalHandling(mainLoop, pipelineManager)

	// 메인 루프 실행
	log.Info("실행 중... Ctrl+C로 종료할 수 있습니다")
	mainLoop.Run()
	log.Info("프로그램 종료")
}
