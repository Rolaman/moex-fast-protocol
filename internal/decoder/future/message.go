package future

import "github.com/shopspring/decimal"

type Message struct {
	TemplateId           uint32
	IncRef               *IncRefMessage
	Snap                 *SnapMessage
	SecurityDefinition   *SecDefMessage
	Reset                *SeqResetMessage
	HeartbeatMessage     *HeartbeatMessage
	TradingSessionStatus *TradingSesStatusMessage
}

type IncRefMessage struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	LastFragment uint32
	NoMDEntries  uint32
	Entries      []*IncRefEntry
}

type IncRefEntry struct {
	MDUpdateAction           uint32
	MDEntryType              string
	SecurityID               uint64
	SecurityIDSource         uint32
	Symbol                   string
	SecurityGroup            string
	ExchangeTradingSessionID uint32
	RptSeq                   uint32
	MarketDepth              uint32
	MDPriceLevel             uint32
	MDEntryID                int64
	MDEntryPx                decimal.Decimal
	MDEntrySize              int64
	MDEntryDate              uint32
	MDEntryTime              uint64
	NumberOfOrders           int32
	MDEntryTradeType         string
	TrdType                  int32
	LastPx                   decimal.Decimal
	MDFlags                  int32
	Currency                 string
	Revision                 uint64
	OrderSide                string
	MDEntrySyntheticSize     int64
}

type SnapMessage struct {
	ApplVerID              string
	MessageType            string
	SenderCompID           string
	MsgSeqNum              uint32
	SendingTime            uint64
	LastFragment           uint32
	RptSeq                 uint32
	TotNumReports          uint32
	LastMsgSeqNumProcessed uint32
	SecurityID             uint64
	SecurityIDSource       uint32
	Symbol                 string
	SecurityGroup          string
	NoMDEntries            uint32
	Entries                []*SnapEntry
}

type SnapEntry struct {
	MDEntryType              string
	ExchangeTradingSessionID uint32
	MDEntryID                int64
	MarketDepth              uint32
	MDEntryPx                decimal.Decimal
	MDEntryDate              uint32
	MDEntryTime              uint64
	MDEntrySize              int64
	MDPriceLevel             uint32
	NumberOfOrders           int32
	MDEntryTradeType         string
	TrdType                  int32
	MDFlags                  int32
	Currency                 string
	OrderSide                string
	MDEntrySyntheticSize     int64
}

type SecDefMessage struct {
	ApplVerID                   string
	MessageType                 string
	SenderCompID                string
	MsgSeqNum                   uint32
	SendingTime                 uint64
	TotNumReports               uint32
	Symbol                      string
	SecurityDesc                string
	SecurityID                  uint64
	SecurityIDSource            uint32
	SecurityAltID               string
	SecurityAltIDSource         string
	SecurityType                string
	CFICode                     string
	StrikePrice                 decimal.Decimal
	ContractMultiplier          decimal.Decimal
	SecurityTradingStatus       uint32
	Currency                    string
	MarketID                    string
	MarketSegmentID             string
	TradingSessionID            uint32
	ExchangeTradingSessionID    uint32
	Volatility                  decimal.Decimal
	NoMDFeedTypes               uint32
	Feeds                       []*SecDefFeed
	NoUnderlyings               uint32
	Underlyings                 []*SecDefUnderlyings
	HighLimitPx                 decimal.Decimal
	LowLimitPx                  decimal.Decimal
	MinPriceIncrement           decimal.Decimal
	MinPriceIncrementAmount     decimal.Decimal
	InitialMarginOnBuy          decimal.Decimal
	InitialMarginOnSell         decimal.Decimal
	InitialMarginSyntetic       decimal.Decimal
	QuotationList               string
	TheorPrice                  decimal.Decimal
	TheorPriceLimit             decimal.Decimal
	NoLegs                      uint32
	InstrumentLegs              []*SecDefInstLegs
	NoInstrAttrib               uint32
	Attrs                       []*SecDefInstAttrs
	UnderlyingQty               decimal.Decimal
	UnderlyingCurrency          string
	NoEvents                    uint32
	EventGroups                 []SecDefEvntGrp
	MaturityDate                uint32
	MaturityTime                uint32
	Flags                       int64
	MinPriceIncrementAmountCurr decimal.Decimal
	SettlPriceOpen              decimal.Decimal
}

type SecDefFeed struct {
	MDFeedType  string
	MarketDepth uint32
	MDBookType  uint32
}

type SecDefUnderlyings struct {
	UnderlyingSymbol     string
	UnderlyingSecurityID uint64
	UnderlyingFutureID   uint64
}

