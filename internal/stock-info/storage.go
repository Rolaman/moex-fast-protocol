package stockinfo

import (
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"time"
)

type Storage struct {
	LastUpdate time.Time
	infoMap    map[string]*decoder.SecurityDefinition
}

func NewStorage() *Storage {
	return &Storage{
		LastUpdate: time.Now(),
		infoMap:    map[string]*decoder.SecurityDefinition{},
	}
}

func (s *Storage) Add(msg *decoder.SecurityDefinition) {
	s.infoMap[msg.Symbol] = msg
	s.LastUpdate = time.Now()
}

func (s *Storage) GetInfos() map[string]*decoder.SecurityDefinition {
	result := make(map[string]*decoder.SecurityDefinition)
	for key, msg := range s.infoMap {
		result[key] = msg
	}
	return result
}
