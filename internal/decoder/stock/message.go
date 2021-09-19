package stock

import "github.com/shopspring/decimal"

type Message struct {
	TemplateId           uint32
	Logon                *Logon
	Logout               *Logout
	WGeneric             *WGeneric
	XGeneric             *XGeneric
	WOLSFOND             *WOLSFOND
	WOLSCURR             *WOLSCURR
	WTLSFOND             *WTLSFOND
	WTLSCURR             *WTLSCURR
	XMSRFOND             *XMSRFOND
	XMSRCURR             *XMSRCURR
	XOLRFOND             *XOLRFOND
	XOLRCURR             *XOLRCURR
	XTLRFOND             *XTLRFOND
	XTLRCURR             *XTLRCURR
	SecurityDefinition   *SecurityDefinition
	SecurityStatus       *SecurityStatus
	TradingSessionStatus *TradingSessionStatus
	Heartbeat            *Heartbeat
}

type Logon struct {
	MessageType      string
	BeginString      string
	SenderCompID     string
	TargetCompID     string
	MsgSeqNum        uint32
	SendingTime      uint64
	HeartBtInt       int32
	Username         string
	Password         string
	DefaultApplVerID string
}

type Logout struct {
	MessageType  string
	BeginString  string
	SenderCompID string
	TargetCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
	Text         string
}

// WGeneric Market Data - Snapshot/Full Refresh Generic
type WGeneric struct {
	MessageType             string
	BeginString             string
	ApplVerID               string
	SenderCompID            string
	MsgSeqNum               uint32
	SendingTime             uint64
	TradingSessionID        string
	Symbol                  string
	LastMsgSeqNumProcessed  uint32
	RptSeq                  int32
	LastFragment            uint32
	RouteFirst              uint32
	TradSesStatus           int32
	MDSecurityTradingStatus int32
	AuctionIndicator        uint32
	NetChgPrevDay           decimal.Decimal
	PriceImprovement        decimal.Decimal
	NoMDEntries             uint32
	GroupMDEntries          []*WGenericGroupMDEntry
}

type WGenericGroupMDEntry struct {
	MDEntryType         string
	MDEntryID           string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	QuoteCondition      string
	TradeCondition      string
	OpenCloseSettlFlag  string
	OrdType             string
	EffectiveTime       uint32
	StartTime           uint32
	AccruedInterestAmt  decimal.Decimal
	ChgFromWAPrice      decimal.Decimal
	ChgOpenInterest     decimal.Decimal
	BidMarketSize       decimal.Decimal
	AskMarketSize       decimal.Decimal
	TotalNumOfTrades    int32
	TradeValue          decimal.Decimal
	Yield               decimal.Decimal
	TotalVolume         decimal.Decimal
	OfferNbOr           int32
	BidNbOr             int32
	ChgFromSettlmnt     decimal.Decimal
	SettlDate           uint32
	SettlDate2          uint32
	SettleType          string
	SumQtyOfBest        uint32
	OrderSide           string
	OrderStatus         string
	MinCurrPx           decimal.Decimal
	IndexCrossRate      decimal.Decimal
	MinCurrPxChgTime    uint32
	TradingSession      string
	SecurityStatistics  string
	VolumeIndicator     uint32
	Price               decimal.Decimal
	PriceType           int32
	NominalValue        decimal.Decimal
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	CXFlag              string
	TradingSessionSubID string
}

// XGeneric Market Data - Incremental Refresh Generic
type XGeneric struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	NoMDEntries    uint32
	GroupMDEntries []*XGenericGroupMDEntry
}

type XGenericGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	RptSeq              int32
	MDEntryDate         uint32
	OrigTime            uint32
	SettlDate           uint32
	SettlDate2          uint32
	SettleType          string
	MDEntryTime         uint32
	EffectiveTime       uint32
	StartTime           uint32
	Symbol              string
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	QuoteCondition      string
	TradeCondition      string
	OpenCloseSettlFlag  string
	OrdType             string
	NetChgPrevDay       decimal.Decimal
	PriceImprovement    decimal.Decimal
	AccruedInterestAmt  decimal.Decimal
	ChgFromWAPrice      decimal.Decimal
	ChgOpenInterest     decimal.Decimal
	BidMarketSize       decimal.Decimal
	AskMarketSize       decimal.Decimal
	TotalNumOfTrades    int32
	TradeValue          decimal.Decimal
	Yield               decimal.Decimal
	TotalVolume         decimal.Decimal
	OfferNbOr           int32
	BidNbOr             int32
	ChgFromSettlmnt     decimal.Decimal
	SumQtyOfBest        int32
	OrderSide           string
	OrderStatus         string
	MinCurrPx           decimal.Decimal
	MinCurrPxChgTime    uint32
	IndexCrossRate      decimal.Decimal
	TradingSession      string
	SecurityStatistics  string
	VolumeIndicator     uint32
	Price               decimal.Decimal
	PriceType           int32
	NominalValue        decimal.Decimal
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	CXFlag              string
	RepoTerm            uint32 //
	TradingSessionID    string
	TradingSessionSubID string
}

