package common

import (
	"fmt"
	"os"
)

type Env struct {
	HlsOutput     string
	GstDebugLevel string

	SqliteDbPath string
}

const (
	hlsOutputKey     = "HLS_OUTPUT"
	gstDebugLevelKey = "GST_DEBUG_LEVEL"
	sqliteDbPathKey  = "SQLITE_DB_PATH"
)

func ParseEnv() (Env, error) {

	hlsOutput := os.Getenv(hlsOutputKey)
	gstDebugLevel := os.Getenv(gstDebugLevelKey)
	sqliteDbPath := os.Getenv(sqliteDbPathKey)

	env := Env{
		HlsOutput:     hlsOutput,
		GstDebugLevel: gstDebugLevel,
		SqliteDbPath:  sqliteDbPath,
	}

	if hlsOutput == "" {
		return Env{}, fmt.Errorf("필수 환경변수가 설정되지 않았습니다: %s 또는 %s", hlsOutputKey)
	}

	return env, nil
}
