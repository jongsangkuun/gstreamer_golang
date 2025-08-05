package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// fileServerHandler는 지정된 디렉터리로부터 파일을 서빙하는 Gin 핸들러를 생성합니다.
func fileServerHandler(directory string) gin.HandlerFunc {
	log.Println("File Server Directory: ", directory)
	fileServer := http.FileServer(http.Dir(directory))

	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if urlPath == "/" {
			urlPath = "" // 루트 경로 처리
		}

		c.Request.URL.Path = urlPath
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {

	// 설정에서 읽어온 HLS 디렉터리 경로를 가져옵니다.
	directory := "./hls_output"

	// Gin 기본 모드 설정 (선택사항: debug, release, test)
	gin.SetMode(gin.ReleaseMode)

	// Gin 라우터 생성
	router := gin.New()

	// 미들웨어 설정
	router.Use(gin.Recovery()) // 패닉 복구 미들웨어

	// CORS 미들웨어 설정
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// 정적 파일 서버 핸들러 등록
	router.Any("/*path", fileServerHandler(directory))

	log.Println("HTTP 파일 서버가 6001번 포트에서 실행 중입니다. CORS는 비활성화 되어 있습니다.")

	// HTTP 서버를 6001번 포트에서 실행하며, 오류가 발생하면 로그와 함께 서버가 종료됩니다.
	err := router.Run(":6001")
	if err != nil {
		log.Fatal(err)
	}
}
