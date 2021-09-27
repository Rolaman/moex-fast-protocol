package currency

import (
	"fmt"
	"github.com/kdt-wolf/moex-fast/internal/config"
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"github.com/kdt-wolf/moex-fast/internal/structs"
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

	return &Listener{
		incrementClients: [2]*protocol.Client{clientIncA, clientIncB},
		snapshotClients:  [2]*protocol.Client{clientSnapshotA, clientSnapshotB},
		storage:          NewStorage(symbols),
	}, nil
}

type Listener struct {
	incrementClients [2]*protocol.Client
	snapshotClients  [2]*protocol.Client
	storage          *Storage
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
	var currWOLSFOND *decoder.WOLSCURR
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
				if decoded.WOLSCURR != nil {
					log.Printf("Currency snap %+v: %+v", decoded.WOLSCURR, decoded.WOLSCURR.GroupMDEntries[0])
				}

				if !ok || decoded.WOLSCURR == nil {
					continue
				}
				if currWOLSFOND == nil {
					currWOLSFOND = decoded.WOLSCURR
				} else {
					currWOLSFOND.GroupMDEntries = append(currWOLSFOND.GroupMDEntries, decoded.WOLSCURR.GroupMDEntries...)
				}
				if decoded.WOLSCURR.LastFragment > 1 {
					panic(fmt.Sprintf("Wrong last fragment %d", decoded.WOLSCURR.LastFragment))
				}
				if decoded.WOLSCURR.LastFragment == 0 {
					continue
				}

				s.storage.AddSnapshot(currWOLSFOND)
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
				case decoded.XOLRCURR != nil:
					s.storage.AddIncrement(flow(i), decoded.XOLRCURR)
				case decoded.Heartbeat != nil:
					log.Printf("(Increment) Heartbeat, sending time: %d", decoded.Heartbeat.SendingTime)
				case decoded.TradingSessionStatus != nil:
					log.Printf("(Increment) Trading session status: %d", decoded.TradingSessionStatus.TradSesStatus)
				case decoded.SecurityDefinition != nil:
					log.Printf("(Increment) Security definition, id: %s, name: %s", decoded.SecurityDefinition.SecurityID, decoded.SecurityDefinition.SecurityDesc)
				}
			}
		}()
	}
	wg.Wait()

}

func (s *Listener) GetStockOrderBook(symbol string) *structs.Book {
	return s.storage.GetStockOrderBook(symbol)
}

func (s *Listener) GetAvailableStockBook() map[string]*decoder.WOLSCURR {
	return s.storage.GetAllCurrencyOrderBook()
}
