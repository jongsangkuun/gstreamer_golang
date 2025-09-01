package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.hds-robotcenter.com/gstreamer-convert/cmd/api/docs"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/common"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/db"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/log"
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/pipeline"
	"gitlab.hds-robotcenter.com/gstreamer-convert/pkg/service"
	"gorm.io/gorm"
)

type responseSchema struct {
	Status  int         `json:"status"`  // HTTP Status Code
	Message string      `json:"message"` // success, fail
	Data    interface{} `json:"data"`    // Additional Data
}

type GetPipeLineResponse struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	RtspURL        string `json:"rtsp_url"`
	HlsURL         string `json:"hls_url"`
	Bitrate        int    `json:"bitrate"`
	UseGPU         bool   `json:"use_gpu"`
	MaxFiles       int    `json:"max_files"`
	PlaylistLength int    `json:"playlist_length"`
	TargetDuration int    `json:"target_duration"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func ConvertDBtoResponse(dbData db.RTSPStream, env common.Env) GetPipeLineResponse {
	return GetPipeLineResponse{
		Id:             int(dbData.ID),
		Name:           dbData.Name,
		RtspURL:        dbData.URL,
		HlsURL:         fmt.Sprintf("%s:%s/%s/index.m3u8", env.FileServerHost, env.FileServerPort, dbData.Name),
		Bitrate:        dbData.Bitrate,
		UseGPU:         dbData.UseGPU,
		MaxFiles:       dbData.MaxFiles,
		PlaylistLength: dbData.PlaylistLength,
		TargetDuration: dbData.TargetDuration,
		IsActive:       dbData.IsActive,
		CreatedAt:      dbData.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      dbData.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertPipelineInfoToResponse(rtspInfo *address.RTSPInformation, env common.Env) GetPipeLineResponse {
	return GetPipeLineResponse{
		Id:             rtspInfo.Id,
		Name:           rtspInfo.Name,
		RtspURL:        rtspInfo.URL,
		HlsURL:         fmt.Sprintf("%s:%s/%s/index.m3u8", env.FileServerHost, env.FileServerPort, rtspInfo.Name),
		Bitrate:        rtspInfo.Bitrate,
		UseGPU:         rtspInfo.UseGPU,
		MaxFiles:       rtspInfo.MaxFiles,
		PlaylistLength: rtspInfo.PlaylistLength,
		TargetDuration: rtspInfo.TargetDuration,
		IsActive:       true,                                     // 기본값 또는 별도 로직 필요
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"), // 임시 값
		UpdatedAt:      time.Now().Format("2006-01-02 15:04:05"), // 임시 값
	}
}

func createResponse(status int, message string, data interface{}) responseSchema {
	return responseSchema{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func healthHandler(c *gin.Context) {
	response := createResponse(http.StatusOK, "success", "")
	c.JSON(http.StatusOK, response)
}

func getAllPipelinesHandler(dbConn *gorm.DB, env common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("getAllPipelinesHandler")
		data, err := db.GetAllRTSPInformations(dbConn)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		var info []GetPipeLineResponse
		for _, stream := range data {
			info = append(info, ConvertDBtoResponse(*stream, env))
		}

		response := createResponse(http.StatusOK, "success", info)
		c.JSON(http.StatusOK, response)
	}
}

func getAllActivePipelinesHandler(pm *pipeline.PipelineManager, env common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("getAllActivePipelinesHandler")
		data := pm.GetAllPipelineInfo()

		var info []GetPipeLineResponse
		for streamName, stream := range data {
			if stream == nil {
				log.WithField("streamName", streamName).Warn("파이프라인 정보가 nil입니다. 건너뜁니다.")
				continue
			}

			if stream.RtspInformation == nil {
				log.WithField("streamName", streamName).Warn("RtspInformation이 nil입니다. 건너뜁니다.")
				continue
			}

			convertedData := ConvertPipelineInfoToResponse(stream.RtspInformation, env)
			info = append(info, convertedData)
		}

		response := createResponse(http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	}
}

func getPipelineByNameHandler(dbConn *gorm.DB, env common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamName := c.Param("streamName")
		stream, err := db.GetRTSPStreamByName(dbConn, streamName)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		info := ConvertDBtoResponse(*stream, env)
		response := createResponse(http.StatusOK, "success", info)
		c.JSON(http.StatusOK, response)
	}
}

func startPipelineHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pipelineName := c.Param("streamName")
		err := pm.UpdatePipelineStatus(pipelineName, pipeline.StatusRunning)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		err = db.UpdateStreamActiveStatus(dbConn, pipelineName, pipeline.StatusRunning)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := createResponse(http.StatusOK, "success", "")
		c.JSON(http.StatusOK, response)
	}
}

func stopPipelineHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pipelineName := c.Param("streamName")
		err := pm.UpdatePipelineStatus(pipelineName, pipeline.StatusStopped)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		err = db.UpdateStreamActiveStatus(dbConn, pipelineName, pipeline.StatusStopped)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := createResponse(http.StatusOK, "success", "")
		c.JSON(http.StatusOK, response)
	}
}

func createPipelineHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager, env common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rtspInfo address.RTSPInformation
		if err := c.ShouldBindJSON(&rtspInfo); err != nil {
			response := createResponse(http.StatusBadRequest, "fail", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if _, exists := pm.GetPipeline(rtspInfo.Name); exists {
			response := createResponse(http.StatusConflict, "fail", "Pipeline already exists")
			c.JSON(http.StatusConflict, response)
			return
		}

		success := pipeline.CreateStreamPipeline(pm, rtspInfo, env.HlsOutput)
		if !success {
			response := createResponse(http.StatusInternalServerError, "fail", "Failed to create pipeline")
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// DB에 RTSP 스트림 정보 저장
		_, err := db.CreateRTSPStream(dbConn, rtspInfo)
		if err != nil {
			response := createResponse(http.StatusInternalServerError, "fail", err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := createResponse(http.StatusOK, "success", "")
		c.JSON(http.StatusOK, response)
	}
}

func updatePipelineHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager, env common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rtspInfo address.RTSPInformation
		if err := c.ShouldBindJSON(&rtspInfo); err != nil {
			response := createResponse(http.StatusBadRequest, "fail", err.Error())
			c.JSON(http.StatusBadRequest, response)
			return
		}

		var rtspData address.RTSPInformation
		if _, exists := pm.GetPipeline(rtspInfo.Name); exists {
			pipeline.CreateStreamPipeline(pm, rtspInfo, env.HlsOutput)
			_, err := db.UpdateRTSPStream(dbConn, rtspInfo)
			if err != nil {
				response := createResponse(http.StatusInternalServerError, "fail", err.Error())
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		} else {
			response := createResponse(http.StatusNotFound, "fail", "Pipeline not found")
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := createResponse(http.StatusOK, "success", rtspData)
		c.JSON(http.StatusOK, response)
	}
}

func deletePipelineHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamName := c.Param("streamName")
		err := pm.UpdatePipelineStatus(streamName, pipeline.StatusStopped)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		err = pm.RemovePipeline(streamName)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		err = db.DeleteHardRTSPStream(dbConn, streamName)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := createResponse(http.StatusOK, "success", "")
		c.JSON(http.StatusOK, response)
	}
}

func deleteAllPipelinesHandler(dbConn *gorm.DB, pm *pipeline.PipelineManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pm.CleanupAllPipelines()

		err := db.DeleteAllRTSPStream(dbConn)
		if err != nil {
			response := createResponse(http.StatusNotFound, "fail", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := createResponse(http.StatusOK, "success", "")
		c.JSON(http.StatusOK, response)
	}
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", healthHandler)

	router.GET("/pipeline/all", getAllPipelinesHandler(dbConn, env))
	router.GET("/pipeline/active/all", getAllActivePipelinesHandler(pm, env))
	router.GET("/pipeline/:streamName", getPipelineByNameHandler(dbConn, env))
	router.GET("/pipeline/:streamName/status", func(c *gin.Context) {})

	router.POST("/pipeline/:streamName/start", startPipelineHandler(dbConn, pm))
	router.POST("/pipeline/:streamName/stop", stopPipelineHandler(dbConn, pm))
	router.POST("/pipeline", updatePipelineHandler(dbConn, pm, env))

	router.PUT("/pipeline", createPipelineHandler(dbConn, pm, env))

	router.DELETE("/pipeline/:streamName", deletePipelineHandler(dbConn, pm))
	router.DELETE("/pipeline/all", deleteAllPipelinesHandler(dbConn, pm))

	router.Run()
}
