package stock

import (
	"github.com/kdt-wolf/moex-fast/internal/decoder"
	"log"
)

func Decode(bytes []byte) (*Message, bool) {
	iterator := decoder.NewIterator(bytes)
	_ = decoder.Number(iterator)

	pmap := decoder.PMap(iterator)
	templateId := iterator.NextUInt32(false)
	if pmap != 64 {
		log.Printf("Pmap is not 64, %d for template %d", pmap, templateId)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Can't decode bytes %+v, %+v", bytes, r)
		}
	}()

	switch templateId {
	case 2101:
		return &Message{Logon: DecodeLogon(iterator), TemplateId: templateId}, true
	case 2102:
		return &Message{Logout: DecodeLogout(iterator), TemplateId: templateId}, true
	case 2603:
		return &Message{WGeneric: DecodeWGeneric(iterator), TemplateId: templateId}, true
	case 2104:
		return &Message{XGeneric: DecodeXGeneric(iterator), TemplateId: templateId}, true
	case 2710:
		return &Message{WOLSFOND: DecodeWOLSFOND(iterator), TemplateId: templateId}, true
	case 2610:
		return &Message{WOLSFOND: DecodeWOLSFOND(iterator), TemplateId: templateId}, true
	case 3600:
		return &Message{WOLSCURR: DecodeWOLSCURR(iterator), TemplateId: templateId}, true
	case 2611:
		return &Message{WTLSFOND: DecodeWTLSFOND(iterator), TemplateId: templateId}, true
	case 3601:
		return &Message{WTLSCURR: DecodeWTLSCURR(iterator), TemplateId: templateId}, true
	case 2623:
		return &Message{XMSRFOND: DecodeXMSRFOND(iterator), TemplateId: templateId}, true
	case 3613:
		return &Message{XMSRCURR: DecodeXMSRCURR(iterator), TemplateId: templateId}, true
	case 2720:
		return &Message{XOLRFOND: DecodeXOLRFOND(iterator), TemplateId: templateId}, true
	case 2620:
		return &Message{XOLRFOND: DecodeXOLRFOND(iterator), TemplateId: templateId}, true
	case 3610:
		return &Message{XOLRCURR: DecodeXOLRCURR(iterator), TemplateId: templateId}, true
	case 2621:
		return &Message{XTLRFOND: DecodeXTLRFOND(iterator), TemplateId: templateId}, true
	case 3611:
		return &Message{XTLRCURR: DecodeXTLRCURR(iterator), TemplateId: templateId}, true
	case 2615:
		return &Message{SecurityDefinition: DecodeSecurityDefinition(iterator)}, true
	case 2106:
		return &Message{SecurityStatus: DecodeSecurityStatus(iterator)}, true
	case 2107:
		return &Message{TradingSessionStatus: DecodeTradingSessionStatus(iterator)}, true
	case 2108:
		return &Message{Heartbeat: DecodeHeartbeat(iterator)}, true
	default:
		log.Printf("Return empty message, templateId: %d", templateId)
		log.Printf("Bytes %+v", bytes)
		return &Message{
			TemplateId: templateId,
		}, true
	}
}

func DecodeLogon(iter *decoder.MessageIterator) *Logon {
	return &Logon{
		TargetCompID:     iter.NextString(),
		MsgSeqNum:        iter.NextUInt32(false),
		SendingTime:      iter.NextUInt64(false),
		HeartBtInt:       iter.NextInt32(false),
		Username:         iter.NextString(),
		Password:         iter.NextString(),
		DefaultApplVerID: iter.NextString(),
	}
}

func DecodeLogout(iter *decoder.MessageIterator) *Logout {
	return &Logout{
		TargetCompID: iter.NextString(),
		MsgSeqNum:    iter.NextUInt32(false),
		SendingTime:  iter.NextUInt64(false),
		Text:         iter.NextString(),
	}
}

func DecodeWGeneric(iter *decoder.MessageIterator) *WGeneric {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	tradingSessionID := iter.NextString()
	symbol := iter.NextString()
	lastMsgSeqNumProcessed := iter.NextUInt32(true)
	rptSeq := iter.NextInt32(false)
	lastFragment := iter.NextUInt32(true)
	routeFirst := iter.NextUInt32(true)
	tradSesStatus := iter.NextInt32(true)
	mdSecurityTradingStatus := iter.NextInt32(true)
	auctionIndicator := iter.NextUInt32(true)
	netChgPrevDay := iter.NextDecimal(true)
	priceImprovement := iter.NextDecimal(true)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*WGenericGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		groupMDEntries[i] = DecodeWGenericGroupMDEntry(iter)
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WGeneric", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &WGeneric{
		MsgSeqNum:               msgSeqNum,
		SendingTime:             sendingTime,
		TradingSessionID:        tradingSessionID,
		Symbol:                  symbol,
		LastMsgSeqNumProcessed:  lastMsgSeqNumProcessed,
		RptSeq:                  rptSeq,
		LastFragment:            lastFragment,
		RouteFirst:              routeFirst,
		TradSesStatus:           tradSesStatus,
		MDSecurityTradingStatus: mdSecurityTradingStatus,
		AuctionIndicator:        auctionIndicator,
		NetChgPrevDay:           netChgPrevDay,
		PriceImprovement:        priceImprovement,
		NoMDEntries:             noMDEntries,
		GroupMDEntries:          groupMDEntries,
	}
}

