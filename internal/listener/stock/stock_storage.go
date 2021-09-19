package stocklistener

import (
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"github.com/kdt-wolf/moex-fast/internal/structs"
	"github.com/shopspring/decimal"
	"log"
	"sync"
)

type Storage struct {
	recoverer        *Recoverer
	subscribedSecIDs map[string]bool

	globalLastIncrementMsgNum uint32
	recovererProcessing       bool
	recovererMsgNums          map[uint32]bool
	// Only healthy (without shifts) messages are counted
	lastRptSeqNum   map[structs.SecurityId]int32
	lastRptSeqNumMu map[structs.SecurityId]*sync.RWMutex

	// 0 - Bid, 1 - Ask
	secIDOrderBook   map[structs.SecurityId]*structs.Book
	secIDOrderBookMu map[structs.SecurityId]*sync.RWMutex

	// Buffer to restore progress after snapshot applied
	restorationIncrementBuffer   map[structs.SecurityId][]*decoder.XOLRFONDGroupMDEntry
	restorationIncrementBufferMu sync.RWMutex

	bufferedRptSeqs map[structs.SecurityId]map[int32]bool

	isHealthStateBySecID map[structs.SecurityId]bool

	snapshotHandler  chan *SnapshotMessage
	incrementHandler chan *IncrementMessage

	// Fetch snapshot only once, use recovery for next missed messages
	snapshotFetchedBySecID   map[structs.SecurityId]bool
	snapshotFetchedBySecIDMu map[structs.SecurityId]*sync.RWMutex

	debugStockOrderBook map[string]decoder.WOLSFOND
}

type SnapshotMessage struct {
	Msg *decoder.WOLSFOND
}

type IncrementMessage struct {
	Msg      *decoder.XOLRFONDGroupMDEntry
	Recovery bool
	MsgNum   uint32
}

func NewStorage(subscribedSecIDs map[string]bool, recoverer *Recoverer) *Storage {
	return &Storage{
		recoverer:        recoverer,
		subscribedSecIDs: subscribedSecIDs,

		lastRptSeqNum:              make(map[structs.SecurityId]int32),
		secIDOrderBook:             make(map[structs.SecurityId]*structs.Book),
		restorationIncrementBuffer: make(map[structs.SecurityId][]*decoder.XOLRFONDGroupMDEntry),
		bufferedRptSeqs:            make(map[structs.SecurityId]map[int32]bool),
		isHealthStateBySecID:       make(map[structs.SecurityId]bool),
		debugStockOrderBook:        make(map[string]decoder.WOLSFOND),

		lastRptSeqNumMu:  make(map[structs.SecurityId]*sync.RWMutex),
		secIDOrderBookMu: make(map[structs.SecurityId]*sync.RWMutex),
		recovererMsgNums: map[uint32]bool{},

		snapshotFetchedBySecID:   make(map[structs.SecurityId]bool),
		snapshotFetchedBySecIDMu: make(map[structs.SecurityId]*sync.RWMutex),
	}
}

func (s *Storage) Launch() {
	go s.processSnapshots()
	go s.processIncrements()
}

func (s *Storage) PushSnapshot(message *decoder.WOLSFOND) {
	if !s.subscribedSecIDs[message.Symbol] {
		return
	}
	msg := &SnapshotMessage{
		Msg: message,
	}
	s.snapshotHandler <- msg
}

func (s *Storage) PushIncrement(message *decoder.XOLRFOND) {
	if message.MsgSeqNum != s.globalLastIncrementMsgNum+1 {
		go func() {
			s.recovererProcessing = true
			for i := s.globalLastIncrementMsgNum + 1; i < message.MsgSeqNum; i++ {
				s.recovererMsgNums[i] = true
			}
			missed := s.recoverer.AddMissed(
				"OLR",
				s.globalLastIncrementMsgNum,
				message.MsgSeqNum-s.globalLastIncrementMsgNum-1,
			)
			for _, entry := range missed {
				for _, mdEntry := range entry.GroupMDEntries {
					if !s.subscribedSecIDs[mdEntry.Symbol] {
						return
					}
					s.incrementHandler <- &IncrementMessage{
						Msg:      mdEntry,
						Recovery: true,
						MsgNum:   entry.MsgSeqNum,
					}
				}
			}
			s.recovererProcessing = false
		}()
	}
	for _, entry := range message.GroupMDEntries {
		if !s.subscribedSecIDs[entry.Symbol] {
			return
		}
		s.incrementHandler <- &IncrementMessage{
			Msg: entry,
		}
	}
}

