package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gitlab.hds-robotcenter.com/gstreamer-convert/internal/address"
	pipe "gitlab.hds-robotcenter.com/gstreamer-convert/internal/pipeline"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// RTSP Stream 정보를 위한 GORM 모델
type RTSPStream struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	URL            string    `gorm:"not null;size:500" json:"url"`
	Bitrate        int       `gorm:"not null;default:1000" json:"bitrate"`
	UseGPU         bool      `gorm:"default:false" json:"use_gpu"`
	MaxFiles       int       `gorm:"default:10" json:"max_files"`
	PlaylistLength int       `gorm:"default:5" json:"playlist_length"`
	TargetDuration int       `gorm:"default:1" json:"target_duration"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// GORM SQLite 연결
func ConnectSQLite(dbPath string) (*gorm.DB, error) {
	// GORM 설정
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 로그 레벨 설정
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// SQLite 드라이버 옵션 설정
	dsn := fmt.Sprintf("%s?cache=shared&mode=rwc&_journal_mode=DELETE&_synchronous=FULL&_cache_size=10000&_foreign_keys=on", dbPath)

	// GORM으로 SQLite 연결
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("SQLite GORM 연결 실패: %v", err)
	}

	// 기본 SQL DB 인스턴스 가져오기 (연결 풀 설정용)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("SQL DB 인스턴스 가져오기 실패: %v", err)
	}

	// 연결 풀 설정
	setupConnectionPool(sqlDB)

	// 연결 테스트
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("SQLite 연결 테스트 실패: %v", err)
	}

	DB = db
	log.Printf("SQLite GORM 데이터베이스 연결 성공: %s", dbPath)
	return db, nil
}

// 연결 풀 설정
func setupConnectionPool(sqlDB *sql.DB) {
	// SQLite는 단일 파일이므로 연결 수를 적게 설정
	sqlDB.SetMaxOpenConns(10)                  // 최대 열린 연결 수
	sqlDB.SetMaxIdleConns(5)                   // 최대 유휴 연결 수
	sqlDB.SetConnMaxLifetime(time.Hour)        // 연결 최대 생명주기
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 유휴 연결 최대 시간
}

// RTSPInformation을 RTSPStream으로 변환하는 메서드
func (r *RTSPStream) FromRTSPInformation(rtsp address.RTSPInformation) {
	r.ID = uint(rtsp.Id)
	r.Name = rtsp.Name
	r.URL = rtsp.RtspUrl
	r.Bitrate = rtsp.Bitrate
	r.UseGPU = rtsp.UseGPU
	r.MaxFiles = rtsp.MaxFiles
	r.PlaylistLength = rtsp.PlaylistLength
	r.TargetDuration = rtsp.TargetDuration
}

// RTSPStream을 RTSPInformation으로 변환하는 메서드
func (r *RTSPStream) ToRTSPInformation() address.RTSPInformation {
	return address.RTSPInformation{
		Id:             int(r.ID),
		Name:           r.Name,
		RtspUrl:        r.URL,
		Bitrate:        r.Bitrate,
		UseGPU:         r.UseGPU,
		MaxFiles:       r.MaxFiles,
		PlaylistLength: r.PlaylistLength,
		TargetDuration: r.TargetDuration,
	}
}

// 자동 마이그레이션
func AutoMigrate(db *gorm.DB) error {
	log.Println("데이터베이스 마이그레이션 시작...")

	// RTSP Stream 모델을 마이그레이션
	err := db.AutoMigrate(&RTSPStream{})
	if err != nil {
		return fmt.Errorf("마이그레이션 실패: %v", err)
	}

	log.Println("데이터베이스 마이그레이션 완료")
	return nil
}

// 초기 RTSP Stream 데이터 삽입
func SeedData(db *gorm.DB) error {
	log.Println("RTSP Stream 초기 데이터 삽입 시작...")

	// address 패키지에서 기본 RTSP 정보 가져오기
	defaultRTSPs := address.DefaultRTSPInformations()

	for _, rtspInfo := range defaultRTSPs {
		var existingStream RTSPStream

		// ID로 먼저 확인
		if err := db.Where("id = ?", rtspInfo.Id).First(&existingStream).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// ID가 없으면 이름으로도 확인
				if err := db.Where("name = ?", rtspInfo.Name).First(&existingStream).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// 새 스트림 생성
						newStream := &RTSPStream{}
						newStream.FromRTSPInformation(rtspInfo)

						if err := db.Create(newStream).Error; err != nil {
							return fmt.Errorf("RTSP Stream 생성 실패 (%s): %v", rtspInfo.Name, err)
						}
						log.Printf("RTSP Stream 생성됨: %s (ID: %d)", rtspInfo.Name, rtspInfo.Id)
					} else {
						return fmt.Errorf("RTSP Stream 조회 실패: %v", err)
					}
				} else {
					log.Printf("RTSP Stream이 이미 존재함 (이름): %s", rtspInfo.Name)
				}
			} else {
				return fmt.Errorf("RTSP Stream 조회 실패: %v", err)
			}
		} else {
			log.Printf("RTSP Stream이 이미 존재함 (ID): %d - %s", rtspInfo.Id, rtspInfo.Name)
		}
	}

	log.Println("RTSP Stream 초기 데이터 삽입 완료")
	return nil
}

// RTSP Stream 관련 데이터베이스 함수들
func GetAllActiveRTSPStreams(db *gorm.DB) ([]*RTSPStream, error) {
	var streams []*RTSPStream
	err := db.Find(&streams).Error
	return streams, err
}

func GetRTSPStreamByName(db *gorm.DB, name string) (*RTSPStream, error) {
	var stream RTSPStream
	err := db.Where("name = ? AND is_active = ?", name, true).First(&stream).Error
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func GetRTSPStreamByID(db *gorm.DB, id uint) (*RTSPStream, error) {
	var stream RTSPStream
	err := db.Where("id = ? AND is_active = ?", id, true).First(&stream).Error
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func CreateRTSPStream(db *gorm.DB, rtspInfo address.RTSPInformation) (*RTSPStream, error) {
	stream := &RTSPStream{}
	stream.FromRTSPInformation(rtspInfo)

	err := db.Create(stream).Error
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func UpdateRTSPStream(db *gorm.DB, rtspInfo address.RTSPInformation) (*RTSPStream, error) {
	var stream RTSPStream
	if err := db.Where("name = ?", rtspInfo.Name).First(&stream).Error; err != nil {
		return nil, err
	}

	// 기존 ID 유지하면서 정보 업데이트
	originalID := stream.ID
	stream.FromRTSPInformation(rtspInfo)
	stream.ID = originalID

	err := db.Save(&stream).Error
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func UpdateStreamActiveStatus(db *gorm.DB, streamName string, status pipe.PipelineStatus) error {
	if status == pipe.StatusRunning {
		return db.Model(&RTSPStream{}).Where("name = ?", streamName).Update("is_active", true).Error
	} else {
		return db.Model(&RTSPStream{}).Where("name = ?", streamName).Update("is_active", false).Error
	}
}

func DeleteSoftRTSPStream(db *gorm.DB, name string) error {
	// 소프트 삭제 (is_active를 false로 설정)
	return db.Model(&RTSPStream{}).Where("name = ?", name).Update("is_active", false).Error
}

func DeleteHardRTSPStream(db *gorm.DB, name string) error {
	// 하드 삭제
	return db.Model(&RTSPStream{}).Where("name = ?", name).Delete(&RTSPStream{}).Error
}

// GORM v1.20+ 에서는 WHERE 조건 없이 모든 레코드를 삭제하는 것을 방지함. 따라서 1=1 같은 조건을 입력
func DeleteAllRTSPStream(db *gorm.DB) error {
	return db.Model(&RTSPStream{}).Where("1 = 1").Delete(&RTSPStream{}).Error
}

// 모든 RTSP Stream을 RTSPInformation 형태로 반환
func GetAllRTSPInformations(db *gorm.DB) ([]*RTSPStream, error) {
	streams, err := GetAllActiveRTSPStreams(db)
	if err != nil {
		return nil, err
	}

	var rtspStream []*RTSPStream
	for _, stream := range streams {
		rtspStream = append(rtspStream, stream)
	}

	return rtspStream, nil
}

// 데이터베이스 전체 초기화
func InitializeDatabase(dbPath string) (*gorm.DB, error) {
	// GORM SQLite 연결
	db, err := ConnectSQLite(dbPath)
	if err != nil {
		return nil, err
	}

	// 자동 마이그레이션
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	// 레코드 수 확인
	var count int64
	if err := db.Model(&RTSPStream{}).Count(&count).Error; err != nil {
		return nil, err
	}

	// DB에 값이 없을 때만 초기 데이터 삽입
	if count == 0 {
		if err := SeedData(db); err != nil {
			log.Printf("초기 데이터 삽입 중 오류: %v", err)
			// 오류가 있어도 계속 진행
		}
	}

	return db, nil
}

// 연결 닫기
func CloseConnection() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("SQL DB 인스턴스 가져오기 실패: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("데이터베이스 연결 닫기 실패: %v", err)
	}

	return nil
}

// Connection Pool 상태 조회
func GetConnectionStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{"error": "데이터베이스 연결이 없습니다"}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{"error": "SQL DB 인스턴스 가져오기 실패"}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