// WOLSFOND Market Data - Snapshot/Full Refresh OLS FOND
type WOLSFOND struct {
	MessageType             string
	BeginString             string
	ApplVerID               string
	SenderCompID            string
	MsgSeqNum               uint32
	SendingTime             uint64
	LastMsgSeqNumProcessed  uint32
	RptSeq                  int32
	LastFragment            uint32
	RouteFirst              uint32
	TradSesStatus           int32
	TradingSessionID        string
	Symbol                  string
	MDSecurityTradingStatus int32
	AuctionIndicator        uint32
	NoMDEntries             uint32
	GroupMDEntries          []*WOLSFONDGroupMDEntry
}

type WOLSFONDGroupMDEntry struct {
	MDEntryType         string //
	MDEntryID           string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	MDEntryPx           decimal.Decimal //
	MDEntrySize         decimal.Decimal
	Yield               decimal.Decimal
	OrderStatus         string
	OrdType             string
	TotalVolume         decimal.Decimal
	TradingSession      string
	SecurityStatistics  string
	TradingSessionSubID string
}

// WOLSCURR Market Data - Snapshot/Full Refresh OLS CURR
type WOLSCURR struct {
	MessageType             string
	BeginString             string
	ApplVerID               string
	SenderCompID            string
	MsgSeqNum               uint32
	SendingTime             uint64
	LastMsgSeqNumProcessed  uint32
	RptSeq                  int32
	LastFragment            uint32
	RouteFirst              uint32
	TradSesStatus           int32
	TradingSessionID        string
	Symbol                  string
	MDSecurityTradingStatus int32
	NoMDEntries             uint32
	GroupMDEntries          []*WOLSCURRGroupMDEntry
}

type WOLSCURRGroupMDEntry struct {
	MDEntryType         string
	MDEntryID           string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	OrderStatus         string
	TradingSession      string
	SecurityStatistics  string
	TradingSessionSubID string
}

// WTLSFOND Market Data - Snapshot/Full Refresh TLS FOND
type WTLSFOND struct {
	MessageType             string
	BeginString             string
	ApplVerID               string
	SenderCompID            string
	MsgSeqNum               uint32
	SendingTime             uint64
	LastMsgSeqNumProcessed  uint32
	RptSeq                  int32
	LastFragment            uint32
	RouteFirst              uint32
	TradSesStatus           int32
	TradingSessionID        string
	Symbol                  string
	MDSecurityTradingStatus int32
	AuctionIndicator        uint32
	NoMDEntries             uint32
	GroupMDEntries          []*WTLSFONDGroupMDEntry
}

type WTLSFONDGroupMDEntry struct {
	MDEntryType         string
	MDEntryID           string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	OrderSide           string
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	AccruedInterestAmt  decimal.Decimal
	TradeValue          decimal.Decimal
	Yield               decimal.Decimal
	SettlDate           uint32
	SettleType          string
	Price               decimal.Decimal
	PriceType           int32
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	TotalVolume         decimal.Decimal
	RepoTerm            uint32
	TradingSession      string
	SecurityStatistics  string
	TradingSessionSubID string
	RefOrderID          string
}

// WTLSCURR Market Data - Snapshot/Full Refresh TLS CURR
type WTLSCURR struct {
	MessageType             string
	BeginString             string
	ApplVerID               string
	SenderCompID            string
	MsgSeqNum               uint32
	SendingTime             uint64
	LastMsgSeqNumProcessed  uint32
	RptSeq                  int32
	LastFragment            uint32
	RouteFirst              uint32
	TradSesStatus           int32
	TradingSessionID        string
	Symbol                  string
	MDSecurityTradingStatus int32
	NoMDEntries             uint32
	GroupMDEntries          []*WTLSCURRGroupMDEntry
}

