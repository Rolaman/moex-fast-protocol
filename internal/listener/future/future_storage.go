package futurelistener

import (
	"github.com/kdt-wolf/moex-fast/internal/structs"
	"github.com/shopspring/decimal"
	"log"
	"sync"

	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/future"
)

type Storage struct {
	depth uint32

	subscribedSecIDs map[uint64]bool

	// only healthy (without shifts) messages are counted
	lastRptSeqNum   map[uint64]uint32
	lastRptSeqNumMu sync.RWMutex

	// 0 - Bid, 1 - Ask
	secIDOrderBook   map[uint64]*structs.SortedBook
	secIDOrderBookMu sync.RWMutex

	// buffer to restore progress after snapshot applied
	restorationIncrementBuffer   map[uint64][]*decoder.IncRefEntry
	restorationIncrementBufferMu sync.RWMutex

	bufferedRptSeqs   map[uint64]map[uint32]bool
	bufferedRptSeqsMu sync.Mutex

	isHealthStateBySecID   map[uint64]*[2]bool
	isHealthStateBySecIDMu sync.Mutex

	unhealthyCount   int // TODO: do on atomics
	unhealthyCountMu sync.RWMutex
}

func NewStorage(depth uint32, subscribedSecIDs map[uint64]bool) *Storage {
	s := &Storage{
		depth: depth,

		subscribedSecIDs: subscribedSecIDs,

		lastRptSeqNum:              make(map[uint64]uint32),
		secIDOrderBook:             make(map[uint64]*structs.SortedBook),
		restorationIncrementBuffer: make(map[uint64][]*decoder.IncRefEntry),
		bufferedRptSeqs:            make(map[uint64]map[uint32]bool),
		isHealthStateBySecID:       make(map[uint64]*[2]bool),
	}
	return s
}

func (s *Storage) GetOrderBook(secID uint64) *structs.SortedBook {
	return s.secIDOrderBook[secID]
}

func (s *Storage) GetAvailableBooks() []uint64 {
	keys := make([]uint64, 0, len(s.secIDOrderBook))
	for k := range s.secIDOrderBook {
		keys = append(keys, k)
	}
	return keys
}

func (s *Storage) GetAllBooks() map[uint64]*structs.SortedBook {
	return s.secIDOrderBook
}

func (s *Storage) AddSnapshot(message *decoder.SnapMessage) {
	if message == nil {
		log.Printf("Snapshot message is nil")
		return
	}

	s.lastRptSeqNumMu.Lock()
	s.isHealthStateBySecIDMu.Lock()
	if s.isHealthStateBySecID[message.SecurityID] != nil && (s.isHealthStateBySecID[message.SecurityID][0] || s.isHealthStateBySecID[message.SecurityID][1]) {
		s.isHealthStateBySecIDMu.Unlock()
		s.lastRptSeqNumMu.Unlock()
		return
	}
	s.isHealthStateBySecIDMu.Unlock()

	// Skip if it is still a lag between latest snapshot and earliest increment
	if message.RptSeq <= s.lastRptSeqNum[message.SecurityID] {
		s.lastRptSeqNumMu.Unlock()
		log.Printf("Snapshot skipped, snap rpt seq: %d, last rpt seq num: %d", message.RptSeq, s.lastRptSeqNum[message.SecurityID])
		return
	}
	log.Printf("%v (Snap)- Restore from snapshot, rptSeq %v", message.SecurityID, message.RptSeq)

	s.lastRptSeqNum[message.SecurityID] = message.RptSeq
	// Fulfil from snapshot and buffer, make ok flag
	s.Restore(message)
	s.lastRptSeqNumMu.Unlock()
}

