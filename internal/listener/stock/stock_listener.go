package stocklistener

import (
	"fmt"
	"github.com/kdt-wolf/moex-fast/internal/config"
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"github.com/kdt-wolf/moex-fast/internal/structs"
	"github.com/kdt-wolf/moex-fast/internal/web"
	"log"
	"sync"
)

type Listener struct {
	incrementClients [2]*protocol.Client
	snapshotClients  [2]*protocol.Client
	storage          *Storage
	subscribedSecIDs map[string]bool
}

func New(options *config.InstrumentOption, symbols map[string]bool) (*Listener, error) {
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
		return nil, err
	}

	recoverer := &Recoverer{
		options: options.RecoveryOptions,
	}

	return &Listener{
		incrementClients: [2]*protocol.Client{clientIncA, clientIncB},
		snapshotClients:  [2]*protocol.Client{clientSnapshotA, clientSnapshotB},
		storage:          NewStorage(symbols, recoverer),
	}, nil
}

func (s *Listener) Launch() {
	var wg sync.WaitGroup
	log.Println("Start listening udp groups")
	wg.Add(1)
	go s.ListenSnapshots()
	go s.ListenIncrements()
	wg.Wait()
}

func (s *Listener) ListenSnapshots() {
	wg := sync.WaitGroup{}
	var currWOLSFOND *decoder.WOLSFOND
	for _, client := range s.snapshotClients {
		wg.Add(1)
		client := client

		log.Println("starting listening snap")
		go func() {
			for {
				if !client.IsConnected {
					if err := client.Connect(); err != nil {
						log.Println("err connecting snap", err.Error())
						continue
					}
				}

				bytes, err := client.ReadNext()
				if err != nil {
					log.Println(err)
					continue
				}

				decoded, ok := decoder.Decode(bytes)
				if !ok || decoded.WOLSFOND == nil {
					continue
				}
				//skip empty snapshots
				if decoded.WOLSFOND.RptSeq == 0 && decoded.WOLSFOND.GroupMDEntries[0].MDEntryType == "J" {
					continue
				}

				//LogSnap(decoded.WOLSFOND, bytes)
				if !s.subscribedSecIDs[decoded.WOLSFOND.Symbol] {
					continue
				}
				if currWOLSFOND == nil {
					currWOLSFOND = decoded.WOLSFOND
				} else {
					currWOLSFOND.GroupMDEntries = append(currWOLSFOND.GroupMDEntries, decoded.WOLSFOND.GroupMDEntries...)
				}
				s.storage.PushSnapshot(currWOLSFOND)
				if decoded.WOLSFOND.LastFragment > 1 {
					panic(fmt.Sprintf("errored last fragment %d", decoded.WOLSFOND.LastFragment))
				}
				if decoded.WOLSFOND.LastFragment == 0 {
					continue
				}
				currWOLSFOND = nil
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

		log.Println("Starting listening inc " + client.String())
		go func() {
			for {
				bytes, err := client.ReadNextN()
				if err != nil {
					log.Println(err)
					continue
				}
				for _, b := range bytes {
					s.ProcessIncrementBytes(i, b)
				}
			}
		}()
	}
	wg.Wait()
}

func (s *Listener) ProcessIncrementBytes(i int, bytes []byte) {
	decoded, ok := decoder.Decode(bytes)
	if !ok || decoded.XOLRFOND == nil {
		return
	}
	LogInc(decoded.XOLRFOND, bytes)
	switch {
	case decoded.XOLRFOND != nil:
		s.storage.PushIncrement(decoded.XOLRFOND)
	default:
		log.Printf("Not WOLSFOND %+v", decoded)
	}
}

func (s *Listener) GetStockOrderBook(symbol string, session string) *web.BookView {
	book := s.storage.GetStockOrderBook(structs.SecurityId{
		Symbol:    symbol,
		SessionId: session,
	})
	return web.ToView(book)
}

func (s *Listener) GetAvailableStockBook() map[string]decoder.WOLSFOND {
	return s.storage.GetAllStockOrderBook()
}
