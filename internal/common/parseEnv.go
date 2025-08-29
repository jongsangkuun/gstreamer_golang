package common

import (
	"fmt"
	"os"
)

type Env struct {
	HlsOutput     string
	GstDebugLevel string

	SqliteDbPath string

	FileServerHost string
	FileServerPort string
}

const (
	hlsOutputKey      = "HLS_OUTPUT"
	gstDebugLevelKey  = "GST_DEBUG_LEVEL"
	sqliteDbPathKey   = "SQLITE_DB_PATH"
	fileServerHostKey = "FILE_SERVER_HOST"
	fileServerPortKey = "FILE_SERVER_PORT"
)

func ParseEnv() (Env, error) {

	hlsOutput := os.Getenv(hlsOutputKey)
	gstDebugLevel := os.Getenv(gstDebugLevelKey)
	sqliteDbPath := os.Getenv(sqliteDbPathKey)
	fileServerHost := os.Getenv(fileServerHostKey)
	fileServerPort := os.Getenv(fileServerPortKey)

	env := Env{
		HlsOutput:      hlsOutput,
		GstDebugLevel:  gstDebugLevel,
		SqliteDbPath:   sqliteDbPath,
		FileServerHost: fileServerHost,
		FileServerPort: fileServerPort,
	}

	if hlsOutput == "" {
		return Env{}, fmt.Errorf("필수 환경변수가 설정되지 않았습니다: %s 또는 %s", hlsOutputKey)
	}

	return env, nil
}