func DecodeWGenericGroupMDEntry(iter *decoder.MessageIterator) *WGenericGroupMDEntry {
	return &WGenericGroupMDEntry{
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		MDEntryDate:         iter.NextUInt32(true),
		MDEntryTime:         iter.NextUInt32(true),
		OrigTime:            iter.NextUInt32(true),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		QuoteCondition:      iter.NextString(),
		TradeCondition:      iter.NextString(),
		OpenCloseSettlFlag:  iter.NextString(),
		OrdType:             iter.NextString(),
		EffectiveTime:       iter.NextUInt32(true),
		StartTime:           iter.NextUInt32(true),
		AccruedInterestAmt:  iter.NextDecimal(true),
		ChgFromWAPrice:      iter.NextDecimal(true),
		ChgOpenInterest:     iter.NextDecimal(true),
		BidMarketSize:       iter.NextDecimal(true),
		AskMarketSize:       iter.NextDecimal(true),
		TotalNumOfTrades:    iter.NextInt32(true),
		TradeValue:          iter.NextDecimal(true),
		Yield:               iter.NextDecimal(true),
		TotalVolume:         iter.NextDecimal(true),
		OfferNbOr:           iter.NextInt32(true),
		BidNbOr:             iter.NextInt32(true),
		ChgFromSettlmnt:     iter.NextDecimal(true),
		SettlDate:           iter.NextUInt32(true),
		SettlDate2:          iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		SumQtyOfBest:        iter.NextUInt32(true),
		OrderSide:           iter.NextString(),
		OrderStatus:         iter.NextString(),
		MinCurrPx:           iter.NextDecimal(true),
		IndexCrossRate:      iter.NextDecimal(true),
		MinCurrPxChgTime:    iter.NextUInt32(true),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		VolumeIndicator:     iter.NextUInt32(true),
		Price:               iter.NextDecimal(true),
		PriceType:           iter.NextInt32(true),
		NominalValue:        iter.NextDecimal(true),
		RepoToPx:            iter.NextDecimal(true),
		BuyBackPx:           iter.NextDecimal(true),
		BuyBackDate:         iter.NextUInt32(true),
		CXFlag:              iter.NextString(),
		TradingSessionSubID: iter.NextString(),
	}
}

func DecodeXGeneric(iter *decoder.MessageIterator) *XGeneric {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*XGenericGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		groupMDEntries[i] = DecodeXGenericGroupMDEntry(iter)
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XGeneric", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XGeneric{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXGenericGroupMDEntry(iter *decoder.MessageIterator) *XGenericGroupMDEntry {
	return &XGenericGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		RptSeq:              iter.NextInt32(true),
		MDEntryDate:         iter.NextUInt32(true),
		OrigTime:            iter.NextUInt32(true),
		SettlDate:           iter.NextUInt32(true),
		SettlDate2:          iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		MDEntryTime:         iter.NextUInt32(true),
		EffectiveTime:       iter.NextUInt32(true),
		StartTime:           iter.NextUInt32(true),
		Symbol:              iter.NextString(),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		QuoteCondition:      iter.NextString(),
		TradeCondition:      iter.NextString(),
		OpenCloseSettlFlag:  iter.NextString(),
		OrdType:             iter.NextString(),
		NetChgPrevDay:       iter.NextDecimal(true),
		PriceImprovement:    iter.NextDecimal(true),
		AccruedInterestAmt:  iter.NextDecimal(true),
		ChgFromWAPrice:      iter.NextDecimal(true),
		ChgOpenInterest:     iter.NextDecimal(true),
		BidMarketSize:       iter.NextDecimal(true),
		AskMarketSize:       iter.NextDecimal(true),
		TotalNumOfTrades:    iter.NextInt32(true),
		TradeValue:          iter.NextDecimal(true),
		Yield:               iter.NextDecimal(true),
		TotalVolume:         iter.NextDecimal(true),
		OfferNbOr:           iter.NextInt32(true),
		BidNbOr:             iter.NextInt32(true),
		ChgFromSettlmnt:     iter.NextDecimal(true),
		SumQtyOfBest:        iter.NextInt32(true),
		OrderSide:           iter.NextString(),
		OrderStatus:         iter.NextString(),
		MinCurrPx:           iter.NextDecimal(true),
		MinCurrPxChgTime:    iter.NextUInt32(true),
		IndexCrossRate:      iter.NextDecimal(true),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		VolumeIndicator:     iter.NextUInt32(true),
		Price:               iter.NextDecimal(true),
		PriceType:           iter.NextInt32(true),
		NominalValue:        iter.NextDecimal(true),
		RepoToPx:            iter.NextDecimal(true),
		BuyBackPx:           iter.NextDecimal(true),
		BuyBackDate:         iter.NextUInt32(true),
		CXFlag:              iter.NextString(),
		RepoTerm:            iter.NextUInt32(true),
		TradingSessionID:    iter.NextString(),
		TradingSessionSubID: iter.NextString(),
	}
}

