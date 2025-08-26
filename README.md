# GStreamer RTSP to HLS 변환기

이 프로젝트는 RTSP 스트림을 HLS(HTTP Live Streaming) 형식으로 변환하는 서비스입니다. CPU와 GPU 환경 모두에서 실행 가능하며, 다중 스트림을 동시에 처리할 수 있습니다.

## 환경 변수 설정 (.env 파일)

프로젝트는 `.env` 파일을 통해 환경 설정을 관리합니다. 주요 환경 변수는 다음과 같습니다:

```
ENV_STATE=dev                  # 환경 상태 (dev 또는 prod)

# 개발 환경 설정
DEV_DEBUG=true                 # 디버그 모드 활성화 여부
DEV_GST_DEBUG_LEVEL=5          # GStreamer 디버그 레벨 (0-5)
DEV_HLS_OUTPUT=/app/hls_output # HLS 출력 디렉토리 경로
DEV_HLS_BACKUP=/app/hls_backup # HLS 백업 디렉토리 경로

# 프로덕션 환경 설정
PROD_DEBUG=false               # 프로덕션 환경에서는 디버그 모드 비활성화
PROD_GST_DEBUG_LEVEL=0         # 프로덕션 환경에서는 GStreamer 디버그 비활성화
PROD_HLS_OUTPUT=/app/hls_output # HLS 출력 디렉토리 경로
PROD_HLS_BACKUP=/app/hls_backup # HLS 백업 디렉토리 경로
```

## Docker Compose 실행 방법

### CPU 환경에서 실행

CPU 환경에서는 다음 명령어로 서비스를 실행합니다:

```bash
docker-compose -f docker-compose.cpu.yml up -d
```

### GPU 환경에서 실행

NVIDIA GPU가 있는 환경에서는 다음 명령어로 GPU 가속을 활용하여 서비스를 실행합니다:

```bash
docker-compose -f docker-compose.gpu.yml up -d
```

### 서비스 중지

서비스를 중지하려면 다음 명령어를 사용합니다:

```bash
# CPU 환경
docker-compose -f docker-compose.cpu.yml down

# GPU 환경
docker-compose -f docker-compose.gpu.yml down
```

## 프로젝트 전체 프로세스

### 1. 초기화 과정

1. 환경 변수 로드 및 검증
2. GStreamer 초기화
3. 하드웨어 감지 (CPU/GPU)
4. 출력 디렉토리 설정

### 2. 파이프라인 생성 및 실행

1. 각 RTSP 스트림에 대해 별도의 파이프라인 생성
2. 하드웨어 환경(CPU/GPU)에 따라 적절한 인코딩 파이프라인 구성
   - CPU: x264enc 인코더 사용
   - GPU: nvh264enc 인코더 사용 (NVIDIA GPU 가속)
3. 각 스트림별로 지정된 출력 디렉토리에 HLS 세그먼트 및 재생 목록 생성

### 3. 스트림 모니터링 및 오류 처리

1. 각 파이프라인의 상태 변화 모니터링
2. 오류 발생 시 자동 재시작 메커니즘 작동 (5초 후 재시작)
3. 종료 신호(SIGINT, SIGTERM) 수신 시 모든 파이프라인 정상 종료

### 4. HLS 스트림 제공

1. 별도의 파일 서버가 HLS 세그먼트 및 재생 목록 제공
2. 포트 6001을 통해 웹 클라이언트에서 접근 가능

## 주요 구성 요소

- **rtsp2hls 서비스**: RTSP 스트림을 HLS로 변환하는 메인 서비스
- **file_server 서비스**: 변환된 HLS 콘텐츠를 HTTP를 통해 제공하는 서비스
- **스트림 설정**: `internal/address/rtsp_address.go`에서 RTSP 스트림 소스 정의
- **파이프라인 관리**: `pkg/pipeline` 패키지에서 GStreamer 파이프라인 생성 및 관리

## 시스템 요구 사항

### CPU 환경
- Docker 및 Docker Compose 설치
- GStreamer 라이브러리 지원

### GPU 환경 (추가 요구 사항)
- NVIDIA GPU
- NVIDIA 드라이버 설치
- NVIDIA Container Toolkit 설치
- Docker 및 Docker Compose 설치

## SQLite
- sqlite3 을 통해서 postgres 대채
- 굳이 postgres 컨테이너를 올려서 작업할 이유가 없음
- 어짜피 미디어서버는 하나의 온프레미스 환경에서 동작할 것으로 보임
- 또한 많은 양의 데이터를 저장할 것은 아님
- 로그에 관련된 데이터는 추후 다른 로그 서버 등의 서비스를 통해서 진행하는게 좋아보임
- connection pool 도 동작하고 백업도 하드백업을 통해서 처리하기 쉬워보임
- 