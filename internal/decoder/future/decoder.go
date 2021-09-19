package future

import (
	"log"

	"github.com/kdt-wolf/moex-fast/internal/decoder"
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
			log.Printf("Can't decode bytes %+v", bytes)
		}
	}()

	switch templateId {
	case 19:
		message := DecodeDefaultIncrementalRefreshMessage(iterator)
		return &Message{
			IncRef:     message,
			TemplateId: templateId,
		}, true
	case 20:
		message := DecodeSnapMessage(iterator)
		return &Message{
			Snap:       message,
			TemplateId: templateId,
		}, true
	case 21:
		message := DecodeSecurityDefinition(iterator)
		return &Message{
			SecurityDefinition: message,
			TemplateId:         templateId,
		}, true
	case 7:
		message := DecodeSeqResetMessage(iterator)
		return &Message{
			Reset:      message,
			TemplateId: templateId,
		}, true
	case 6:
		message := DecodeHeartbeatMessage(iterator)
		return &Message{
			HeartbeatMessage: message,
			TemplateId:       templateId,
		}, true
	case 8:
		message := DecodeTradingSessionStatus(iterator)
		return &Message{
			TradingSessionStatus: message,
			TemplateId:           templateId,
		}, true
	default:
		log.Printf("Return empty message, templateId: %d", templateId)
		return &Message{}, true
	}
}

func DecodeHeartbeatMessage(iter *decoder.MessageIterator) *HeartbeatMessage {
	return &HeartbeatMessage{
		MsgSeqNum:   iter.NextUInt32(false),
		SendingTime: iter.NextUInt64(false),
	}
}

func DecodeSeqResetMessage(iter *decoder.MessageIterator) *SeqResetMessage {
	return &SeqResetMessage{
		MsgSeqNum:   iter.NextUInt32(false),
		SendingTime: iter.NextUInt64(false),
		NewSeqNo:    iter.NextUInt32(false),
	}
}

func DecodeDefaultIncrementalRefreshMessage(iter *decoder.MessageIterator) *IncRefMessage {
	seqNum := iter.NextUInt32(false)
	sendTime := iter.NextUInt64(false)
	lastFragment := iter.NextUInt32(true)
	nEntries := iter.NextUInt32(false)
	groupMDEntries := make([]*IncRefEntry, nEntries)
	for i := uint32(0); i < nEntries; i++ {
		groupMDEntries[i] = DecodeIncRefEntry(iter)
	}
	return &IncRefMessage{
		MsgSeqNum:    seqNum,
		SendingTime:  sendTime,
		LastFragment: lastFragment,
		NoMDEntries:  nEntries,
		Entries:      groupMDEntries,
	}
}

func DecodeSnapMessage(iter *decoder.MessageIterator) *SnapMessage {
	seqNum := iter.NextUInt32(false)
	sendTime := iter.NextUInt64(false)
	lastFragment := iter.NextUInt32(true)
	rptSeq := iter.NextUInt32(false)
	totNumRep := iter.NextUInt32(false)
	lastMsgSeqNumProcessed := iter.NextUInt32(false)
	secId := iter.NextUInt64(true)
	symbol := iter.NextString()
	secGp := iter.NextString()
	nEntries := iter.NextUInt32(false)
	entries := make([]*SnapEntry, nEntries)
	for i := uint32(0); i < nEntries; i++ {
		entries[i] = DecodeSnapEntry(iter)
	}
	return &SnapMessage{
		MsgSeqNum:              seqNum,
		SendingTime:            sendTime,
		LastFragment:           lastFragment,
		RptSeq:                 rptSeq,
		TotNumReports:          totNumRep,
		LastMsgSeqNumProcessed: lastMsgSeqNumProcessed,
		SecurityID:             secId,
		Symbol:                 symbol,
		SecurityGroup:          secGp,
		NoMDEntries:            nEntries,
		Entries:                entries,
	}
}

func DecodeSnapEntry(iter *decoder.MessageIterator) *SnapEntry {
	return &SnapEntry{
		MDEntryType:              iter.NextString(),
		ExchangeTradingSessionID: iter.NextUInt32(true),
		MDEntryID:                iter.NextInt64(true),
		MarketDepth:              iter.NextUInt32(true),
		MDEntryPx:                iter.NextDecimal(true),
		MDEntryDate:              iter.NextUInt32(true),
		MDEntryTime:              iter.NextUInt64(false),
		MDEntrySize:              iter.NextInt64(true),
		MDPriceLevel:             iter.NextUInt32(true),
		NumberOfOrders:           iter.NextInt32(true),
		MDEntryTradeType:         iter.NextString(),
		TrdType:                  iter.NextInt32(true),
		MDFlags:                  iter.NextInt32(true),
		Currency:                 iter.NextString(),
		OrderSide:                iter.NextString(),
		MDEntrySyntheticSize:     iter.NextInt64(true),
	}
}

