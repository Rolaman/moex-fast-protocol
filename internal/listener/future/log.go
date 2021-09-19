package futurelistener

import (
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/future"
	"log"
)

var LogFutures = map[uint64]bool{
	// Hardcoded value to debug on one instrument
	752864: true,
}

func LogSnap(snap *decoder.SnapMessage, msg string) {
	if !contains(snap.SecurityID) {
		return
	}
	for _, entry := range snap.Entries {
		log.Printf("%v (SnapEntry:%v)- rptSeq %v, type %v, price %v, size %v, lvl %v",
			snap.SecurityID, msg, snap.RptSeq,
			entry.MDEntryType, entry.MDEntryPx, entry.MDEntrySize, entry.MDPriceLevel)
	}
}

func LogInc(inc *decoder.IncRefMessage, storage *Storage) {
	for _, entry := range inc.Entries {
		if !contains(entry.SecurityID) {
			return
		}
		log.Printf("%v (IncEntry)- rptSeq %v, action %v, type %v, priceLvl %v, price %v, size %v",
			entry.SecurityID, entry.RptSeq, entry.MDUpdateAction, entry.MDEntryType,
			entry.MDPriceLevel, entry.MDEntryPx, entry.MDEntrySize)
		book := storage.GetOrderBook(entry.SecurityID)
		if book == nil {
			return
		}
		log.Printf("Current bids")
		bookRows := book.Bid
		for _, row := range bookRows {
			log.Printf("Row %+v", row)
		}
		log.Printf("Current asks")
		bookRows2 := book.Ask
		for _, row := range bookRows2 {
			log.Printf("Row %+v", row)
		}
	}
}

func LogHandleIncEntry(entry *decoder.IncRefEntry, storage *Storage) {
	if !contains(entry.SecurityID) {
		return
	}
	log.Printf("%v (IncEntry:handle)- rptSeq %v, action %v, type %v, priceLvl %v, price %v, size %v",
		entry.SecurityID, entry.RptSeq, entry.MDUpdateAction, entry.MDEntryType,
		entry.MDPriceLevel, entry.MDEntryPx, entry.MDEntrySize)
	book := storage.GetOrderBook(entry.SecurityID)
	if book == nil {
		return
	}
	log.Printf("Current bids")
	bookRows := book.Bid
	for _, row := range bookRows {
		log.Printf("Row %+v", row)
	}
	log.Printf("Current asks")
	bookRows2 := book.Ask
	for _, row := range bookRows2 {
		log.Printf("Row %+v", row)
	}
}

func LogIncError(entry *decoder.IncRefEntry, lastRptSeqNum uint32) {
	if !contains(entry.SecurityID) {
		return
	}
	log.Printf("%v (ErrorInc)- blank message between %d and %d", entry.SecurityID,
		lastRptSeqNum, entry.RptSeq)
}

func LogBufferRestore(entry *decoder.IncRefEntry) {
	if !contains(entry.SecurityID) {
		return
	}
	log.Printf("%v (RestoreBuffer)- rptSeq %v, action %v, type %v, priceLvl %v, price %v, size %v",
		entry.SecurityID, entry.RptSeq, entry.MDUpdateAction, entry.MDEntryType,
		entry.MDPriceLevel, entry.MDEntryPx, entry.MDEntrySize)
}

func LogLastRptSeq(secId uint64, rptSeq uint32) {
	if !contains(secId) {
		return
	}
	log.Printf("%v (LastRptSeq)- rptSeq %v)", secId, rptSeq)
}

func contains(e uint64) bool {
	_, contain := LogFutures[e]
	return contain
}
