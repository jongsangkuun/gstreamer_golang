package pipeline

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-gst/go-gst/gst"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
)

const OutputDirMode = 0755

// PipelineManager는 파이프라인 관리를 위한 구조체입니다.
type PipelineManager struct {
	Pipelines map[string]*PipelineInfo // key는 RtspInformation.Name
	mu        sync.RWMutex
}

// PipelineInfo는 개별 파이프라인 정보를 담는 구조체입니다.
type PipelineInfo struct {
	pipeline        *gst.Pipeline            `json:"pipeline"`
	RtspInformation *address.RTSPInformation `json:"RtspInformation"`
	Status          PipelineStatus           `json:"status"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type PipelineStatus int

const (
	StatusCreating PipelineStatus = iota
	StatusRunning
	StatusStopped
	StatusError
)

func NewPipelineManager() *PipelineManager {
	return &PipelineManager{
		Pipelines: make(map[string]*PipelineInfo),
	}
}

// 파이프라인을 추가합니다.
func (pm *PipelineManager) AddPipeline(pipeline *gst.Pipeline, rtspInformation *address.RTSPInformation) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	streamName := rtspInformation.Name
	pipelineInfo := &PipelineInfo{
		pipeline:        pipeline,
		RtspInformation: rtspInformation,
		Status:          StatusCreating,
		CreatedAt:       time.Now(),
	}

	pm.Pipelines[streamName] = pipelineInfo
	log.Info(fmt.Sprintf("[%s] 파이프라인 추가됨", streamName))
}

// 특정 파이프라인 조회 (스트림 이름으로)
func (pm *PipelineManager) GetPipeline(streamName string) (*PipelineInfo, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, exists := pm.Pipelines[streamName]
	return info, exists
}

func (pm *PipelineManager) UpdatePipeline(streamName string, pipeline *gst.Pipeline) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if _, exists := pm.Pipelines[streamName]; !exists {
		return fmt.Errorf("streamName %s not found", streamName)
	}

	pm.Pipelines[streamName].pipeline = pipeline

	return nil
}

// 파이프라인 상태 업데이트
func (pm *PipelineManager) UpdatePipelineStatus(streamName string, status PipelineStatus) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if info, exists := pm.Pipelines[streamName]; exists {
		info.Status = status
		log.Info(fmt.Sprintf("[%s] 파이프라인 상태 변경: %d", streamName, status))
	}
}

// 파이프라인 제거
func (pm *PipelineManager) RemovePipeline(streamName string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	info, exists := pm.Pipelines[streamName]
	if !exists {
		return fmt.Errorf("streamName %s not found", streamName)
	}

	info.pipeline.BlockSetState(gst.StateNull)
	delete(pm.Pipelines, streamName)
	log.Info(fmt.Sprintf("[%s] 파이프라인 제거됨", streamName))
	return nil
}

// 모든 파이프라인 정보 조회
func (pm *PipelineManager) GetAllPipelineInfo() map[string]*PipelineInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*PipelineInfo)
	for name, info := range pm.Pipelines {
		result[name] = info
	}
	return result
}

// 모든 파이프라인을 정리합니다.
func (pm *PipelineManager) CleanupAllPipelines() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for streamName, info := range pm.Pipelines {
		log.Info(fmt.Sprintf("[%s] 정리 중", streamName))
		info.pipeline.BlockSetState(gst.StateNull)
	}
	pm.Pipelines = make(map[string]*PipelineInfo) // 맵 초기화
}

// 파이프라인 상태 모니터링
func (pm *PipelineManager) MonitorPipelines() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pm.mu.RLock()
		running := 0
		stopped := 0
		errors := 0

		log.Info("=== 파이프라인 상태 요약 ===")
		for streamName, info := range pm.Pipelines {
			switch info.Status {
			case StatusRunning:
				running++
			case StatusStopped:
				stopped++
			case StatusError:
				errors++
			}

			log.Info(fmt.Sprintf("스트림: %s | 상태: %d | 생성시간: %s",
				streamName, info.Status, info.CreatedAt.Format("15:04:05")))
		}

		log.Info(fmt.Sprintf("실행중: %d, 정지: %d, 오류: %d", running, stopped, errors))
		pm.mu.RUnlock()
	}
}
