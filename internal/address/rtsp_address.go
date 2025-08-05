package address

type RTSPStream struct {
	Name    string
	URL     string
	Bitrate int
}

const (
	//인코딩 설정 - 화질별 비트레이트 (kbps)
	BitrateLD  = 500  // 저화질 (480p)
	BitrateSD  = 1000 // 표준화질 (720p)
	BitrateHD  = 2000 // 고화질 (1080p)
	BitrateFHD = 4000 // 풀HD
)

var RtspStreams = []RTSPStream{
	{"hanwha1", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD},
	{"hanwha2", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD},
	{"hanwha3", "rtsp://admin:hdshds3112@192.168.1.25:558/LiveChannel/0/media.smp", BitrateSD},
	{"bluecop1-1", "rtsp://192.168.1.123:8554/stream1", BitrateSD},
	{"bluecop2-1", "rtsp://192.168.1.123:8554/stream2", BitrateSD},
	{"bluecop3-1", "rtsp://192.168.1.123:8554/stream3", BitrateSD},
	{"bluecop4-1", "rtsp://192.168.1.123:8554/stream4", BitrateSD},
	{"bluecop5-1", "rtsp://192.168.1.123:8554/stream5", BitrateSD},
	{"bluecop6-1", "rtsp://192.168.1.123:8554/stream6", BitrateSD},
	{"bluecop7-1", "rtsp://192.168.1.123:8554/stream7", BitrateSD},
}