func (s *Storage) AddIncrement(flow flow, message *decoder.IncRefMessage) {
	if message == nil {
		log.Printf("Increment message is nil")
		return
	}
	for _, entry := range message.Entries {

		s.lastRptSeqNumMu.Lock()
		lastRptSeqNum := s.lastRptSeqNum[entry.SecurityID]
		if lastRptSeqNum+1 != entry.RptSeq {
			s.lastRptSeqNumMu.Unlock()
			if lastRptSeqNum+1 > entry.RptSeq {
				continue
			}

			s.isHealthStateBySecIDMu.Lock()
			s.bufferedRptSeqsMu.Lock()
			if _, ok := s.bufferedRptSeqs[entry.SecurityID]; !ok {
				s.bufferedRptSeqs[entry.SecurityID] = map[uint32]bool{}
			}
			if !s.bufferedRptSeqs[entry.SecurityID][entry.RptSeq] {
				s.AddToBuffer(entry)
				s.bufferedRptSeqs[entry.SecurityID][entry.RptSeq] = true
			}
			s.bufferedRptSeqsMu.Unlock()

			health, ok := s.isHealthStateBySecID[entry.SecurityID]
			if !ok {
				health = &[2]bool{
					true,
					true,
				}
				s.isHealthStateBySecID[entry.SecurityID] = health
			}
			if !health[flow] {
				s.isHealthStateBySecIDMu.Unlock()
				continue
			}
			health[flow] = false
			if !health[flow.oppositeFlow()] {
				s.unhealthyCountMu.Lock()
				s.unhealthyCount++
				s.unhealthyCountMu.Unlock()
				LogIncError(entry, lastRptSeqNum)
			}
			s.isHealthStateBySecIDMu.Unlock()
			continue
		}
		s.isHealthStateBySecIDMu.Lock()
		if _, ok := s.isHealthStateBySecID[entry.SecurityID]; !ok {
			s.isHealthStateBySecID[entry.SecurityID] = &[2]bool{}
		}
		s.isHealthStateBySecID[entry.SecurityID][flow] = true
		s.isHealthStateBySecIDMu.Unlock()
		s.lastRptSeqNum[entry.SecurityID] = entry.RptSeq
		s.HandleEntry(entry)
		s.lastRptSeqNumMu.Unlock()
	}
}

func (s *Storage) Restore(message *decoder.SnapMessage) {
	s.secIDOrderBookMu.Lock()
	orderBook := s.snapMessageToOrderBook(message)
	s.secIDOrderBook[message.SecurityID] = orderBook
	s.secIDOrderBookMu.Unlock()

	s.AppendFromBuffer(message.SecurityID, message.RptSeq)

	s.unhealthyCountMu.Lock()
	s.unhealthyCount--
	s.unhealthyCountMu.Unlock()
	s.isHealthStateBySecIDMu.Lock()
	if _, ok := s.isHealthStateBySecID[message.SecurityID]; !ok {
		s.isHealthStateBySecID[message.SecurityID] = &[2]bool{}
	}
	s.isHealthStateBySecID[message.SecurityID][0] = true
	s.isHealthStateBySecID[message.SecurityID][1] = true
	s.isHealthStateBySecIDMu.Unlock()
}

func (s *Storage) AddToBuffer(entry *decoder.IncRefEntry) {
	s.restorationIncrementBufferMu.Lock()
	s.restorationIncrementBuffer[entry.SecurityID] = append(s.restorationIncrementBuffer[entry.SecurityID], entry)
	s.restorationIncrementBufferMu.Unlock()
}

func (s *Storage) AppendFromBuffer(secId uint64, snapRptSeq uint32) {
	lastRptSeq := uint32(0)
	s.restorationIncrementBufferMu.Lock()
	s.bufferedRptSeqsMu.Lock()
	for _, entry := range s.restorationIncrementBuffer[secId] { // TODO: handle possible repeated shift
		if entry.RptSeq <= snapRptSeq {
			continue
		}
		LogBufferRestore(entry)
		s.HandleEntry(entry)
		lastRptSeq = entry.RptSeq
	}
	if lastRptSeq != 0 {
		LogLastRptSeq(secId, lastRptSeq)
		s.lastRptSeqNum[secId] = lastRptSeq // Lock in upper level, rewrite
	}
	s.bufferedRptSeqs[secId] = map[uint32]bool{}
	s.bufferedRptSeqsMu.Unlock()
	s.restorationIncrementBuffer[secId] = []*decoder.IncRefEntry{}
	s.restorationIncrementBufferMu.Unlock()
}

func (s *Storage) HandleEntry(entry *decoder.IncRefEntry) {
	switch entry.MDUpdateAction {
	case 0:
		s.AddIncrementNew(entry)
	case 1:
		s.AddIncrementChange(entry)
	case 2:
		s.AddIncrementDelete(entry)
	default:
		log.Printf("%v (ERROR) Unknown action %v", entry.SecurityID, entry.MDUpdateAction)
	}
	LogHandleIncEntry(entry, s)
}