func DecodeWOLSFOND(iter *decoder.MessageIterator) *WOLSFOND {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastMsgSeqNumProcessed := iter.NextUInt32(true)
	rptSeq := iter.NextInt32(false)
	lastFragment := iter.NextUInt32(true)
	routeFirst := iter.NextUInt32(true)
	tradSesStatus := iter.NextInt32(true)
	tradingSessionID := iter.NextString()
	symbol := iter.NextString()
	mdSecurityTradingStatus := iter.NextInt32(true)
	auctionIndicator := iter.NextUInt32(true)
	noMDEntries := iter.NextUInt32(false)
	lastMDEntry := &WOLSFONDGroupMDEntry{}
	groupMDEntries := make([]*WOLSFONDGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeWOLSFONDGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WOLSFOND", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &WOLSFOND{
		MsgSeqNum:               msgSeqNum,
		SendingTime:             sendingTime,
		LastMsgSeqNumProcessed:  lastMsgSeqNumProcessed,
		RptSeq:                  rptSeq,
		LastFragment:            lastFragment,
		RouteFirst:              routeFirst,
		TradSesStatus:           tradSesStatus,
		TradingSessionID:        tradingSessionID,
		Symbol:                  symbol,
		MDSecurityTradingStatus: mdSecurityTradingStatus,
		AuctionIndicator:        auctionIndicator,
		NoMDEntries:             noMDEntries,
		GroupMDEntries:          groupMDEntries,
	}
}

func DecodeWOLSFONDGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *WOLSFONDGroupMDEntry,
) *WOLSFONDGroupMDEntry {
	return &WOLSFONDGroupMDEntry{
		MDEntryType:         pmap.NextString(iter.NextString, prev.MDEntryType),
		MDEntryID:           iter.NextString(),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		Yield:               pmap.NextDecimal(iter.NextDecimalOptional, prev.Yield),
		OrderStatus:         pmap.NextString(iter.NextString, prev.OrderStatus),
		OrdType:             pmap.NextString(iter.NextString, prev.OrdType),
		TotalVolume:         pmap.NextDecimal(iter.NextDecimalOptional, prev.TotalVolume),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
	}
}

func DecodeWOLSCURR(iter *decoder.MessageIterator) *WOLSCURR {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastMsgSeqNumProcessed := iter.NextUInt32(true)
	rptSeq := iter.NextInt32(false)
	lastFragment := iter.NextUInt32(true)
	routeFirst := iter.NextUInt32(true)
	tradSesStatus := iter.NextInt32(true)
	tradingSessionID := iter.NextString()
	symbol := iter.NextString()
	mdSecurityTradingStatus := iter.NextInt32(true)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*WOLSCURRGroupMDEntry, noMDEntries)
	lastMDEntry := &WOLSCURRGroupMDEntry{}
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeWOLSCURRGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WOLSCURR", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &WOLSCURR{
		MsgSeqNum:               msgSeqNum,
		SendingTime:             sendingTime,
		LastMsgSeqNumProcessed:  lastMsgSeqNumProcessed,
		RptSeq:                  rptSeq,
		LastFragment:            lastFragment,
		RouteFirst:              routeFirst,
		TradSesStatus:           tradSesStatus,
		TradingSessionID:        tradingSessionID,
		Symbol:                  symbol,
		MDSecurityTradingStatus: mdSecurityTradingStatus,
		NoMDEntries:             noMDEntries,
		GroupMDEntries:          groupMDEntries,
	}
}

func DecodeWOLSCURRGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *WOLSCURRGroupMDEntry,
) *WOLSCURRGroupMDEntry {
	return &WOLSCURRGroupMDEntry{
		MDEntryType:         pmap.NextString(iter.NextString, prev.MDEntryType),
		MDEntryID:           iter.NextString(),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		OrderStatus:         pmap.NextString(iter.NextString, prev.OrderStatus),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
	}
}

