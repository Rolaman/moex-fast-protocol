package currency_info

import (
	"github.com/kdt-wolf/moex-fast/internal/config"
	"github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"log"
	"sync"
	"time"
)

type Listener struct {
	client  *protocol.Client
	wg      *sync.WaitGroup
	storage *Storage
}

func New(options *config.ClientOptions) (*Listener, error) {
	client, err := protocol.NewClient(
		&protocol.Options{
			Group:  options.GroupIP,
			Source: options.SourceIP,
			Port:   options.Port,
		})
	if err != nil {
		return nil, err
	}
	return &Listener{
		client:  client,
		storage: NewStorage(),
	}, nil
}

func (s *Listener) Launch() {
	s.wg = new(sync.WaitGroup)
	log.Println("Start listening udp groups for stock info")
	s.wg.Add(1)
	go s.listenSnapshots()
	s.wg.Wait()
}

func (s *Listener) listenSnapshots() {
	for {
		if !s.client.IsConnected {
			if err := s.client.Connect(); err != nil {
				log.Println("Error connecting info snap", err.Error())
				continue
			}
		}
		bytes, err := s.client.ReadNext()
		if err != nil {
			log.Println(err)
			continue
		}

		decoded, ok := stock.Decode(bytes)
		if !ok || decoded.SecurityDefinition == nil {
			continue
		}
		if decoded.SecurityDefinition == nil {
			continue
		}
		oldLen := len(s.storage.GetInfos())
		s.storage.Add(decoded.SecurityDefinition)
		if oldLen == len(s.storage.GetInfos()) {
			time.Sleep(time.Second * 10)
		}

	}
}

func (s *Listener) Stop() {
	log.Println("Stop listening udp groups for stock info")
	s.wg.Done()
}

func (s *Listener) GetInfos() map[string]*stock.SecurityDefinition {
	storage := s.storage
	return storage.GetInfos()
}