func (s *Storage) AddIncrementNew(entry *decoder.IncRefEntry) {
	if entry.MDPriceLevel-1 > s.depth {
		return
	}

	s.secIDOrderBookMu.Lock()

	switch entry.MDEntryType {
	case "0":
		for i := s.depth - 1; i >= entry.MDPriceLevel; i-- {
			s.secIDOrderBook[entry.SecurityID].Bid[i] = s.secIDOrderBook[entry.SecurityID].Bid[i-1]
		}

		s.secIDOrderBook[entry.SecurityID].Bid[entry.MDPriceLevel-1] = &structs.SortedBookRow{
			Quantity: decimal.NewFromInt(entry.MDEntrySize),
			Price:    entry.MDEntryPx,
		}
	case "1":
		for i := s.depth - 1; i >= entry.MDPriceLevel; i-- {
			s.secIDOrderBook[entry.SecurityID].Ask[i] = s.secIDOrderBook[entry.SecurityID].Bid[i-1]
		}

		s.secIDOrderBook[entry.SecurityID].Ask[entry.MDPriceLevel-1] = &structs.SortedBookRow{
			Quantity: decimal.NewFromInt(entry.MDEntrySize),
			Price:    entry.MDEntryPx,
		}
	default:
		log.Println("specific entry type", entry)
		s.secIDOrderBookMu.Unlock()
		return
	}

	s.secIDOrderBookMu.Unlock()
}

func (s *Storage) AddIncrementChange(entry *decoder.IncRefEntry) {
	if entry.MDPriceLevel-1 > s.depth {
		return
	}

	s.secIDOrderBookMu.Lock()
	switch entry.MDEntryType {
	case "0":
		s.secIDOrderBook[entry.SecurityID].Bid[entry.MDPriceLevel-1] = &structs.SortedBookRow{
			Quantity: decimal.NewFromInt(entry.MDEntrySize),
			Price:    entry.MDEntryPx,
		}
	case "1":
		s.secIDOrderBook[entry.SecurityID].Ask[entry.MDPriceLevel-1] = &structs.SortedBookRow{
			Quantity: decimal.NewFromInt(entry.MDEntrySize),
			Price:    entry.MDEntryPx,
		}
	default:
		log.Println("specific entry type", entry)
		s.secIDOrderBookMu.Unlock()
		return
	}

	s.secIDOrderBookMu.Unlock()
}

func (s *Storage) AddIncrementDelete(entry *decoder.IncRefEntry) {
	if entry.MDPriceLevel-1 > s.depth {
		return
	}

	s.secIDOrderBookMu.Lock()
	switch entry.MDEntryType {
	case "0":
		for i := entry.MDPriceLevel - 1; i < s.depth-1; i++ {
			s.secIDOrderBook[entry.SecurityID].Bid[i] = s.secIDOrderBook[entry.SecurityID].Bid[i+1]
		}
		s.secIDOrderBook[entry.SecurityID].Bid[s.depth-1] = nil
	case "1":
		for i := entry.MDPriceLevel - 1; i < s.depth-1; i++ {
			s.secIDOrderBook[entry.SecurityID].Ask[i] = s.secIDOrderBook[entry.SecurityID].Ask[i+1]
		}
		s.secIDOrderBook[entry.SecurityID].Ask[s.depth-1] = nil
	default:
		log.Println("specific entry type", entry)
		s.secIDOrderBookMu.Unlock()
		return
	}
	s.secIDOrderBookMu.Unlock()
}

func (s *Storage) snapMessageToOrderBook(message *decoder.SnapMessage) *structs.SortedBook {
	orderBook := &structs.SortedBook{
		Ask:  make([]*structs.SortedBookRow, s.depth),
		Bid:  make([]*structs.SortedBookRow, s.depth),
		Id:   message.Symbol,
		Type: structs.Future,
	}
	for _, entry := range message.Entries {
		switch entry.MDEntryType {
		case "0":
			orderBook.Bid[entry.MDPriceLevel-1] = &structs.SortedBookRow{
				Quantity: decimal.NewFromInt(entry.MDEntrySize),
				Price:    entry.MDEntryPx,
			}
		case "1":
			orderBook.Ask[entry.MDPriceLevel-1] = &structs.SortedBookRow{
				Quantity: decimal.NewFromInt(entry.MDEntrySize),
				Price:    entry.MDEntryPx,
			}
		default:
			log.Println("specific entry type", entry)
			continue
		}
		if entry.MDPriceLevel > s.depth {
			continue
		}
	}
	return orderBook
}