func DecodeWTLSFOND(iter *decoder.MessageIterator) *WTLSFOND {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastMsgSeqNumProcessed := iter.NextUInt32(true)
	rptSeq := iter.NextInt32(false)
	lastFragment := iter.NextUInt32(true)
	routeFirst := iter.NextUInt32(true)
	tradSesStatus := iter.NextInt32(true)
	tradingSessionID := iter.NextString()
	symbol := iter.NextString()
	mdSecurityTradingStatus := iter.NextInt32(true)
	auctionIndicator := iter.NextUInt32(true)
	noMDEntries := iter.NextUInt32(false)
	lastMDEntry := &WTLSFONDGroupMDEntry{}
	groupMDEntries := make([]*WTLSFONDGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeWTLSFONDGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WTLSFOND", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &WTLSFOND{
		MsgSeqNum:               msgSeqNum,
		SendingTime:             sendingTime,
		LastMsgSeqNumProcessed:  lastMsgSeqNumProcessed,
		RptSeq:                  rptSeq,
		LastFragment:            lastFragment,
		RouteFirst:              routeFirst,
		TradSesStatus:           tradSesStatus,
		TradingSessionID:        tradingSessionID,
		Symbol:                  symbol,
		MDSecurityTradingStatus: mdSecurityTradingStatus,
		AuctionIndicator:        auctionIndicator,
		NoMDEntries:             noMDEntries,
		GroupMDEntries:          groupMDEntries,
	}
}

func DecodeWTLSFONDGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *WTLSFONDGroupMDEntry,
) *WTLSFONDGroupMDEntry {
	return &WTLSFONDGroupMDEntry{
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		OrderSide:           pmap.NextString(iter.NextString, prev.OrderSide),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		AccruedInterestAmt:  pmap.NextDecimal(iter.NextDecimalOptional, prev.AccruedInterestAmt),
		TradeValue:          pmap.NextDecimal(iter.NextDecimalOptional, prev.TradeValue),
		Yield:               pmap.NextDecimal(iter.NextDecimalOptional, prev.Yield),
		SettlDate:           pmap.NextUInt32(iter.NextUInt32Optional, prev.SettlDate),
		SettleType:          pmap.NextString(iter.NextString, prev.SettleType),
		Price:               pmap.NextDecimal(iter.NextDecimalOptional, prev.Price),
		PriceType:           pmap.NextInt32(iter.NextInt32Optional, prev.PriceType),
		RepoToPx:            pmap.NextDecimal(iter.NextDecimalOptional, prev.RepoToPx),
		BuyBackPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.BuyBackPx),
		BuyBackDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.BuyBackDate),
		TotalVolume:         iter.NextDecimal(true),
		RepoTerm:            pmap.NextUInt32(iter.NextUInt32Optional, prev.RepoTerm),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
		RefOrderID:          pmap.NextString(iter.NextString, prev.RefOrderID),
	}
}

func DecodeWTLSCURR(iter *decoder.MessageIterator) *WTLSCURR {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastMsgSeqNumProcessed := iter.NextUInt32(true)
	rptSeq := iter.NextInt32(false)
	lastFragment := iter.NextUInt32(true)
	routeFirst := iter.NextUInt32(true)
	tradSesStatus := iter.NextInt32(true)
	tradingSessionID := iter.NextString()
	symbol := iter.NextString()
	mdSecurityTradingStatus := iter.NextInt32(true)
	noMDEntries := iter.NextUInt32(false)
	lastMDEntry := &WTLSCURRGroupMDEntry{}
	groupMDEntries := make([]*WTLSCURRGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeWTLSCURRGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WTLSCURR", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &WTLSCURR{
		MsgSeqNum:               msgSeqNum,
		SendingTime:             sendingTime,
		LastMsgSeqNumProcessed:  lastMsgSeqNumProcessed,
		RptSeq:                  rptSeq,
		LastFragment:            lastFragment,
		RouteFirst:              routeFirst,
		TradSesStatus:           tradSesStatus,
		TradingSessionID:        tradingSessionID,
		Symbol:                  symbol,
		MDSecurityTradingStatus: mdSecurityTradingStatus,
		NoMDEntries:             noMDEntries,
		GroupMDEntries:          groupMDEntries,
	}
}

func DecodeWTLSCURRGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *WTLSCURRGroupMDEntry,
) *WTLSCURRGroupMDEntry {
	return &WTLSCURRGroupMDEntry{
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		OrderSide:           pmap.NextString(iter.NextString, prev.OrderSide),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		TradeValue:          pmap.NextDecimal(iter.NextDecimalOptional, prev.TradeValue),
		SettlDate:           pmap.NextUInt32(iter.NextUInt32Optional, prev.SettlDate),
		SettleType:          pmap.NextString(iter.NextString, prev.SettleType),
		Price:               pmap.NextDecimal(iter.NextDecimalOptional, prev.Price),
		PriceType:           pmap.NextInt32(iter.NextInt32Optional, prev.PriceType),
		RepoToPx:            pmap.NextDecimal(iter.NextDecimalOptional, prev.RepoToPx),
		BuyBackPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.BuyBackPx),
		BuyBackDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.BuyBackDate),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
		RefOrderID:          pmap.NextString(iter.NextString, prev.RefOrderID),
	}
}

