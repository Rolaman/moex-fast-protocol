package futurelistener

import (
	"fmt"
	"github.com/kdt-wolf/moex-fast/internal/config"
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/future"
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"github.com/kdt-wolf/moex-fast/internal/web"
	"log"
	"sync"
)

type flow byte

func (f flow) oppositeFlow() flow {
	if f == 0 {
		return 1
	}
	return 0
}

type Listener struct {
	incrementClients [2]*protocol.Client
	snapshotClients  [2]*protocol.Client

	storage          *Storage
	currentSnapMutex sync.Mutex
}

func New(options *config.InstrumentOption, secIDs map[uint64]bool) (*Listener, error) {
	clientIncA, err := protocol.NewClient(
		&protocol.Options{
			Group:  options.IncrementalClientAOptions.GroupIP,
			Source: options.IncrementalClientAOptions.SourceIP,
			Port:   options.IncrementalClientAOptions.Port,
		})
	if err != nil {
		return nil, err
	}

	clientIncB, err := protocol.NewClient(
		&protocol.Options{
			Group:  options.IncrementalClientBOptions.GroupIP,
			Source: options.IncrementalClientBOptions.SourceIP,
			Port:   options.IncrementalClientBOptions.Port,
		})
	if err != nil {
		return nil, err
	}

	clientSnapshotA, err := protocol.NewClient(
		&protocol.Options{
			Group:  options.SnapshotClientAOptions.GroupIP,
			Source: options.SnapshotClientAOptions.SourceIP,
			Port:   options.SnapshotClientAOptions.Port,
		})
	if err != nil {
		return nil, err
	}

	clientSnapshotB, err := protocol.NewClient(
		&protocol.Options{
			Group:  options.SnapshotClientBOptions.GroupIP,
			Source: options.SnapshotClientBOptions.SourceIP,
			Port:   options.SnapshotClientBOptions.Port,
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Listener{
		incrementClients: [2]*protocol.Client{clientIncA, clientIncB},
		snapshotClients:  [2]*protocol.Client{clientSnapshotA, clientSnapshotB},
		storage:          NewStorage(options.Depth, secIDs),
	}, nil
}

func (s *Listener) GetOrderBook(secID uint64) *web.BookView {
	return s.storage.GetOrderBook(secID)
}

func (s *Listener) GetAvailableBook() []uint64 {
	return s.storage.GetAvailableBooks()
}

func (s *Listener) Launch() {
	var wg sync.WaitGroup
	log.Println("Start listening udp groups")
	wg.Add(1)
	go s.ListenSnapshots()
	//go s.ListenIncrements()
	wg.Wait()
}

func (s *Listener) ListenSnapshots() {
	wg := sync.WaitGroup{}
	var currSnap *decoder.SnapMessage
	for _, client := range s.snapshotClients {
		wg.Add(1)
		client := client

		log.Println("Starting listening snap")
		go func() {
			for {
				if !client.IsConnected {
					if err := client.Connect(); err != nil {
						log.Println("Err connecting snap", err.Error())
						continue
					}
				}

				bytes, err := client.ReadNext()
				if err != nil {
					log.Println(err)
					continue
				}

				decoded, ok := decoder.Decode(bytes)

				if !ok || decoded.Snap == nil {
					continue
				}

				s.currentSnapMutex.Lock()
				if currSnap == nil {
					currSnap = decoded.Snap
				} else {
					log.Printf("%+v, %+v", currSnap.Entries, decoded.Snap)
					currSnap.Entries = append(currSnap.Entries, decoded.Snap.Entries...)
				}
				if decoded.Snap.LastFragment > 1 {
					panic(fmt.Sprintf("errored last fragment %d", decoded.Snap.LastFragment))
				}
				if decoded.Snap.LastFragment == 0 {
					continue
				}

				s.storage.AddSnapshot(currSnap)
				currSnap = nil
				s.currentSnapMutex.Unlock()
			}
		}()
	}
	wg.Wait()
}

func (s *Listener) ListenIncrements() {
	wg := sync.WaitGroup{}
	for i, client := range s.incrementClients {
		wg.Add(1)
		i := i
		client := client
		if err := client.Connect(); err != nil {
			panic(err)
		}

		log.Println("starting listening inc")
		go func() {
			for {
				bytes, err := client.ReadNext()
				if err != nil {
					log.Println(err)
					continue
				}

				decoded, ok := decoder.Decode(bytes)
				if !ok {
					continue
				}
				switch {
				case decoded.IncRef != nil:
					s.storage.AddIncrement(flow(i), decoded.IncRef)
				case decoded.HeartbeatMessage != nil:
					log.Printf("(Increment) Heartbeat, sending time: %d", decoded.HeartbeatMessage.SendingTime)
				case decoded.TradingSessionStatus != nil:
					log.Printf("(Increment) Trading session status: %d", decoded.TradingSessionStatus.TradSesStatus)
				case decoded.SecurityDefinition != nil:
					log.Printf("(Increment) Security definition, id: %d, name: %s", decoded.SecurityDefinition.SecurityID, decoded.SecurityDefinition.SecurityDesc)
				}
			}
		}()
	}

	wg.Wait()
}
