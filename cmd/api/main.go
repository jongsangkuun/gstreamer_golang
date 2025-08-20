package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/pipeline"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/service"
)

func main() {
	env, err := common.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Init()
	mainLoop, pm, err := service.GstServiceStart(env)
	if err != nil {
		log.Fatal(err)
	}
	defer service.GstServiceStop(mainLoop, pm)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// 모든 파이브라인 정보 조회
	router.GET("/pipeline/all", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": pm.GetAllPipelineInfo(),
			// data 필드는 스웨거 완성 전 임시로 넣어둔 데이터
			"data": "0 == StatusCreating, 1 == StatusRunning, 2 == StatusStopped, 3 == StatusError",
		})

	})

	// streamName의 파이프라인 정보 조회
	router.GET("/pipeline/:streamName", func(c *gin.Context) {
		streamName := c.Param("streamName")
		pipelineInfo, exists := pm.GetPipeline(streamName)
		if !exists {
			c.JSON(404, gin.H{
				"error":      "Pipeline not found",
				"streamName": streamName,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": pipelineInfo,
		})
	})

	// streamName 파이프라인 정지
	router.POST("/pipeline/:streamName/stop", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// streamName 파이프라인 생성
	router.PUT("/pipeline", func(c *gin.Context) {
		var rtspInfo address.RTSPInformation
		if err := c.ShouldBindJSON(&rtspInfo); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			return
		}

		// 이미 같은 이름의 파이프라인이 존재하는지 확인
		if _, exists := pm.GetPipeline(rtspInfo.Name); exists {
			c.JSON(409, gin.H{
				"error":      "Pipeline already exists",
				"streamName": rtspInfo.Name,
			})
			return
		}

		success := pipeline.CreateStreamPipeline(pm, rtspInfo, env.HlsOutput)

		if !success {
			c.JSON(500, gin.H{
				"error":      "Failed to create pipeline",
				"streamName": rtspInfo.Name,
			})
			return
		}

		c.JSON(200, gin.H{
			"message":    "Pipeline created successfully",
			"streamName": rtspInfo.Name,
			"outputPath": fmt.Sprintf("%s/%s/index.m3u8", env.HlsOutput, rtspInfo.Name),
		})

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
			c.JSON(200, gin.H{
				"message": "ok",
			})
			return
		} else {
			c.JSON(404, gin.H{
				"error":      "Pipeline not found",
				"streamName": rtspInfo.Name,
			})
		}

	})

	// streamName 파이프라인 삭제
	router.DELETE("/pipeline/:streamName", func(c *gin.Context) {
		streamName := c.Param("streamName")
		exists := pm.RemovePipeline(streamName)
		if exists != nil {
			c.JSON(404, gin.H{
				"error":      "Pipeline not found",
				"streamName": streamName,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	// streamName 파이프라인 삭제
	router.DELETE("/pipeline/all", func(c *gin.Context) {
		pm.CleanupAllPipelines()
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