func DecodeXMSRFOND(iter *decoder.MessageIterator) *XMSRFOND {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastUpdateTime := iter.NextUInt64(true)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*XMSRFONDGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		groupMDEntries[i] = DecodeXMSRFONDGroupMDEntry(iter)
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in WMSRFOND", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XMSRFOND{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		LastUpdateTime: lastUpdateTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXMSRFONDGroupMDEntry(iter *decoder.MessageIterator) *XMSRFONDGroupMDEntry {
	return &XMSRFONDGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		Symbol:              iter.NextString(),
		RptSeq:              iter.NextInt32(true),
		TotalNumOfTrades:    iter.NextInt32(true),
		TradeValue:          iter.NextDecimal(true),
		OfferNbOr:           iter.NextInt32(true),
		BidNbOr:             iter.NextInt32(true),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		MDEntryDate:         iter.NextUInt32(true),
		MDEntryTime:         iter.NextUInt32(true),
		StartTime:           iter.NextUInt32(true),
		QuoteCondition:      iter.NextString(),
		TradeCondition:      iter.NextString(),
		OpenCloseSettlFlag:  iter.NextString(),
		NetChgPrevDay:       iter.NextDecimal(true),
		PriceImprovement:    iter.NextDecimal(true),
		AccruedInterestAmt:  iter.NextDecimal(true),
		ChgFromWAPrice:      iter.NextDecimal(true),
		ChgOpenInterest:     iter.NextDecimal(true),
		BidMarketSize:       iter.NextDecimal(true),
		AskMarketSize:       iter.NextDecimal(true),
		Yield:               iter.NextDecimal(true),
		ChgFromSettlmnt:     iter.NextDecimal(true),
		MinCurrPx:           iter.NextDecimal(true),
		MinCurrPxChgTime:    iter.NextUInt32(true),
		IndexCrossRate:      iter.NextDecimal(true),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		VolumeIndicator:     iter.NextUInt32(true),
		SettlDate:           iter.NextUInt32(true),
		SettlDate2:          iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		CXFlag:              iter.NextString(),
		TradingSessionID:    iter.NextString(),
		TradingSessionSubID: iter.NextString(),
	}
}

func DecodeXMSRCURR(iter *decoder.MessageIterator) *XMSRCURR {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	lastUpdateTime := iter.NextUInt64(true)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*XMSRCURRGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		groupMDEntries[i] = DecodeXMSRCURRGroupMDEntry(iter)
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XMSRCURR", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XMSRCURR{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		LastUpdateTime: lastUpdateTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXMSRCURRGroupMDEntry(iter *decoder.MessageIterator) *XMSRCURRGroupMDEntry {
	return &XMSRCURRGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		Symbol:              iter.NextString(),
		RptSeq:              iter.NextInt32(true),
		TotalNumOfTrades:    iter.NextInt32(true),
		TradeValue:          iter.NextDecimal(true),
		OfferNbOr:           iter.NextInt32(true),
		BidNbOr:             iter.NextInt32(true),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		MDEntryDate:         iter.NextUInt32(true),
		MDEntryTime:         iter.NextUInt32(true),
		StartTime:           iter.NextUInt32(true),
		QuoteCondition:      iter.NextString(),
		TradeCondition:      iter.NextString(),
		OpenCloseSettlFlag:  iter.NextString(),
		NetChgPrevDay:       iter.NextDecimal(true),
		ChgFromWAPrice:      iter.NextDecimal(true),
		ChgOpenInterest:     iter.NextDecimal(true),
		BidMarketSize:       iter.NextDecimal(true),
		AskMarketSize:       iter.NextDecimal(true),
		ChgFromSettlmnt:     iter.NextDecimal(true),
		SettlDate:           iter.NextUInt32(true),
		SettlDate2:          iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		CXFlag:              iter.NextString(),
		TradingSessionID:    iter.NextString(),
		TradingSessionSubID: iter.NextString(),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
	}
}