func (s *Storage) AddIncrement(entry *decoder.XOLRFONDGroupMDEntry) {
	if !s.subscribedSecIDs[entry.Symbol] {
		return
	}

	if s.lastRptSeqNumMu[IncId(entry)] == nil {
		s.lastRptSeqNumMu[IncId(entry)] = &sync.RWMutex{}
	}

	s.lastRptSeqNumMu[IncId(entry)].Lock()
	lastRptSeqNum := s.lastRptSeqNum[IncId(entry)]
	if lastRptSeqNum+1 != entry.RptSeq {
		s.lastRptSeqNumMu[IncId(entry)].Unlock()
		if lastRptSeqNum+1 > entry.RptSeq {
			return
		}

		if _, ok := s.bufferedRptSeqs[IncId(entry)]; !ok {
			s.bufferedRptSeqs[IncId(entry)] = map[int32]bool{}
		}
		if !s.bufferedRptSeqs[IncId(entry)][entry.RptSeq] {
			s.AddToBuffer(entry)
			s.bufferedRptSeqs[IncId(entry)][entry.RptSeq] = true
		}

		// TODO: change health monitoring, with recovery
		s.isHealthStateBySecID[IncId(entry)] = false
		return
	}
	if _, ok := s.isHealthStateBySecID[IncId(entry)]; !ok {
		s.isHealthStateBySecID[IncId(entry)] = false
	}
	s.isHealthStateBySecID[IncId(entry)] = true
	s.lastRptSeqNum[IncId(entry)] = entry.RptSeq
	s.HandleEntry(entry)
	s.lastRptSeqNumMu[IncId(entry)].Unlock()
}

func (s *Storage) AddRecoveryIncrement(entry *decoder.XOLRFONDGroupMDEntry, msgNum uint32) {
	if !s.subscribedSecIDs[entry.Symbol] {
		return
	}
	id := IncId(entry)
	if _, ok := s.bufferedRptSeqs[id]; !ok {
		s.bufferedRptSeqs[id] = map[int32]bool{}
	}
	if !s.bufferedRptSeqs[id][entry.RptSeq] {
		s.AddToBuffer(entry)
		s.bufferedRptSeqs[id][entry.RptSeq] = true
	}
	if s.checkBufferReady(id) {
		s.RestoreFromBuffer(id)
	}

	// TODO: change health monitoring, with recovery
	s.isHealthStateBySecID[id] = false
	return
}

func (s *Storage) AddSnapshot(message *decoder.WOLSFOND) {
	if !s.subscribedSecIDs[message.Symbol] {
		return
	}
	s.debugStockOrderBook[message.Symbol] = *message
	s.secIDOrderBookMu[SnapId(message)] = &sync.RWMutex{}
	if s.lastRptSeqNumMu[SnapId(message)] == nil {
		s.lastRptSeqNumMu[SnapId(message)] = &sync.RWMutex{}
	}
	s.lastRptSeqNumMu[SnapId(message)].Lock()
	if s.isHealthStateBySecID[SnapId(message)] {
		s.lastRptSeqNumMu[SnapId(message)].Unlock()
		return
	}

	// Skip if it is still a lag between latest snapshot and earliest increment
	if message.RptSeq <= s.lastRptSeqNum[SnapId(message)] {
		s.lastRptSeqNumMu[SnapId(message)].Unlock()
		log.Printf("Snapshot skipped, snap rpt seq: %d, last rpt seq num: %d",
			message.RptSeq, s.lastRptSeqNum[SnapId(message)])
		return
	}

	s.lastRptSeqNum[SnapId(message)] = message.RptSeq
	// Fulfil from snapshot and buffer, make ok flag
	s.Restore(message)
	s.snapshotFetchedBySecID[SnapId(message)] = true
	s.lastRptSeqNumMu[SnapId(message)].Unlock()
}

func (s *Storage) Restore(message *decoder.WOLSFOND) {
	s.secIDOrderBookMu[SnapId(message)].Lock()
	orderBook := s.snapMessageToOrderBook(message)
	s.secIDOrderBook[SnapId(message)] = orderBook
	s.secIDOrderBookMu[SnapId(message)].Unlock()
}

func (s *Storage) snapMessageToOrderBook(message *decoder.WOLSFOND) *structs.Book {
	orderBook := &structs.Book{
		Id:   message.Symbol,
		Type: structs.Stock,
		Ask:  map[decimal.Decimal]*structs.BookRow{},
		Bid:  map[decimal.Decimal]*structs.BookRow{},
	}
	for _, entry := range message.GroupMDEntries {
		switch entry.MDEntryType {
		case "0":
			orderBook.Bid[entry.MDEntryPx] = &structs.BookRow{
				Quantity: entry.MDEntrySize,
				Price:    entry.MDEntryPx,
			}
		case "1":
			orderBook.Ask[entry.MDEntryPx] = &structs.BookRow{
				Quantity: entry.MDEntrySize,
				Price:    entry.MDEntryPx,
			}
		default:
			log.Println("specific entry type", entry)
			continue
		}
	}
	return orderBook
}

