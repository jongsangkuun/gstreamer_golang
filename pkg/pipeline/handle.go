package pipeline

import (
	"fmt"
	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/validation"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const restartDelay = 5 * time.Second
const outputDirMode = 0755

// PipelineManager는 파이프라인 관리를 위한 구조체입니다.
type PipelineManager struct {
	Pipelines []*PipelineInfo
	mu        sync.Mutex
}

// PipelineInfo는 개별 파이프라인 정보를 담는 구조체입니다.
type PipelineInfo struct {
	pipeline *gst.Pipeline
	name     string
}

// 새 파이프라인 매니저를 생성합니다.
func NewPipelineManager() *PipelineManager {
	return &PipelineManager{
		Pipelines: make([]*PipelineInfo, 0),
	}
}

// 파이프라인을 추가합니다.
func (pm *PipelineManager) AddPipeline(pipeline *gst.Pipeline, name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.Pipelines = append(pm.Pipelines, &PipelineInfo{pipeline: pipeline, name: name})
}

// 모든 파이프라인을 정리합니다.
func (pm *PipelineManager) CleanupAllPipelines() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for _, p := range pm.Pipelines {
		log.Info("정리 중:", p.name)
		p.pipeline.BlockSetState(gst.StateNull)
	}
}

// 각 RTSP 스트림에 대한 파이프라인을 생성합니다.
func CreatePipelines(pm *PipelineManager, wg *sync.WaitGroup, hwStatus validation.HardwareType, baseOutputDir string) {
	for _, stream := range address.RtspStreams {
		wg.Add(1)
		go func(stream address.RTSPStream) {
			defer wg.Done()
			createStreamPipeline(pm, stream, hwStatus, baseOutputDir)
		}(stream)
	}
}

// 단일 스트림에 대한 파이프라인을 생성합니다.
func createStreamPipeline(pm *PipelineManager, stream address.RTSPStream, hwStatus validation.HardwareType, baseOutputDir string) {
	// 스트림별 출력 디렉토리 생성
	outputDir := filepath.Join(baseOutputDir, stream.Name)
	if err := os.MkdirAll(outputDir, outputDirMode); err != nil {
		log.Error(fmt.Sprintf("[%s] 디렉토리 생성 오류: %v", stream.Name, err))
		return
	}

	// 파이프라인 문자열 생성
	var pipelineStr string
	if hwStatus == validation.HardwareTypeCPU {
		pipelineStr = CpuPipeline(stream.URL, outputDir, stream.Bitrate)
	} else {
		pipelineStr = GpuPipeline(stream.URL, outputDir, stream.Bitrate)
	}
	log.Info(fmt.Sprintf("[%s] 파이프라인 생성: %s", stream.Name, pipelineStr))

	// 파이프라인 생성
	pipelineFromString, err := gst.NewPipelineFromString(pipelineStr)
	if err != nil {
		log.Error(fmt.Sprintf("[%s] 파이프라인 생성 오류: %v", stream.Name, err))
		return
	}

	// 메시지 핸들러 추가
	setupPipelineMessageHandler(pipelineFromString, stream.Name)

	// 파이프라인 시작
	pipelineFromString.SetState(gst.StatePlaying)
	log.Info(fmt.Sprintf("[%s] 파이프라인 시작됨. HLS 경로: %s/index.m3u8", stream.Name, outputDir))

	// 전역 관리를 위해 파이프라인 저장
	pm.AddPipeline(pipelineFromString, stream.Name)
}

// 파이프라인 메시지 핸들러를 설정합니다.
func setupPipelineMessageHandler(pipeline *gst.Pipeline, streamName string) {
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			log.Info(fmt.Sprintf("[%s] 스트림 종료", streamName))
			if err := pipeline.BlockSetState(gst.StateNull); err != nil {
				return false
			}
		case gst.MessageError:
			err := msg.ParseError()
			log.Error(fmt.Sprintf("[%s] 오류: %s", streamName, err.Error()))
			if debug := err.DebugString(); debug != "" {
				log.Error(fmt.Sprintf("[%s] 디버그: %s", streamName, debug))
			}
			handlePipelineError(pipeline, streamName)
		case gst.MessageStateChanged:
			if msg.Source() == pipeline.GetName() {
				oldState, newState := msg.ParseStateChanged()
				log.Info(fmt.Sprintf("[%s] 상태 변경: %s -> %s",
					streamName, oldState.String(), newState.String()))
			}
		}
		return true
	})
}

// 파이프라인 오류를 처리합니다.
func handlePipelineError(pipeline *gst.Pipeline, streamName string) {
	pipeline.BlockSetState(gst.StateNull)
	go func() {
		log.Info(fmt.Sprintf("[%s] %s 후 재시작 시도...", streamName, restartDelay))
		time.Sleep(restartDelay)
		log.Info(fmt.Sprintf("[%s] 재시작 중...", streamName))
		pipeline.SetState(gst.StatePlaying)
	}()
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