func DecodeXOLRFOND(iter *decoder.MessageIterator) *XOLRFOND {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*XOLRFONDGroupMDEntry, noMDEntries)
	lastMDEntry := &XOLRFONDGroupMDEntry{}
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeXOLRFONDGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XOLRFOND", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XOLRFOND{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXOLRFONDGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *XOLRFONDGroupMDEntry,
) *XOLRFONDGroupMDEntry {
	return &XOLRFONDGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         pmap.NextString(iter.NextString, prev.MDEntryType),
		MDEntryID:           iter.NextString(),
		Symbol:              pmap.NextString(iter.NextString, prev.Symbol),
		RptSeq:              iter.NextInt32(true),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		Yield:               pmap.NextDecimal(iter.NextDecimalOptional, prev.Yield),
		OrderStatus:         pmap.NextString(iter.NextString, prev.OrderStatus),
		OrdType:             pmap.NextString(iter.NextString, prev.OrdType),
		TotalVolume:         pmap.NextDecimal(iter.NextDecimalOptional, prev.TotalVolume),
		TradingSession:      pmap.NextString(iter.NextString, prev.TradingSession),
		SecurityStatistics:  pmap.NextString(iter.NextString, prev.SecurityStatistics),
		TradingSessionID:    pmap.NextString(iter.NextString, prev.TradingSessionID),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
	}
}

func DecodeXOLRCURR(iter *decoder.MessageIterator) *XOLRCURR {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	noMDEntries := iter.NextUInt32(false)
	lastMDEntry := &XOLRCURRGroupMDEntry{}
	groupMDEntries := make([]*XOLRCURRGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeXOLRCURRGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XOLRCURR", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XOLRCURR{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXOLRCURRGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *XOLRCURRGroupMDEntry,
) *XOLRCURRGroupMDEntry {
	return &XOLRCURRGroupMDEntry{
		MDUpdateAction:      pmap.NextUInt32(iter.NextUInt32Optional, prev.MDUpdateAction),
		MDEntryType:         pmap.NextString(iter.NextString, prev.MDEntryType),
		MDEntryID:           iter.NextString(),
		Symbol:              pmap.NextString(iter.NextString, prev.Symbol),
		RptSeq:              iter.NextInt32(true),
		MDEntryPx:           pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntryPx),
		MDEntrySize:         pmap.NextDecimal(iter.NextDecimalOptional, prev.MDEntrySize),
		MDEntryDate:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryDate),
		MDEntryTime:         pmap.NextUInt32(iter.NextUInt32Optional, prev.MDEntryTime),
		OrigTime:            pmap.NextUInt32(iter.NextUInt32Optional, prev.OrigTime),
		OrderStatus:         pmap.NextString(iter.NextString, prev.OrderStatus),
		TradingSessionID:    pmap.NextString(iter.NextString, prev.TradingSessionID),
		TradingSessionSubID: pmap.NextString(iter.NextString, prev.TradingSessionSubID),
		TradingSession:      pmap.NextString(iter.NextString, prev.TradingSession),
		SecurityStatistics:  pmap.NextString(iter.NextString, prev.SecurityStatistics),
	}
}

func DecodeXTLRFOND(iter *decoder.MessageIterator) *XTLRFOND {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	noMDEntries := iter.NextUInt32(false)
	lastMDEntry := &XTLRFONDGroupMDEntry{}
	groupMDEntries := make([]*XTLRFONDGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		pmap := iter.NextPmap()
		groupMDEntries[i] = DecodeXTLRFONDGroupMDEntry(iter, pmap, lastMDEntry)
		lastMDEntry = groupMDEntries[i]
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XTLRFOND", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XTLRFOND{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXTLRFONDGroupMDEntry(
	iter *decoder.MessageIterator,
	pmap *decoder.Pmap,
	prev *XTLRFONDGroupMDEntry,
) *XTLRFONDGroupMDEntry {
	return &XTLRFONDGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		Symbol:              iter.NextString(),
		RptSeq:              iter.NextInt32(true),
		OrderSide:           iter.NextString(),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		RefOrderID:          iter.NextString(),
		MDEntryDate:         iter.NextUInt32(true),
		MDEntryTime:         iter.NextUInt32(true),
		OrigTime:            iter.NextUInt32(true),
		AccruedInterestAmt:  iter.NextDecimal(true),
		TradeValue:          iter.NextDecimal(true),
		Yield:               iter.NextDecimal(true),
		SettlDate:           iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		Price:               iter.NextDecimal(true),
		PriceType:           iter.NextInt32(true),
		RepoToPx:            iter.NextDecimal(true),
		BuyBackPx:           iter.NextDecimal(true),
		BuyBackDate:         iter.NextUInt32(true),
		RepoTerm:            pmap.NextUInt32(iter.NextUInt32Optional, prev.RepoTerm),
		TotalVolume:         iter.NextDecimal(true),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
		TradingSessionID:    iter.NextString(),
		TradingSessionSubID: iter.NextString(),
	}
}
func DecodeXTLRCURR(iter *decoder.MessageIterator) *XTLRCURR {
	msgSeqNum := iter.NextUInt32(false)
	sendingTime := iter.NextUInt64(false)
	noMDEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*XTLRCURRGroupMDEntry, noMDEntries)
	for i := uint32(0); i < noMDEntries; i++ {
		groupMDEntries[i] = DecodeXTLRCURRGroupMDEntry(iter)
	}
	if iter.HasNext() {
		log.Printf("Lasts bytes: %+v in XTLRCURR", iter.Lasts())
		panic("after scanning group md entries by number left more bytes")
	}
	return &XTLRCURR{
		MsgSeqNum:      msgSeqNum,
		SendingTime:    sendingTime,
		NoMDEntries:    noMDEntries,
		GroupMDEntries: groupMDEntries,
	}
}