type WTLSCURRGroupMDEntry struct {
	MDEntryType         string
	MDEntryID           string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	OrderSide           string
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	TradeValue          decimal.Decimal
	SettlDate           uint32
	SettleType          string
	Price               decimal.Decimal
	PriceType           int32
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	TradingSession      string
	SecurityStatistics  string
	TradingSessionSubID string
	RefOrderID          string
}

// XMSRFOND Market Data - Incremental Refresh MSR FOND
type XMSRFOND struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	LastUpdateTime uint64
	NoMDEntries    uint32
	GroupMDEntries []*XMSRFONDGroupMDEntry
}

type XMSRFONDGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	TotalNumOfTrades    int32
	TradeValue          decimal.Decimal
	OfferNbOr           int32
	BidNbOr             int32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	MDEntryDate         uint32
	MDEntryTime         uint32
	StartTime           uint32
	QuoteCondition      string
	TradeCondition      string
	OpenCloseSettlFlag  string
	NetChgPrevDay       decimal.Decimal
	PriceImprovement    decimal.Decimal
	AccruedInterestAmt  decimal.Decimal
	ChgFromWAPrice      decimal.Decimal
	ChgOpenInterest     decimal.Decimal
	BidMarketSize       decimal.Decimal
	AskMarketSize       decimal.Decimal
	Yield               decimal.Decimal
	ChgFromSettlmnt     decimal.Decimal
	MinCurrPx           decimal.Decimal
	MinCurrPxChgTime    uint32
	IndexCrossRate      decimal.Decimal
	TradingSession      string
	SecurityStatistics  string
	VolumeIndicator     uint32
	SettlDate           uint32
	SettlDate2          uint32
	SettleType          string
	CXFlag              string
	TradingSessionID    string
	TradingSessionSubID string
}

// XMSRCURR Market Data - Incremental Refresh MSR CURR
type XMSRCURR struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	LastUpdateTime uint64
	NoMDEntries    uint32
	GroupMDEntries []*XMSRCURRGroupMDEntry
}

type XMSRCURRGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	TotalNumOfTrades    int32
	TradeValue          decimal.Decimal
	OfferNbOr           int32
	BidNbOr             int32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	MDEntryDate         uint32
	MDEntryTime         uint32
	StartTime           uint32
	QuoteCondition      string
	TradeCondition      string
	OpenCloseSettlFlag  string
	NetChgPrevDay       decimal.Decimal
	ChgFromWAPrice      decimal.Decimal
	ChgOpenInterest     decimal.Decimal
	BidMarketSize       decimal.Decimal
	AskMarketSize       decimal.Decimal
	ChgFromSettlmnt     decimal.Decimal
	SettlDate           uint32
	SettlDate2          uint32
	SettleType          string
	CXFlag              string
	TradingSessionID    string
	TradingSessionSubID string
	TradingSession      string
	SecurityStatistics  string
}

// XOLRFOND Market Data - Incremental Refresh OLR FOND
type XOLRFOND struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	NoMDEntries    uint32
	GroupMDEntries []*XOLRFONDGroupMDEntry
}

type XOLRFONDGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	Yield               decimal.Decimal
	OrderStatus         string
	OrdType             string
	TotalVolume         decimal.Decimal
	TradingSession      string
	SecurityStatistics  string
	TradingSessionID    string
	TradingSessionSubID string
}

// XOLRCURR Market Data - Incremental Refresh OLR CURR
type XOLRCURR struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	NoMDEntries    uint32
	GroupMDEntries []*XOLRCURRGroupMDEntry
}

type XOLRCURRGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	OrderStatus         string
	TradingSessionID    string
	TradingSessionSubID string
	TradingSession      string
	SecurityStatistics  string
}

// XTLRFOND Market Data - Incremental Refresh TLR FOND
type XTLRFOND struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	NoMDEntries    uint32
	GroupMDEntries []*XTLRFONDGroupMDEntry
}

type XTLRFONDGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	OrderSide           string
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	RefOrderID          string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	AccruedInterestAmt  decimal.Decimal
	TradeValue          decimal.Decimal
	Yield               decimal.Decimal
	SettlDate           uint32
	SettleType          string
	Price               decimal.Decimal
	PriceType           int32
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	RepoTerm            uint32
	TotalVolume         decimal.Decimal
	TradingSession      string
	SecurityStatistics  string
	TradingSessionID    string
	TradingSessionSubID string
}

