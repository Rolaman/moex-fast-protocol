package stocklistener

import (
	"fmt"
	decoder "github.com/kdt-wolf/moex-fast/internal/decoder/stock"
	"log"
	"strings"
)

var logStocks = map[string]bool{
	"BANE": true,
}

func LogSnap(snap *decoder.WOLSFOND, b []byte) {
	if !contains(snap.Symbol) {
		return
	}
	if snap.RptSeq == 0 && snap.GroupMDEntries[0].MDEntryType == "J" {
		return
	}
	for _, entry := range snap.GroupMDEntries {
		log.Printf("%v %v (SnapEntry:WOLSFOND)- rptSeq %v, type %v, price %v, size %v",
			snap.Symbol, snap.TradingSessionID, snap.RptSeq,
			entry.MDEntryType, entry.MDEntryPx, entry.MDEntrySize)
	}
}

func LogInc(inc *decoder.XOLRFOND, b []byte) {
	for _, entry := range inc.GroupMDEntries {
		if !contains(entry.Symbol) {
			return
		}
		log.Printf("%v %v (IncEntry:WOLSFOND)- rptSeq %v, type %v, action: %v, price %v, size %v",
			entry.Symbol, entry.TradingSessionID, entry.RptSeq,
			entry.MDEntryType, entry.MDUpdateAction, entry.MDEntryPx, entry.MDEntrySize)
	}
}

func LogBytes(b []byte) {
	builder := new(strings.Builder)
	for i, el := range b {
		if b[i] == 0 && len(b) > i+2 && b[i+1] == 0 && b[i+2] == 0 {
			break
		}
		builder.WriteString(fmt.Sprintf("%x", el))
		builder.WriteString(" ")
	}
	log.Printf(builder.String())
}

func contains(e string) bool {
	_, contain := logStocks[e]
	return contain
}