func DecodeXTLRCURRGroupMDEntry(iter *decoder.MessageIterator) *XTLRCURRGroupMDEntry {
	return &XTLRCURRGroupMDEntry{
		MDUpdateAction:      iter.NextUInt32(true),
		MDEntryType:         iter.NextString(),
		MDEntryID:           iter.NextString(),
		Symbol:              iter.NextString(),
		RptSeq:              iter.NextInt32(true),
		OrderSide:           iter.NextString(),
		MDEntryPx:           iter.NextDecimal(true),
		MDEntrySize:         iter.NextDecimal(true),
		RefOrderID:          iter.NextString(),
		MDEntryDate:         iter.NextUInt32(true),
		MDEntryTime:         iter.NextUInt32(true),
		OrigTime:            iter.NextUInt32(true),
		TradeValue:          iter.NextDecimal(true),
		SettlDate:           iter.NextUInt32(true),
		SettleType:          iter.NextString(),
		Price:               iter.NextDecimal(true),
		PriceType:           iter.NextInt32(true),
		RepoToPx:            iter.NextDecimal(true),
		BuyBackPx:           iter.NextDecimal(true),
		BuyBackDate:         iter.NextUInt32(true),
		TradingSessionID:    iter.NextString(),
		TradingSessionSubID: iter.NextString(),
		TradingSession:      iter.NextString(),
		SecurityStatistics:  iter.NextString(),
	}
}

func DecodeSecurityDefinition(iter *decoder.MessageIterator) *SecurityDefinition {
	msgSeqNum := iter.NextUInt32(true) // increment operator, need to check pmap
	sendingTime := iter.NextUInt64(false)
	totNumReports := iter.NextInt32(true)
	symbol := iter.NextString()
	securityID := iter.NextString()
	securityIDSource := iter.NextString()
	product := iter.NextInt32(true)
	eveningSession := iter.NextString()
	cfiCode := iter.NextByteVectorAsString(true)
	securityType := iter.NextByteVectorAsString(true)
	maturityDate := iter.NextUInt32(true)
	settlDate := iter.NextUInt32(true)
	settleType := iter.NextString()
	origIssueAmt := iter.NextDecimal(true)
	couponPaymentDate := iter.NextUInt32(true)
	couponRate := iter.NextDecimal(true)
	settlFixingDate := iter.NextUInt32(true)
	dividendNetPx := iter.NextDecimal(true)
	securityDesc := iter.NextByteVectorAsString(true)
	encodedSecurityDesc := iter.NextByteVectorAsString(true)
	quoteText := iter.NextByteVectorAsString(true)

	noInstrAttrib := iter.NextUInt32(true)
	groupInstrAttribs := make([]*GroupInstrAttrib, noInstrAttrib)
	for i := uint32(0); i < noInstrAttrib; i++ {
		groupInstrAttribs[i] = DecodeGroupInstrAttrib(iter)
	}

	currency := iter.NextString()

	noMarketSegments := iter.NextUInt32(true)
	marketSegmentGrps := make([]*MarketSegmentGrp, noMarketSegments)
	for i := uint32(0); i < noMarketSegments; i++ {
		marketSegmentGrps[i] = DecodeMarketSegmentGrp(iter)
	}

	settlCurrency := iter.NextString()
	priceType := iter.NextInt32(true)
	stateSecurityID := iter.NextString()
	encodedShortSecurityDesc := iter.NextByteVectorAsString(true)
	marketCode := iter.NextByteVectorAsString(true)
	minPriceIncrement := iter.NextDecimal(true)
	mktShareLimit := iter.NextDecimal(true)
	mktShareThreshold := iter.NextDecimal(true)
	maxOrdersVolume := iter.NextDecimal(true)
	priceMvmLimit := iter.NextDecimal(true)
	faceValue := iter.NextDecimal(true)
	baseSwapPx := iter.NextDecimal(true)
	repoToPx := iter.NextDecimal(true)
	buyBackPx := iter.NextDecimal(true)
	buyBackDate := iter.NextUInt32(true)
	noSharesIssued := iter.NextDecimal(true)
	highLimit := iter.NextDecimal(true)
	lowLimit := iter.NextDecimal(true)
	numOfDaysToMaturity := iter.NextInt32(true)
	return &SecurityDefinition{
		MsgSeqNum:                msgSeqNum,
		SendingTime:              sendingTime,
		TotNumReports:            totNumReports,
		Symbol:                   symbol,
		SecurityID:               securityID,
		SecurityIDSource:         securityIDSource,
		Product:                  product,
		EveningSession:           eveningSession,
		CFICode:                  cfiCode,
		SecurityType:             securityType,
		MaturityDate:             maturityDate,
		SettlDate:                settlDate,
		SettleType:               settleType,
		OrigIssueAmt:             origIssueAmt,
		CouponPaymentDate:        couponPaymentDate,
		CouponRate:               couponRate,
		SettlFixingDate:          settlFixingDate,
		DividendNetPx:            dividendNetPx,
		SecurityDesc:             securityDesc,
		EncodedSecurityDesc:      encodedSecurityDesc,
		QuoteText:                quoteText,
		NoInstrAttrib:            noInstrAttrib,
		GroupInstrAttribs:        groupInstrAttribs,
		Currency:                 currency,
		NoMarketSegments:         noMarketSegments,
		MarketSegmentGrps:        marketSegmentGrps,
		SettlCurrency:            settlCurrency,
		PriceType:                priceType,
		StateSecurityID:          stateSecurityID,
		EncodedShortSecurityDesc: encodedShortSecurityDesc,
		MarketCode:               marketCode,
		MinPriceIncrement:        minPriceIncrement,
		MktShareLimit:            mktShareLimit,
		MktShareThreshold:        mktShareThreshold,
		MaxOrdersVolume:          maxOrdersVolume,
		PriceMvmLimit:            priceMvmLimit,
		FaceValue:                faceValue,
		BaseSwapPx:               baseSwapPx,
		RepoToPx:                 repoToPx,
		BuyBackPx:                buyBackPx,
		BuyBackDate:              buyBackDate,
		NoSharesIssued:           noSharesIssued,
		HighLimit:                highLimit,
		LowLimit:                 lowLimit,
		NumOfDaysToMaturity:      numOfDaysToMaturity,
	}
}