type SecDefInstLegs struct {
	LegSymbol     string
	LegSecurityID uint64
	LegRatioQTY   decimal.Decimal
}

type SecDefInstAttrs struct {
	InstrAttribType  int32
	InstrAttribValue string
}

type SecDefEvntGrp struct {
	EventType int32
	EventDate uint32
	EventTime uint64
}

type SecDefUpdMessage struct {
	ApplVerID        string
	MessageType      string
	SenderCompID     string
	MsgSeqNum        uint32
	SendingTime      uint64
	SecurityID       uint64
	SecurityIDSource uint32
	Volatility       decimal.Decimal
	TheorPrice       decimal.Decimal
	TheorPriceLimit  decimal.Decimal
}

type SecStatusMessage struct {
	ApplVerID             string
	MessageType           string
	SenderCompID          string
	MsgSeqNum             uint32
	SendingTime           uint64
	SecurityID            uint64
	SecurityIDSource      uint32
	Symbol                string
	SecurityTradingStatus uint32
	HighLimitPx           decimal.Decimal
	LowLimitPx            decimal.Decimal
	InitialMarginOnBuy    decimal.Decimal
	InitialMarginOnSell   decimal.Decimal
	InitialMarginSyntetic decimal.Decimal
}

type HeartbeatMessage struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
}

type SeqResetMessage struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	NewSeqNo     uint32
}

type TradingSesStatusMessage struct {
	ApplVerID                      string
	MessageType                    string
	SenderCompID                   string
	MsgSeqNum                      uint32
	SendingTime                    uint64
	TradSesOpenTime                uint64
	TradSesCloseTime               uint64
	TradSesIntermClearingStartTime uint64
	TradSesIntermClearingEndTime   uint64
	TradingSessionID               uint32
	ExchangeTradingSessionID       uint32
	TradSesStatus                  uint32
	MarketID                       string
	MarketSegmentID                string
	TradSesEvent                   int32
}

type NewsMessage struct {
	ApplVerID       string
	MessageType     string
	SenderCompID    string
	MsgSeqNum       uint32
	SendingTime     uint64
	LastFragment    uint32
	NewsId          string
	OrigTime        uint64
	LanguageCode    string
	Urgency         uint32
	Headline        string
	MarketID        string
	MarketSegmentID string
	NoLinesOfText   uint32
}

type NewsText struct {
	Text string
}

type OrdersLogMessage struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	LastFragment uint32
	NoMDEntries  uint32
	MDEntries    []*OrdersLogEntry
}

type OrdersLogEntry struct {
	MDUpdateAction           uint32
	MDEntryType              string
	MDEntryID                int64
	SecurityID               uint64
	SecurityIDSource         uint32
	RptSeq                   uint32
	MDEntryDate              uint32
	MDEntryTime              uint64
	MDEntryPx                decimal.Decimal
	MDEntrySize              int64
	LastPx                   decimal.Decimal
	LastQty                  int64
	TradeID                  int64
	ExchangeTradingSessionID uint32
	MDFlags                  int64
	Revision                 uint64
}

type QuotesLogMessage struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	LastFragment uint32
	NoMDEntries  uint32
}

type QuotesLogEntry struct {
	MDUpdateAction           uint32
	MDEntryType              string
	MDEntryID                int64
	SecurityID               uint64
	SecurityIDSource         uint32
	RptSeq                   uint32
	MDEntryDate              uint32
	MDEntryTime              uint64
	MDEntryPx                decimal.Decimal
	MDEntrySize              int64
	LastPx                   decimal.Decimal
	LastQty                  int64
	TradeID                  int64
	ExchangeTradingSessionID uint32
	MDFlags                  int64
	Revision                 uint64
	OrderID                  uint64
	TrdMatchID               uint64
}

type BookMessage struct {
	ApplVerID                string
	MessageType              string
	SenderCompID             string
	MsgSeqNum                uint32
	SendingTime              uint64
	LastMsgSeqNumProcessed   uint32
	RptSeq                   uint32
	LastFragment             uint32
	RouteFirst               uint32
	ExchangeTradingSessionID uint32
	SecurityID               uint64
	SecurityIDSource         uint32
	NoMDEntries              uint32
}

type BookEntry struct {
	MDEntryType string
	MDEntryID   int64
	MDEntryDate uint32
	MDEntryTime uint64
	MDEntryPx   decimal.Decimal
	MDEntrySize int64
	TradeID     int64
	MDFlags     int64
}

type Logon struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
}

type Logout struct {
	ApplVerID    string
	MessageType  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	Text         string
}