func DecodeIncRefEntry(iter *decoder.MessageIterator) *IncRefEntry {
	return &IncRefEntry{
		MDUpdateAction:           iter.NextUInt32(false),
		MDEntryType:              iter.NextString(),
		SecurityID:               iter.NextUInt64(true),
		Symbol:                   iter.NextString(),
		SecurityGroup:            iter.NextString(),
		ExchangeTradingSessionID: iter.NextUInt32(true),
		RptSeq:                   iter.NextUInt32(false),
		MarketDepth:              iter.NextUInt32(true),
		MDPriceLevel:             iter.NextUInt32(true),
		MDEntryID:                iter.NextInt64(true),
		MDEntryPx:                iter.NextDecimal(true),
		MDEntrySize:              iter.NextInt64(true),
		MDEntryDate:              iter.NextUInt32(true),
		MDEntryTime:              iter.NextUInt64(false),
		NumberOfOrders:           iter.NextInt32(true),
		MDEntryTradeType:         iter.NextString(),
		TrdType:                  iter.NextInt32(true),
		LastPx:                   iter.NextDecimal(true),
		MDFlags:                  iter.NextInt32(true),
		Currency:                 iter.NextString(),
		OrderSide:                iter.NextString(),
		Revision:                 iter.NextUInt64(true),
		MDEntrySyntheticSize:     iter.NextInt64(true),
	}
}

func DecodeSecurityDefinition(iter *decoder.MessageIterator) *SecDefMessage {
	msgSeqNum := iter.NextUInt32(false)
	sentTime := iter.NextUInt64(false)
	totNumRep := iter.NextUInt32(false)
	symbol := iter.NextString()
	secDesc := iter.NextUnicode()
	secId := iter.NextUInt64(false)
	secAltId := iter.NextString()
	secAltIDSource := iter.NextString()
	secType := iter.NextString()
	cfiCode := iter.NextString()
	strikePrice := iter.NextDecimal(true)
	contractMult := iter.NextDecimal(true)
	secTradStat := iter.NextUInt32(true)
	currency := iter.NextString()
	marketSegId := iter.NextString()
	tradSessId := iter.NextUInt32(true)
	exchTradSessId := iter.NextUInt32(true)
	volat := iter.NextDecimal(true)
	nFeeds := iter.NextUInt32(false)
	feeds := make([]*SecDefFeed, nFeeds)
	for i := uint32(0); i < nFeeds; i++ {
		feeds[i] = DecodeSecDefFeed(iter)
	}
	nUnderl := iter.NextUInt32(true)
	underlyings := make([]*SecDefUnderlyings, nUnderl)
	for i := uint32(0); i < nUnderl; i++ {
		underlyings[i] = DecodeUnderlying(iter)
	}
	highLimPx := iter.NextDecimal(true)
	lowLimPx := iter.NextDecimal(true)
	minPrInc := iter.NextDecimal(true)
	minPrIncAm := iter.NextDecimal(true)
	initMargBuy := iter.NextDecimal(true)
	initMargSell := iter.NextDecimal(true)
	initMargSynt := iter.NextDecimal(true)
	quotList := iter.NextString()
	theorPr := iter.NextDecimal(true)
	treorPrLim := iter.NextDecimal(true)
	nInstLegs := iter.NextUInt32(true)
	instLegs := make([]*SecDefInstLegs, nInstLegs)
	for i := uint32(0); i < nInstLegs; i++ {
		instLegs[i] = DecodeInstLegs(iter)
	}
	nInstAttr := iter.NextUInt32(true)
	attr := make([]*SecDefInstAttrs, nInstAttr)
	for i := uint32(0); i < nInstAttr; i++ {
		attr[i] = DecodeInstAttr(iter)
	}
	underlyingQty := iter.NextDecimal(true)
	underlyingCurrency := iter.NextString()
	nEvent := iter.NextUInt32(true)
	events := make([]SecDefEvntGrp, nEvent)
	for i := uint32(0); i < nEvent; i++ {
		events[i] = DecodeEvent(iter)
	}
	matDate := iter.NextUInt32(true)
	matTime := iter.NextUInt32(true)
	flags := iter.NextInt64(true)
	minPrIncAmount := iter.NextDecimal(true)
	settPrOpen := iter.NextDecimal(true)
	return &SecDefMessage{
		MsgSeqNum:                   msgSeqNum,
		SendingTime:                 sentTime,
		TotNumReports:               totNumRep,
		Symbol:                      symbol,
		SecurityDesc:                secDesc,
		SecurityID:                  secId,
		SecurityAltID:               secAltId,
		SecurityAltIDSource:         secAltIDSource,
		SecurityType:                secType,
		CFICode:                     cfiCode,
		StrikePrice:                 strikePrice,
		ContractMultiplier:          contractMult,
		SecurityTradingStatus:       secTradStat,
		Currency:                    currency,
		MarketSegmentID:             marketSegId,
		TradingSessionID:            tradSessId,
		ExchangeTradingSessionID:    exchTradSessId,
		Volatility:                  volat,
		NoMDFeedTypes:               nFeeds,
		Feeds:                       feeds,
		NoUnderlyings:               nUnderl,
		Underlyings:                 underlyings,
		HighLimitPx:                 highLimPx,
		LowLimitPx:                  lowLimPx,
		MinPriceIncrement:           minPrInc,
		MinPriceIncrementAmount:     minPrIncAm,
		InitialMarginOnBuy:          initMargBuy,
		InitialMarginOnSell:         initMargSell,
		InitialMarginSyntetic:       initMargSynt,
		QuotationList:               quotList,
		TheorPrice:                  theorPr,
		TheorPriceLimit:             treorPrLim,
		NoLegs:                      nInstLegs,
		InstrumentLegs:              instLegs,
		NoInstrAttrib:               nInstAttr,
		Attrs:                       attr,
		UnderlyingQty:               underlyingQty,
		UnderlyingCurrency:          underlyingCurrency,
		NoEvents:                    nEvent,
		EventGroups:                 events,
		MaturityDate:                matDate,
		MaturityTime:                matTime,
		Flags:                       flags,
		MinPriceIncrementAmountCurr: minPrIncAmount,
		SettlPriceOpen:              settPrOpen,
	}
}

