package validation

import (
	"os/exec"
	"runtime"
)

// Architecture는 시스템 아키텍처를 나타내는 타입입니다.
type Architecture string
type HardwareType string

const (
	ArchitectureAMD64 Architecture = "amd64"
	ArchitectureARM64 Architecture = "arm64"
	ArchitectureX86   Architecture = "x86_64"
)

const (
	HardwareTypeCPU HardwareType = "CPU"
	HardwareTypeGPU HardwareType = "GPU"
)

// ProcessorInfo는 프로세서 정보를 담는 구조체입니다.
type ProcessorInfo struct {
	Architecture Architecture // 문자열 대신 Architecture 타입 사용
	HardwareType HardwareType // 문자열 대신 HardwareType 타입 사용
}

// HasNvidiaGPU는 시스템에 NVIDIA GPU가 있는지 확인합니다.
func HasNvidiaGPU() HardwareType {
	cmd := exec.Command("nvidia-smi")
	err := cmd.Run()
	if err != nil {
		return HardwareTypeCPU
	}
	return HardwareTypeGPU
}

// 하드웨어 타입을 감지합니다.
func DetectHardware() HardwareType {
	if HasNvidiaGPU() == "GPU" && runtime.GOARCH == "amd64" {
		return HardwareTypeGPU
	} else {
		return HardwareTypeCPU
	}
}
