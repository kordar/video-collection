package command

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	StartStatusReady = iota
	StartStatusRunning
	RetryStatusReady
	RetryStatusFinish
	RetryStatusExit
)

type StreamManager struct {
	commands map[string]ICommand // 视频命令
	buffer   chan string
	locker   sync.Mutex
	status   map[string]int // 0 未启动 1 已启动
}

func NewStreamManager(bufSize int) *StreamManager {
	return &StreamManager{
		commands: make(map[string]ICommand),
		buffer:   make(chan string, bufSize),
		status:   make(map[string]int),
	}
}

func (s *StreamManager) Add(c ICommand) bool {
	s.locker.Lock()
	defer s.locker.Unlock()
	id := c.GetId()
	// 如果流不存在进入准备开启状态
	if s.commands[id] == nil {
		s.commands[id] = c
		s.status[id] = StartStatusReady
		s.buffer <- id
		return true
	}
	return false
}

func (s *StreamManager) Run() {
	go func() {
		for {
			id := <-s.buffer
			status := s.status[id]
			if status == StartStatusReady {
				go s.start(id)
			}
		}
	}()
}

func (s *StreamManager) start(id string) {
	if s.commands[id] == nil {
		delete(s.status, id)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorln(r)
			// 抛出异常尝试接触id关系绑定
			delete(s.status, id)
			delete(s.commands, id)
		}
	}()

	s.status[id] = StartStatusRunning
	stream := s.commands[id]
	err := stream.Execute()
	if err != nil {
		log.Errorln(err)
	}
	// 启动失败或结束后，进行重启操作
	status := stream.GetStatus()
	if status == "restart" {
		log.Infof("************ 尝试重启服务, Id = %s **************", id)
		stream.Refresh()
		s.status[id] = StartStatusReady
		s.buffer <- id
	}
	if status == "exit" {
		log.Infof("************ 退出服务, Id = %s **************", id)
		delete(s.status, id)
		delete(s.commands, id)
	}
}

// Stop 停止Progress
func (s *StreamManager) Stop(id string) {
	if s.commands[id] == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorln(r)
		}
	}()

	s.commands[id].Stop()
}

func (s *StreamManager) StartCheckDeath(spec string) {
	c := cron.New()
	_, _ = c.AddFunc(spec, func() {
		for id, cmd := range s.commands {
			now := time.Now().Unix()
			if cmd.GetProgressRestartSeconds() > 0 && now-cmd.GetProgressRefreshTime().Unix() > cmd.GetProgressRestartSeconds() {
				log.Info("强制停止编码器，id = ", id)
				cmd.Stop()
			}
		}
	})
	c.Start()
}

func (s *StreamManager) GetStreamData() []ICommand {
	var data = make([]ICommand, 0)
	for _, command := range s.commands {
		data = append(data, command)
	}
	return data
}
