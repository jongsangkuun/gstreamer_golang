package main

import (
	address "gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
)

// --- Swagger용 타입 로드 ---
var _ address.RTSPInformation
var _ responseSchema
var _ GetPipeLineResponse
var _ RTSPInformationRequest

// 이 파일은 Swagger 문서 생성을 위한 더미(handler) 정의만 포함합니다.
// 실제 라우팅/로직은 main.go에 있고, 여기 함수들은 호출되지 않습니다.

// @title           GStreamer Media Server API
// @version         1.0
// @description     RTSP → HLS 변환 파이프라인 관리 및 상태 조회 API
// @BasePath        /
// @schemes         http

// Health check
// @Summary Health check
// @Description 서버 상태를 확인합니다
// @Tags    system
// @Produce json
// @Success 200 {object} responseSchema{data=string} "서버 정상"
// @Router  /health [get]
func docHealth() {}

// 전체 파이프라인 조회
// @Summary Get all pipelines
// @Description 모든 파이프라인 설정 정보를 DB에서 조회하여 GetPipeLineResponse 배열로 반환합니다
// @Tags    pipeline
// @Produce json
// @Success 200 {object} responseSchema{data=[]GetPipeLineResponse} "모든 파이프라인 조회 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인 조회 실패"
// @Router  /pipeline/all [get]
func docGetAllPipelines() {}

// 활성 파이프라인 조회
// @Summary Get all active pipelines
// @Description 동작중인 모든 활성 파이프라인 정보를 조회하여 GetPipeLineResponse 배열로 반환합니다
// @Tags    pipeline
// @Produce json
// @Success 200 {object} responseSchema{data=[]GetPipeLineResponse} "활성 파이프라인 조회 성공"
// @Router  /pipeline/active/all [get]
func docGetAllActivePipelines() {}

// 특정 파이프라인 조회
// @Summary Get a pipeline by name
// @Description streamName에 해당하는 파이프라인 정보를 DB에서 조회하여 반환합니다
// @Tags    pipeline
// @Param   streamName path string true "파이프라인 이름"
// @Produce json
// @Success 200 {object} responseSchema{data=GetPipeLineResponse} "파이프라인 조회 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인을 찾을 수 없음"
// @Router  /pipeline/{streamName} [get]
func docGetPipeline() {}

// 특정 파이프라인 상태 조회 (구현 예정)
// @Summary Get pipeline status by name
// @Description streamName 파이프라인의 실시간 상태 정보를 반환합니다 (TODO: 구현 예정)
// @Tags    pipeline
// @Param   streamName path string true "파이프라인 이름"
// @Produce json
// @Success 200 {object} responseSchema "파이프라인 상태 조회"
// @Router  /pipeline/{streamName}/status [get]
func docGetPipelineStatus() {}

// 파이프라인 시작
// @Summary Start a pipeline
// @Description streamName 파이프라인을 실행 상태로 변경하고 DB에 상태를 업데이트합니다
// @Tags    pipeline
// @Param   streamName path string true "시작할 파이프라인 이름"
// @Produce json
// @Success 200 {object} responseSchema{data=string} "파이프라인 시작 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인 시작 실패"
// @Router  /pipeline/{streamName}/start [put]
func docStartPipeline() {}

// 파이프라인 정지
// @Summary Stop a pipeline
// @Description streamName 파이프라인을 정지 상태로 변경하고 DB에 상태를 업데이트합니다
// @Tags    pipeline
// @Param   streamName path string true "정지할 파이프라인 이름"
// @Produce json
// @Success 200 {object} responseSchema{data=string} "파이프라인 정지 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인 정지 실패"
// @Router  /pipeline/{streamName}/stop [put]
func docStopPipeline() {}

// 파이프라인 업데이트 (PUT)
// @Summary Update a pipeline
// @Description 기존 파이프라인을 새 설정으로 업데이트하고 DB에 반영합니다
// @Tags    pipeline
// @Accept  json
// @Produce json
// @Param   rtsp body RTSPInformationRequest true "업데이트할 RTSP 정보 (ID 필드 제외)"
// @Success 200 {object} responseSchema{data=RTSPInformationRequest} "파이프라인 업데이트 성공"
// @Failure 400 {object} responseSchema{data=string} "잘못된 요청 본문"
// @Failure 404 {object} responseSchema{data=string} "파이프라인을 찾을 수 없음"
// @Failure 500 {object} responseSchema{data=string} "파이프라인 업데이트 또는 DB 반영 실패"
// @Router  /pipeline [put]
func docUpdatePipelinePut() {}

// 파이프라인 생성 (POST)
// @Summary Create a new pipeline
// @Description RTSP 정보를 받아 새 파이프라인을 생성하고 DB에 저장합니다
// @Tags    pipeline
// @Accept  json
// @Produce json
// @Param   rtsp body RTSPInformationRequest true "새 파이프라인 RTSP 정보 (ID 필드 제외)"
// @Success 200 {object} responseSchema{data=string} "파이프라인 생성 성공"
// @Failure 400 {object} responseSchema{data=string} "잘못된 요청 본문"
// @Failure 409 {object} responseSchema{data=string} "파이프라인이 이미 존재함"
// @Failure 500 {object} responseSchema{data=string} "파이프라인 생성 또는 DB 저장 실패"
// @Router  /pipeline [post]
func docCreatePipelinePost() {}

// 파이프라인 삭제
// @Summary Delete a pipeline
// @Description streamName 파이프라인을 정지하고 삭제한 후 DB에서도 하드 삭제합니다
// @Tags    pipeline
// @Param   streamName path string true "삭제할 파이프라인 이름"
// @Produce json
// @Success 200 {object} responseSchema{data=string} "파이프라인 삭제 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인 삭제 실패"
// @Router  /pipeline/{streamName} [delete]
func docDeletePipeline() {}

// 모든 파이프라인 삭제
// @Summary Delete all pipelines
// @Description 모든 파이프라인을 정리하고 DB에서도 모두 삭제합니다
// @Tags    pipeline
// @Produce json
// @Success 200 {object} responseSchema{data=string} "모든 파이프라인 삭제 성공"
// @Failure 404 {object} responseSchema{data=string} "파이프라인 삭제 실패"
// @Router  /pipeline/all [delete]
func docDeleteAllPipelines() {}
