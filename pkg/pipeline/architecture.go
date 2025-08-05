package pipeline

import "fmt"

// 파이프라인 설정을 위한 상수
const (
	// HLS 설정
	defaultMaxFiles       = 10
	defaultPlaylistLength = 5
	defaultTargetDuration = 1
)

// PipelineConfig는 파이프라인 구성에 필요한 설정을 담는 구조체
type PipelineConfig struct {
	URL            string
	OutputDir      string
	MaxFiles       int
	PlaylistLength int
	TargetDuration int
	Bitrate        int
	UseGPU         bool
}

// NewPipelineConfig는 기본 설정값으로 PipelineConfig를 생성합니다
func NewPipelineConfig(url, outputDir string, useGPU bool, bitrate int) *PipelineConfig {
	return &PipelineConfig{
		URL:            url,
		OutputDir:      outputDir,
		MaxFiles:       defaultMaxFiles,
		PlaylistLength: defaultPlaylistLength,
		TargetDuration: defaultTargetDuration,
		UseGPU:         useGPU,
		Bitrate:        bitrate,
	}
}

// BuildPipeline은 설정에 따라 파이프라인 문자열을 생성합니다
func (pc *PipelineConfig) BuildPipeline() string {
	// 소스 부분 구성
	source := fmt.Sprintf("rtspsrc location=\"%s\" latency=0 buffer-mode=auto do-rtcp=true do-rtsp-keep-alive=true", pc.URL)

	// 처리 파이프라인 구성 (CPU/GPU에 따라 다름)
	var process string
	if pc.UseGPU {
		process = "decodebin ! videoconvert ! " +
			"videoconvert ! clockoverlay valignment=bottom halignment=right font-desc=\"Sans, 18\" shaded-background=true ! " +
			fmt.Sprintf("videoconvert ! nvh264enc bitrate=%d preset=low-latency-hq rc-mode=cbr "+
				"gop-size=30", pc.Bitrate)
	} else {
		process = "rtph264depay ! h264parse config-interval=1 ! " +
			"decodebin ! videoconvert ! " +
			"clockoverlay valignment=bottom halignment=right font-desc=\"Sans, 18\" shaded-background=true ! " +
			fmt.Sprintf("x264enc tune=zerolatency bitrate=%d speed-preset=superfast key-int-max=30", pc.Bitrate)
	}

	// 출력 부분 구성
	output := "h264parse config-interval=1 ! " +
		"mpegtsmux ! hlssink playlist-root=\".\" " +
		fmt.Sprintf("max-files=%d playlist-length=%d target-duration=%d "+
			"location=%s/segment_%%05d.ts playlist-location=%s/index.m3u8",
			pc.MaxFiles, pc.PlaylistLength, pc.TargetDuration,
			pc.OutputDir, pc.OutputDir)

	// 전체 파이프라인 문자열 반환
	return fmt.Sprintf("%s ! %s ! %s", source, process, output)
}

// CpuPipeline는 CPU용 파이프라인 문자열을 생성합니다
func CpuPipeline(url string, outputDir string, bitrate int) string {
	config := NewPipelineConfig(url, outputDir, false, bitrate)
	return config.BuildPipeline()
}

// GpuPipeline은 GPU용 파이프라인 문자열을 생성합니다
func GpuPipeline(url string, outputDir string, bitrate int) string {
	config := NewPipelineConfig(url, outputDir, true, bitrate)
	// GPU 파이프라인은 최대 파일 수를 20으로 설정
	config.MaxFiles = 20
	return config.BuildPipeline()
}
