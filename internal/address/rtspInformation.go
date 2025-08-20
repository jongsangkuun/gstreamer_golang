package address

const (
	//인코딩 설정 - 화질별 비트레이트 (kbps)
	BitrateLD  = 500  // 저화질 (480p)
	BitrateSD  = 1000 // 표준화질 (720p)
	BitrateHD  = 2000 // 고화질 (1080p)
	BitrateFHD = 4000 // 풀HD
)

type RTSPInformation struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Bitrate        int    `json:"bitrate"`
	UseGPU         bool   `json:"use_gpu"`
	MaxFiles       int    `json:"max_files"`
	PlaylistLength int    `json:"playlist_length"`
	TargetDuration int    `json:"target_duration"`
}

type RTSPInformationList []RTSPInformation

func DefaultRTSPInformations() RTSPInformationList {
	rtspInformationList := RTSPInformationList{
		{0, "hanwha1", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD, false, 10, 5, 1},
		{1, "hanwha2", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD, false, 10, 5, 1},
		{2, "hanwha3", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD, false, 10, 5, 1},
		{3, "bluecop1-1", "rtsp://192.168.1.123:8554/stream1", BitrateSD, false, 10, 5, 1},
		{4, "bluecop2-1", "rtsp://192.168.1.123:8554/stream2", BitrateSD, false, 10, 5, 1},
		{5, "bluecop3-1", "rtsp://192.168.1.123:8554/stream3", BitrateSD, false, 10, 5, 1},
		{6, "bluecop4-1", "rtsp://192.168.1.123:8554/stream4", BitrateSD, false, 10, 5, 1},
		{7, "bluecop5-1", "rtsp://192.168.1.123:8554/stream5", BitrateSD, false, 10, 5, 1},
		{8, "bluecop6-1", "rtsp://192.168.1.123:8554/stream6", BitrateSD, false, 10, 5, 1},
		{9, "bluecop7-1", "rtsp://192.168.1.123:8554/stream7", BitrateSD, false, 10, 5, 1},
	}

	return rtspInformationList
}