func DecodeGroupInstrAttrib(iter *decoder.MessageIterator) *GroupInstrAttrib {
	return &GroupInstrAttrib{
		InstrAttribType:  iter.NextInt32(false),
		InstrAttribValue: iter.NextByteVector(true),
	}
}

func DecodeMarketSegmentGrp(iter *decoder.MessageIterator) *MarketSegmentGrp {
	roundLot := iter.NextDecimal(true)
	lotDivider := iter.NextUInt32(true)
	noTradingSessionRules := iter.NextUInt32(true)
	tradingSessionRulesGrps := make([]*TradingSessionRulesGrp, noTradingSessionRules)
	for i := uint32(0); i < noTradingSessionRules; i++ {
		tradingSessionRulesGrps[i] = DecodeTradingSessionRulesGrp(iter)
	}
	return &MarketSegmentGrp{
		RoundLot:                roundLot,
		LotDivider:              lotDivider,
		NoTradingSessionRules:   noTradingSessionRules,
		TradingSessionRulesGrps: tradingSessionRulesGrps,
	}
}

func DecodeTradingSessionRulesGrp(iter *decoder.MessageIterator) *TradingSessionRulesGrp {
	return &TradingSessionRulesGrp{
		TradingSessionID:      iter.NextString(),
		TradingSessionSubID:   iter.NextString(),
		SecurityTradingStatus: iter.NextInt32(true),
		OrderNote:             iter.NextInt32(true),
	}
}

func DecodeSecurityStatus(iter *decoder.MessageIterator) *SecurityStatus {
	return &SecurityStatus{
		MsgSeqNum:             iter.NextUInt32(false),
		SendingTime:           iter.NextUInt64(false),
		Symbol:                iter.NextString(),
		TradingSessionID:      iter.NextString(),
		TradingSessionSubID:   iter.NextString(),
		SecurityTradingStatus: iter.NextInt32(true),
		AuctionIndicator:      iter.NextUInt32(true),
	}
}

func DecodeTradingSessionStatus(iter *decoder.MessageIterator) *TradingSessionStatus {
	return &TradingSessionStatus{
		MsgSeqNum:        iter.NextUInt32(false),
		SendingTime:      iter.NextUInt64(false),
		TradSesStatus:    iter.NextInt32(false),
		Text:             iter.NextString(),
		TradingSessionID: iter.NextString(),
	}
}

func DecodeHeartbeat(iter *decoder.MessageIterator) *Heartbeat {
	return &Heartbeat{
		MsgSeqNum:   iter.NextUInt32(false),
		SendingTime: iter.NextUInt64(false),
	}
}