func (s *Storage) GetStockOrderBook(id structs.SecurityId) *structs.Book {
	return s.secIDOrderBook[id]
}

func (s *Storage) GetAllStockOrderBook() map[string]decoder.WOLSFOND {
	return s.debugStockOrderBook
}

func (s *Storage) AddToBuffer(entry *decoder.XOLRFONDGroupMDEntry) {
	s.restorationIncrementBufferMu.Lock()
	s.restorationIncrementBuffer[IncId(entry)] = append(s.restorationIncrementBuffer[IncId(entry)], entry)
	s.restorationIncrementBufferMu.Unlock()
}

func (s *Storage) HandleEntry(entry *decoder.XOLRFONDGroupMDEntry) {
	switch entry.MDUpdateAction {
	case 0:
		s.AddIncrementNew(entry)
	case 1:
		s.AddIncrementChange(entry)
	case 2:
		s.AddIncrementDelete(entry)
	default:
		log.Printf("%v (ERROR) Unknown action %v", entry.Symbol, entry.MDUpdateAction)
	}
}

func (s *Storage) AddIncrementNew(entry *decoder.XOLRFONDGroupMDEntry) {
	s.secIDOrderBookMu[IncId(entry)].Lock()

	switch entry.MDEntryType {
	case "0":
		bookRows := s.secIDOrderBook[IncId(entry)].Bid
		bookRows[entry.MDEntryPx] = &structs.BookRow{
			Quantity: entry.MDEntrySize,
			Price:    entry.MDEntryPx,
		}
	case "1":
		bookRows := s.secIDOrderBook[IncId(entry)].Ask
		bookRows[entry.MDEntryPx] = &structs.BookRow{
			Quantity: entry.MDEntrySize,
			Price:    entry.MDEntryPx,
		}
	default:
		log.Println("Specific entry type", entry)
	}

	s.secIDOrderBookMu[IncId(entry)].Unlock()
}

func (s *Storage) AddIncrementChange(entry *decoder.XOLRFONDGroupMDEntry) {
	s.secIDOrderBookMu[IncId(entry)].Lock()

	switch entry.MDEntryType {
	case "0":
		bookRows := s.secIDOrderBook[IncId(entry)].Bid
		bookRows[entry.MDEntryPx] = &structs.BookRow{
			Quantity: entry.MDEntrySize,
			Price:    entry.MDEntryPx,
		}
	case "1":
		bookRows := s.secIDOrderBook[IncId(entry)].Ask
		bookRows[entry.MDEntryPx] = &structs.BookRow{
			Quantity: entry.MDEntrySize,
			Price:    entry.MDEntryPx,
		}
	default:
		log.Println("Specific entry type", entry)
	}

	s.secIDOrderBookMu[IncId(entry)].Unlock()
}

func (s *Storage) AddIncrementDelete(entry *decoder.XOLRFONDGroupMDEntry) {
	s.secIDOrderBookMu[IncId(entry)].Lock()

	switch entry.MDEntryType {
	case "0":
		delete(s.secIDOrderBook[IncId(entry)].Bid, entry.MDEntryPx)
	case "1":
		delete(s.secIDOrderBook[IncId(entry)].Ask, entry.MDEntryPx)
	default:
		log.Println("Specific entry type", entry)
	}

	s.secIDOrderBookMu[IncId(entry)].Unlock()
}

func (s *Storage) processSnapshots() {
	for {
		snapMessage := <-s.snapshotHandler
		msg := snapMessage.Msg
		s.AddSnapshot(msg)
	}
}

func (s *Storage) processIncrements() {
	for {
		incrementMessage := <-s.incrementHandler
		msg := incrementMessage.Msg
		s.AddIncrement(msg)
		if incrementMessage.Recovery {
			s.recovererMsgNums[incrementMessage.MsgNum] = false
		}
	}
}

func (s *Storage) checkBufferReady(id structs.SecurityId) bool {
	if len(s.bufferedRptSeqs[id]) == 0 {
		return false
	}
	var minNum int32 = 0
	var maxNum int32 = 0
	for seq := range s.bufferedRptSeqs[id] {
		if seq < minNum {
			minNum = seq
		}
		if seq > maxNum {
			maxNum = seq
		}
	}
	return maxNum-minNum == int32(len(s.bufferedRptSeqs[id]))
}

func (s *Storage) RestoreFromBuffer(id structs.SecurityId) {
	// TODO: implement restoring
}

func SnapId(message *decoder.WOLSFOND) structs.SecurityId {
	return structs.SecurityId{
		Symbol:    message.Symbol,
		SessionId: message.TradingSessionID,
	}
}

func IncId(message *decoder.XOLRFONDGroupMDEntry) structs.SecurityId {
	return structs.SecurityId{
		Symbol:    message.Symbol,
		SessionId: message.TradingSessionID,
	}
}
