package futureinfo

import (
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/future"
	"time"
)

type Storage struct {
	LastUpdate time.Time
	infoMap    map[uint64]*decoder.SecDefMessage
}

func NewStorage() *Storage {
	return &Storage{
		LastUpdate: time.Now(),
		infoMap:    map[uint64]*decoder.SecDefMessage{},
	}
}

func (s *Storage) Add(msg *decoder.SecDefMessage) {
	s.infoMap[msg.SecurityID] = msg
	s.LastUpdate = time.Now()
}

func (s *Storage) GetInfos() map[uint64]*decoder.SecDefMessage {
	result := make(map[uint64]*decoder.SecDefMessage)
	for key, msg := range s.infoMap {
		result[key] = msg
	}
	return result
}