// XTLRCURR Market Data - Incremental Refresh TLR CURR
type XTLRCURR struct {
	MessageType    string
	ApplVerID      string
	BeginString    string
	SenderCompID   string
	MsgSeqNum      uint32
	SendingTime    uint64
	NoMDEntries    uint32
	GroupMDEntries []*XTLRCURRGroupMDEntry
}

type XTLRCURRGroupMDEntry struct {
	MDUpdateAction      uint32
	MDEntryType         string
	MDEntryID           string
	Symbol              string
	RptSeq              int32
	OrderSide           string
	MDEntryPx           decimal.Decimal
	MDEntrySize         decimal.Decimal
	RefOrderID          string
	MDEntryDate         uint32
	MDEntryTime         uint32
	OrigTime            uint32
	TradeValue          decimal.Decimal
	SettlDate           uint32
	SettleType          string
	Price               decimal.Decimal
	PriceType           int32
	RepoToPx            decimal.Decimal
	BuyBackPx           decimal.Decimal
	BuyBackDate         uint32
	TradingSessionID    string
	TradingSessionSubID string
	TradingSession      string
	SecurityStatistics  string
}

type SecurityDefinition struct {
	MessageType              string
	ApplVerID                string
	BeginString              string
	SenderCompID             string
	MsgSeqNum                uint32
	SendingTime              uint64
	MessageEncoding          string
	TotNumReports            int32
	Symbol                   string
	SecurityID               string
	SecurityIDSource         string
	Product                  int32
	EveningSession           string
	CFICode                  string
	SecurityType             string
	MaturityDate             uint32
	SettlDate                uint32
	SettleType               string
	OrigIssueAmt             decimal.Decimal
	CouponPaymentDate        uint32
	CouponRate               decimal.Decimal
	SettlFixingDate          uint32
	DividendNetPx            decimal.Decimal
	SecurityDesc             string
	EncodedSecurityDesc      string
	QuoteText                string
	NoInstrAttrib            uint32
	GroupInstrAttribs        []*GroupInstrAttrib
	Currency                 string
	NoMarketSegments         uint32
	MarketSegmentGrps        []*MarketSegmentGrp
	SettlCurrency            string
	PriceType                int32
	StateSecurityID          string
	EncodedShortSecurityDesc string
	MarketCode               string
	MinPriceIncrement        decimal.Decimal
	MktShareLimit            decimal.Decimal
	MktShareThreshold        decimal.Decimal
	MaxOrdersVolume          decimal.Decimal
	PriceMvmLimit            decimal.Decimal
	FaceValue                decimal.Decimal
	BaseSwapPx               decimal.Decimal
	RepoToPx                 decimal.Decimal
	BuyBackPx                decimal.Decimal
	BuyBackDate              uint32
	NoSharesIssued           decimal.Decimal
	HighLimit                decimal.Decimal
	LowLimit                 decimal.Decimal
	NumOfDaysToMaturity      int32
}

type GroupInstrAttrib struct {
	InstrAttribType  int32
	InstrAttribValue uint32
}

type MarketSegmentGrp struct {
	RoundLot                decimal.Decimal
	LotDivider              uint32
	NoTradingSessionRules   uint32
	TradingSessionRulesGrps []*TradingSessionRulesGrp
}

type TradingSessionRulesGrp struct {
	TradingSessionID      string
	TradingSessionSubID   string
	SecurityTradingStatus int32
	OrderNote             int32
}

type SecurityStatus struct {
	MessageType           string
	ApplVerID             string
	BeginString           string
	SenderCompID          string
	MsgSeqNum             uint32
	SendingTime           uint64
	Symbol                string
	TradingSessionID      string
	TradingSessionSubID   string
	SecurityTradingStatus int32
	AuctionIndicator      uint32
}

type TradingSessionStatus struct {
	MessageType      string
	ApplVerID        string
	BeginString      string
	SenderCompID     string
	MsgSeqNum        uint32
	SendingTime      uint64
	TradSesStatus    int32
	Text             string
	TradingSessionID string
}

type Heartbeat struct {
	MessageType  string
	BeginString  string
	SenderCompID string
	MsgSeqNum    uint32
	SendingTime  uint64
}
