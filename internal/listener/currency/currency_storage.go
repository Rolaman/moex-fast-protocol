package currency

import (
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/structs"
	"log"
	"sync"
)

type Storage struct {
	subscribedSecIDs map[string]bool

	// Only healthy (without shifts) messages are counted
	lastRptSeqNum   map[string]int32
	lastRptSeqNumMu sync.RWMutex

	// 0 - Bid, 1 - Ask
	secIDOrderBook   map[string]*structs.Book
	secIDOrderBookMu sync.RWMutex

	// Buffer to restore progress after snapshot applied
	restorationIncrementBuffer   map[string][]*decoder.XOLRCURR
	restorationIncrementBufferMu sync.RWMutex

	bufferedRptSeqs   map[string]map[int32]bool
	bufferedRptSeqsMu sync.Mutex

	isHealthStateBySecID   map[string]*[2]bool
	isHealthStateBySecIDMu sync.Mutex

	unhealthyCount   int // TODO: do on atomics
	unhealthyCountMu sync.RWMutex

	depth uint32

	debugStockOrderBook map[string]*decoder.WOLSCURR
}

func NewStorage(subscribedSecIDs map[string]bool) *Storage {
	return &Storage{
		subscribedSecIDs: subscribedSecIDs,

		lastRptSeqNum:              make(map[string]int32),
		secIDOrderBook:             make(map[string]*structs.Book),
		restorationIncrementBuffer: make(map[string][]*decoder.XOLRCURR),
		bufferedRptSeqs:            make(map[string]map[int32]bool),
		isHealthStateBySecID:       make(map[string]*[2]bool),
		depth:                      5,
		debugStockOrderBook:        make(map[string]*decoder.WOLSCURR),
	}
}

func (s *Storage) AddSnapshot(message *decoder.WOLSCURR) {
	s.debugStockOrderBook[message.Symbol] = message
}

func (s *Storage) AddToBuffer(entry *decoder.XOLRCURRGroupMDEntry) {
	// TODO: implement it
}

func (s *Storage) HandleEntry(entry *decoder.XOLRCURRGroupMDEntry) {
	// TODO: implement it
}

func (s *Storage) AddIncrement(flow flow, message *decoder.XOLRCURR) {
	if message == nil {
		log.Printf("Increment message is nil")
		return
	}
	for _, entry := range message.GroupMDEntries {
		if !s.subscribedSecIDs[entry.Symbol] {
			continue
		}

		s.lastRptSeqNumMu.Lock()
		lastRptSeqNum := s.lastRptSeqNum[entry.Symbol]
		if lastRptSeqNum+1 != entry.RptSeq {
			s.lastRptSeqNumMu.Unlock()
			if lastRptSeqNum+1 > entry.RptSeq {
				continue
			}

			s.isHealthStateBySecIDMu.Lock()
			s.bufferedRptSeqsMu.Lock()
			if _, ok := s.bufferedRptSeqs[entry.Symbol]; !ok {
				s.bufferedRptSeqs[entry.Symbol] = map[int32]bool{}
			}
			if !s.bufferedRptSeqs[entry.Symbol][entry.RptSeq] {
				s.AddToBuffer(entry)
				s.bufferedRptSeqs[entry.Symbol][entry.RptSeq] = true
			}
			s.bufferedRptSeqsMu.Unlock()

			health, ok := s.isHealthStateBySecID[entry.Symbol]
			if !ok {
				health = &[2]bool{
					true,
					true,
				}
				s.isHealthStateBySecID[entry.Symbol] = health
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
			}
			s.isHealthStateBySecIDMu.Unlock()
			continue
		}
		s.isHealthStateBySecIDMu.Lock()
		if _, ok := s.isHealthStateBySecID[entry.Symbol]; !ok {
			s.isHealthStateBySecID[entry.Symbol] = &[2]bool{}
		}
		s.isHealthStateBySecID[entry.Symbol][flow] = true
		s.isHealthStateBySecIDMu.Unlock()
		s.lastRptSeqNum[entry.Symbol] = entry.RptSeq
		s.HandleEntry(entry)
		s.lastRptSeqNumMu.Unlock()
	}
}

func (s *Storage) GetAllCurrencyOrderBook() map[string]*decoder.WOLSCURR {
	return s.debugStockOrderBook
}

func (s *Storage) GetStockOrderBook(symbol string) *structs.Book {
	return s.secIDOrderBook[symbol]
}
