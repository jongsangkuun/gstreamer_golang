package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/db"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/pipeline"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/service"
)

type responseSchema struct {
	Status  int         `json:"status"`  // HTTP Status Code
	Message string      `json:"message"` // OK, ERROR
	Data    interface{} `json:"data"`    // Additional Data
}

func main() {
	env, err := common.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Init()

	dbConn, err := db.ConnectSQLite(env.SqliteDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection()

	_, err = db.InitializeDatabase(env.SqliteDbPath)
	if err != nil {
		log.Fatal(err)
	}

	mainLoop, pm, err := service.GstServiceStart(env, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer service.GstServiceStop(mainLoop, pm)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    "",
		}
		c.JSON(200, response)
	})

	// 모든 파이브라인 정보 조회
	router.GET("/pipeline/all", func(c *gin.Context) {
		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    pm.GetAllPipelineInfo(),
		}

		c.JSON(200, response)
	})

	// streamName의 파이프라인 정보 조회
	router.GET("/pipeline/:streamName", func(c *gin.Context) {
		streamName := c.Param("streamName")
		pipelineInfo, exists := pm.GetPipeline(streamName)
		if !exists {
			response := responseSchema{
				Status:  404,
				Message: "Pipeline not found",
				Data:    streamName,
			}
			c.JSON(404, response)
			return
		}

		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    pipelineInfo,
		}
		c.JSON(200, response)
	})

	// streamName 파이프라인 정지
	router.POST("/pipeline/:streamName/stop", func(c *gin.Context) {
		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    "",
		}
		c.JSON(200, response)
	})

	// streamName 파이프라인 생성
	router.PUT("/pipeline", func(c *gin.Context) {
		var rtspInfo address.RTSPInformation
		if err := c.ShouldBindJSON(&rtspInfo); err != nil {
			response := responseSchema{
				Status:  400,
				Message: "fail",
				Data:    err.Error(),
			}
			c.JSON(400, response)
			return
		}

		// 이미 같은 이름의 파이프라인이 존재하는지 확인
		if _, exists := pm.GetPipeline(rtspInfo.Name); exists {
			response := responseSchema{
				Status:  409,
				Message: "fail",
				Data:    "Pipeline already exists",
			}
			c.JSON(409, response)
			return
		}

		success := pipeline.CreateStreamPipeline(pm, rtspInfo, env.HlsOutput)

		if !success {
			response := responseSchema{
				Status:  500,
				Message: "fail",
				Data:    "Failed to create pipeline",
			}
			c.JSON(500, response)
			return
		}

		_, err = db.CreateRTSPStream(dbConn, rtspInfo)
		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    "",
		}
		c.JSON(200, response)
	})

	// streamName 파이프라인 변경
	router.POST("/pipeline", func(c *gin.Context) {
		var rtspInfo address.RTSPInformation
		if err := c.ShouldBindJSON(&rtspInfo); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			return
		}

		// 같은 이름의 파이프라인이 존재하는지 확인
		// 파이프라인 존재 시 create로 덮어씌움
		// 이름이 없다면 에러 리턴
		if _, exists := pm.GetPipeline(rtspInfo.Name); exists {
			pipeline.CreateStreamPipeline(pm, rtspInfo, env.HlsOutput)
			response := responseSchema{
				Status:  200,
				Message: "ok",
				Data:    "",
			}
			c.JSON(200, response)
			return
		} else {
			response := responseSchema{
				Status:  404,
				Message: "fail",
				Data:    "Pipeline not found",
			}
			c.JSON(404, response)
		}

	})

	// streamName 파이프라인 삭제
	router.DELETE("/pipeline/:streamName", func(c *gin.Context) {
		streamName := c.Param("streamName")
		exists := pm.RemovePipeline(streamName)
		if exists != nil {
			response := responseSchema{
				Status:  404,
				Message: "fail",
				Data:    exists.Error(),
			}
			c.JSON(404, response)
			return
		}

		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    "",
		}
		c.JSON(200, response)
	})

	// streamName 파이프라인 삭제
	router.DELETE("/pipeline/all", func(c *gin.Context) {
		pm.CleanupAllPipelines()
		response := responseSchema{
			Status:  200,
			Message: "ok",
			Data:    "",
		}
		c.JSON(200, response)
	})

	err = router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(err)
	}
}
