package common

import (
	"fmt"
	"os"
)

type Env struct {
	HlsOutput     string
	HlsBackup     string
	GstDebugLevel string
}

func ParseEnv() (Env, error) {
	const (
		hlsOutputKey     = "HLS_OUTPUT"
		hlsBackupKey     = "HLS_BACKUP"
		gstDebugLevelKey = "GST_DEBUG_LEVEL"
	)

	hlsOutput := os.Getenv(hlsOutputKey)
	hlsBackup := os.Getenv(hlsBackupKey)
	gstDebugLevel := os.Getenv(gstDebugLevelKey)

	env := Env{
		HlsOutput:     hlsOutput,
		HlsBackup:     hlsBackup,
		GstDebugLevel: gstDebugLevel,
	}

	if hlsOutput == "" && hlsBackup == "" {
		return Env{}, fmt.Errorf("필수 환경변수가 설정되지 않았습니다: %s 또는 %s", hlsOutputKey, hlsBackupKey)
	}

	return env, nil

}
