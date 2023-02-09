package video

import (
	"github.com/kordar/video-collection/strategy"
	log "github.com/sirupsen/logrus"
)

var (
	ReadingStatus  uint32 = 1
	RunningStatus  uint32 = 2
	FinishedStatus uint32 = 3
)

type StreamManager struct {
	baseUrl string
	streams map[string]strategy.Strategy // 视频策略
	status  map[string]uint32            // 状态 1:reading, 2:running, 3:finished
	buffer  chan string
}

func NewStreamManager(baseUrl string, bufSize int) *StreamManager {
	return &StreamManager{
		streams: make(map[string]strategy.Strategy),
		status:  make(map[string]uint32),
		baseUrl: baseUrl,
		buffer:  make(chan string, bufSize),
	}
}

func (m *StreamManager) Add(s strategy.Strategy) {
	if s.GetBaseDir() == "" {
		s.SetBaseDir(m.baseUrl)
	}
	id := s.GetId()
	// 如果流不存在进入准备开启状态
	if m.streams[id] == nil {
		m.streams[id] = s
		m.status[id] = ReadingStatus
		m.buffer <- id
	}
}

func (m *StreamManager) Run() {
	go func() {
		strategy.WaitRestartHandler.Run()
		for {
			id := <-m.buffer
			state := m.status[id]
			if state == ReadingStatus {
				go m.start(id)
			}
		}
	}()
}

func (m *StreamManager) start(id string) {
	if m.streams[id] == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorln(r)
		}
	}()

	m.status[id] = RunningStatus
	stream := m.streams[id]
	err := stream.Execute()
	if err != nil {
		log.Errorln(err)
	}
	m.status[id] = FinishedStatus
	// 启动失败或结束后，进行重启操作
	status := stream.GetStatus()
	if status == "restart" {
		log.Infof("************ 尝试重启服务, Id = %s **************", id)
		stream.Refresh()
		m.status[id] = ReadingStatus
		m.buffer <- id
	}
	if status == "exit" {
		log.Infof("************ 退出服务, Id = %s **************", id)
		delete(m.status, id)
		delete(m.streams, id)
	}
}

// Stop 停止Progress
func (m *StreamManager) Stop(id string) {
	if m.streams[id] == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorln(r)
		}
	}()

	m.streams[id].Stop()
}