func DecodeInstAttr(iter *decoder.MessageIterator) *SecDefInstAttrs {
	return &SecDefInstAttrs{
		InstrAttribType:  iter.NextInt32(false),
		InstrAttribValue: iter.NextString(),
	}
}

func DecodeInstLegs(iter *decoder.MessageIterator) *SecDefInstLegs {
	return &SecDefInstLegs{
		LegSymbol:     iter.NextString(),
		LegSecurityID: iter.NextUInt64(false),
		LegRatioQTY:   iter.NextDecimal(false),
	}
}

func DecodeSecDefFeed(iter *decoder.MessageIterator) *SecDefFeed {
	return &SecDefFeed{
		MDFeedType:  iter.NextString(),
		MarketDepth: iter.NextUInt32(true),
		MDBookType:  iter.NextUInt32(true),
	}
}

func DecodeUnderlying(iter *decoder.MessageIterator) *SecDefUnderlyings {
	return &SecDefUnderlyings{
		UnderlyingSymbol:     iter.NextString(),
		UnderlyingSecurityID: iter.NextUInt64(true),
		UnderlyingFutureID:   iter.NextUInt64(true),
	}
}

func DecodeEvent(iter *decoder.MessageIterator) SecDefEvntGrp {
	return SecDefEvntGrp{
		EventType: iter.NextInt32(false),
		EventDate: iter.NextUInt32(false),
		EventTime: iter.NextUInt64(false),
	}
}

func DecodeTradingSessionStatus(iter *decoder.MessageIterator) *TradingSesStatusMessage {
	return &TradingSesStatusMessage{
		MsgSeqNum:                      iter.NextUInt32(false),
		SendingTime:                    iter.NextUInt64(false),
		TradSesOpenTime:                iter.NextUInt64(false),
		TradSesCloseTime:               iter.NextUInt64(false),
		TradSesIntermClearingStartTime: iter.NextUInt64(true),
		TradSesIntermClearingEndTime:   iter.NextUInt64(true),
		TradingSessionID:               iter.NextUInt32(false),
		ExchangeTradingSessionID:       iter.NextUInt32(true),
		TradSesStatus:                  iter.NextUInt32(false),
		MarketSegmentID:                iter.NextString(),
		TradSesEvent:                   iter.NextInt32(true),
	}
}
