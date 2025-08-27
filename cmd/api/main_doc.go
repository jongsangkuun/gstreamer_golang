package main

import (
	address "gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
)

// --- Swagger용 타입 로드 ---
var _ address.RTSPInformation

// 이 파일은 Swagger 문서 생성을 위한 더미(handler) 정의만 포함합니다.
// 실제 라우팅/로직은 main.go에 있고, 여기 함수들은 호출되지 않습니다.

// @title           GStreamer Media Server API
// @version         1.0
// @description     RTSP → HLS 변환 파이프라인 관리 및 상태 조회 API
// @BasePath        /
// @schemes         http

// Health
// @Summary Health check
// @Tags    system
// @Produce json
// @Success 200 {object} responseSchema
// @Router  /health [get]
func docHealth() {}

// 전체 파이프라인 조회
// @Summary Get all pipelines
// @Description 모든 파이프라인 정보를 반환합니다.
// @Tags    pipeline
// @Produce json
// @Success 200 {object} responseSchema
// @Router  /pipeline/all [get]
func docGetAllPipelines() {}

// 특정 파이프라인 조회
// @Summary Get a pipeline by name
// @Description streamName에 해당하는 파이프라인 정보를 반환합니다.
// @Tags    pipeline
// @Param   streamName path string true "Pipeline (stream) name"
// @Produce json
// @Success 200 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Router  /pipeline/{streamName} [get]
func docGetPipeline() {}

// 파이프라인 시작
// @Summary Start a pipeline
// @Description streamName 파이프라인을 실행 상태로 변경합니다.
// @Tags    pipeline
// @Param   streamName path string true "Pipeline (stream) name"
// @Produce json
// @Success 200 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Router  /pipeline/{streamName}/start [post]
func docStartPipeline() {}

// 파이프라인 정지
// @Summary Stop a pipeline
// @Description streamName 파이프라인을 정지 상태로 변경합니다.
// @Tags    pipeline
// @Param   streamName path string true "Pipeline (stream) name"
// @Produce json
// @Success 200 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Router  /pipeline/{streamName}/stop [post]
func docStopPipeline() {}

// 파이프라인 생성
// @Summary Create a pipeline
// @Description RTSP 정보를 받아 새 파이프라인을 생성합니다.
// @Tags    pipeline
// @Accept  json
// @Produce json
// @Param   rtsp body address.RTSPInformation true "RTSP info"
// @Success 200 {object} responseSchema
// @Failure 400 {object} responseSchema
// @Failure 409 {object} responseSchema
// @Failure 500 {object} responseSchema
// @Router  /pipeline [put]
func docCreatePipeline() {}

// 파이프라인 변경(덮어쓰기)
// @Summary Update (recreate) a pipeline
// @Description 같은 이름의 파이프라인이 있으면 새 설정으로 덮어씁니다.
// @Tags    pipeline
// @Accept  json
// @Produce json
// @Param   rtsp body address.RTSPInformation true "RTSP info"
// @Success 200 {object} responseSchema
// @Failure 400 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Failure 500 {object} responseSchema
// @Router  /pipeline [post]
func docUpdatePipeline() {}

// 파이프라인 삭제
// @Summary Delete a pipeline
// @Tags    pipeline
// @Param   streamName path string true "Pipeline (stream) name"
// @Produce json
// @Success 200 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Router  /pipeline/{streamName} [delete]
func docDeletePipeline() {}

// 모든 파이프라인 삭제
// @Summary Delete all pipelines
// @Tags    pipeline
// @Produce json
// @Success 200 {object} responseSchema
// @Failure 404 {object} responseSchema
// @Router  /pipeline/all [delete]
func docDeleteAllPipelines() {}
