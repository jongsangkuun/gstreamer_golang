package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-gst/go-gst/gst"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
)

const outputDirMode = 0755

// PipelineConfig는 파이프라인 구성에 필요한 설정을 담는 구조체
type PipelineConfig struct {
	RtspInformation *address.RTSPInformation
	OutputDir       string
}

// NewPipelineConfig는 기본 설정값으로 PipelineConfig를 생성합니다
func NewPipelineConfig(information address.RTSPInformation, hlsOutputDir string) *PipelineConfig {
	return &PipelineConfig{
		RtspInformation: &information,
		OutputDir:       hlsOutputDir,
	}
}

// BuildPipeline은 설정에 따라 파이프라인 문자열을 생성합니다
func BuildPipeline(rtspInfo *address.RTSPInformation, outputDir string) string {
	// 소스 부분 구성
	source := fmt.Sprintf("rtspsrc location=\"%s\" latency=0 buffer-mode=auto do-rtcp=true do-rtsp-keep-alive=true", rtspInfo.URL)

	// 처리 파이프라인 구성 (CPU/GPU에 따라 다름)
	var process string
	if rtspInfo.UseGPU {
		process = "decodebin ! videoconvert ! " +
			"videoconvert ! clockoverlay valignment=bottom halignment=right font-desc=\"Sans, 18\" shaded-background=true ! " +
			fmt.Sprintf("videoconvert ! nvh264enc bitrate=%d preset=low-latency-hq rc-mode=cbr "+
				"gop-size=30", rtspInfo.Bitrate)
	} else {
		process = "rtph264depay ! h264parse config-interval=1 ! " +
			"decodebin ! videoconvert ! " +
			"clockoverlay valignment=bottom halignment=right font-desc=\"Sans, 18\" shaded-background=true ! " +
			fmt.Sprintf("x264enc tune=zerolatency bitrate=%d speed-preset=superfast key-int-max=30", rtspInfo.Bitrate)
	}

	// 출력 부분 구성
	output := "h264parse config-interval=1 ! " +
		"mpegtsmux ! hlssink playlist-root=\".\" " +
		fmt.Sprintf("max-files=%d playlist-length=%d target-duration=%d "+
			"location=%s/segment_%%05d.ts playlist-location=%s/index.m3u8",
			rtspInfo.MaxFiles, rtspInfo.PlaylistLength, rtspInfo.TargetDuration,
			outputDir, outputDir)

	// 전체 파이프라인 문자열 반환
	return fmt.Sprintf("%s ! %s ! %s", source, process, output)
}

func CreatePipelines(pm *PipelineManager, wg *sync.WaitGroup, baseOutputDir string) {
	// RTSP 설정 정보 가져오기 (10개 스트림)
	rtspConfigs := address.DefaultRTSPInformations()

	if len(rtspConfigs) == 0 {
		log.Info("생성할 RTSP 스트림이 없습니다")
		return
	}

	log.Info(fmt.Sprintf("총 %d개의 RTSP 스트림 파이프라인 생성을 시작합니다", len(rtspConfigs)))

	for i, info := range rtspConfigs {
		wg.Add(1)
		go func(index int, info address.RTSPInformation) {
			defer wg.Done()

			streamName := info.Name
			log.Info(fmt.Sprintf("[%s] 파이프라인 #%d 생성 시작", streamName, index+1))

			success := CreateStreamPipeline(pm, info, baseOutputDir)
			if success {
				log.Info(fmt.Sprintf("[%s] 파이프라인 #%d 생성 완료", streamName, index+1))
			} else {
				log.Error(fmt.Sprintf("[%s] 파이프라인 #%d 생성 실패", streamName, index+1))
			}
		}(i, info)
	}
}

// 단일 스트림에 대한 파이프라인을 생성합니다.
func CreateStreamPipeline(pm *PipelineManager, rtspInfo address.RTSPInformation,
	baseOutputDir string) bool {

	streamName := rtspInfo.Name

	// 스트림별 출력 디렉토리 생성
	outputDir := filepath.Join(baseOutputDir, streamName)
	if err := os.MkdirAll(outputDir, outputDirMode); err != nil {
		log.Error(fmt.Sprintf("[%s] 디렉토리 생성 오류: %v", streamName, err))
		return false
	}

	pipelineStr := BuildPipeline(&rtspInfo, outputDir)
	log.Info(fmt.Sprintf("[%s] 파이프라인 문자열: %s", streamName, pipelineStr))

	// 파이프라인 생성
	gstPipeline, err := gst.NewPipelineFromString(pipelineStr)
	if err != nil {
		log.Error(fmt.Sprintf("[%s] 파이프라인 생성 오류: %v", streamName, err))
		return false
	}

	// 파이프라인 이름 설정
	pipelineName := fmt.Sprintf("pipeline-%s", streamName)
	gstPipeline.Object.SetProperty("name", pipelineName)

	// 전역 관리를 위해 파이프라인 저장 (시작 전에 저장)
	pm.AddPipeline(gstPipeline, &rtspInfo)

	// 메시지 핸들러 추가
	setupPipelineMessageHandler(gstPipeline, streamName, pm)

	// 파이프라인 시작
	if err := gstPipeline.SetState(gst.StatePlaying); err != nil {
		log.Error(fmt.Sprintf("[%s] 파이프라인 시작 오류: %v", streamName, err))
		// 실패 시 파이프라인 제거
		pm.RemovePipeline(streamName)
		return false
	}

	pm.UpdatePipelineStatus(streamName, StatusRunning)

	log.Info(fmt.Sprintf("[%s] 파이프라인 시작됨. HLS 경로: %s/index.m3u8",
		streamName, outputDir))

	return true
}

// 파이프라인 메시지 핸들러를 설정합니다.
func setupPipelineMessageHandler(pipeline *gst.Pipeline, streamName string, pm *PipelineManager) {
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			log.Info(fmt.Sprintf("[%s] 스트림 종료", streamName))
			pm.UpdatePipelineStatus(streamName, StatusStopped)
			if err := pipeline.BlockSetState(gst.StateNull); err != nil {
				log.Error(fmt.Sprintf("[%s] 파이프라인 정지 실패: %v", streamName, err))
				return false
			}

		case gst.MessageError:
			err := msg.ParseError()
			log.Error(fmt.Sprintf("[%s] 오류: %s", streamName, err.Error()))
			pm.UpdatePipelineStatus(streamName, StatusError)

			if debug := err.DebugString(); debug != "" {
				log.Error(fmt.Sprintf("[%s] 디버그: %s", streamName, debug))
			}
			handlePipelineError(pipeline, streamName, pm)

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
