package config

// Options

var T0Top50Options = InstrumentOption{
	Depth: 50,

	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.12.7",
		SourceIP: "91.203.253.242",
		Port:     48007,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.140.7",
		SourceIP: "91.203.255.242",
		Port:     49007,
	},
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.12.8",
		SourceIP: "91.203.253.242",
		Port:     48008,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.140.8",
		SourceIP: "91.203.255.242",
		Port:     49008,
	},
}

var T1Top50Options = InstrumentOption{
	Depth: 50,

	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.9.7",
		SourceIP: "91.203.253.235",
		Port:     42007,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.137.7",
		SourceIP: "91.203.255.235",
		Port:     43007,
	},
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.9.8",
		SourceIP: "91.203.253.235",
		Port:     42008,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.137.8",
		SourceIP: "91.203.255.235",
		Port:     43008,
	},
}

var T0Top5Options = InstrumentOption{
	Depth: 5,

	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.12.3",
		SourceIP: "91.203.253.242",
		Port:     48003,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.140.3",
		SourceIP: "91.203.255.242",
		Port:     49003,
	},
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.12.4",
		SourceIP: "91.203.253.242",
		Port:     48004,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.140.4",
		SourceIP: "91.203.255.242",
		Port:     49004,
	},
}

var T1Top5Options = InstrumentOption{
	Depth: 5,

	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.9.3",
		SourceIP: "91.203.253.235",
		Port:     42003,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.137.3",
		SourceIP: "91.203.255.235",
		Port:     43003,
	},
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.9.4",
		SourceIP: "91.203.253.235",
		Port:     42004,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.137.4",
		SourceIP: "91.203.255.235",
		Port:     43004,
	},
}

var T0InfoOptions = FutureInfoOptions{
	SnapshotClient: &ClientOptions{
		GroupIP:  "239.195.12.11",
		SourceIP: "91.203.253.242",
		Port:     48011,
	},
}

var T1InfoOptions = FutureInfoOptions{
	SnapshotClient: &ClientOptions{
		GroupIP:  "239.195.9.11",
		SourceIP: "91.203.253.235",
		Port:     42011,
	},
}

// Currency

var AstsNextCurrencyOrderOptions = InstrumentOption{
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.72",
		SourceIP: "91.203.253.238",
		Port:     16072,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.200",
		SourceIP: "91.203.255.238",
		Port:     17072,
	},
	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.71",
		SourceIP: "91.203.253.238",
		Port:     16071,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.199",
		SourceIP: "91.203.255.238",
		Port:     17071,
	},
}

var AstsUatCurrencyOrderOptions = InstrumentOption{
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.82",
		SourceIP: "91.203.253.238",
		Port:     16082,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.210",
		SourceIP: "91.203.253.238",
		Port:     17082,
	},
	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.81",
		SourceIP: "91.203.253.238",
		Port:     16081,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.209",
		SourceIP: "91.203.253.238",
		Port:     17081,
	},
}

// Stock

var AstsNextStockOrderOptions = InstrumentOption{
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.102",
		SourceIP: "91.203.253.239",
		Port:     16102,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.230",
		SourceIP: "91.203.253.239",
		Port:     17102,
	},
	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.101",
		SourceIP: "91.203.253.239",
		Port:     16101,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.229",
		SourceIP: "91.203.253.239",
		Port:     17101,
	},
}

var AstsUatStockOrderOptions = InstrumentOption{
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.114",
		SourceIP: "91.203.253.239",
		Port:     16114,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.242",
		SourceIP: "91.203.253.239",
		Port:     17114,
	},
	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.101",
		SourceIP: "91.203.253.239",
		Port:     16101,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.229",
		SourceIP: "91.203.253.239",
		Port:     17101,
	},
}

var AstsNextStockTradeOptions = InstrumentOption{
	SnapshotClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.106",
		SourceIP: "239.195.1.234",
		Port:     16106,
	},
	SnapshotClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.200",
		SourceIP: "91.203.253.239",
		Port:     17106,
	},
	IncrementalClientAOptions: &ClientOptions{
		GroupIP:  "239.195.1.105",
		SourceIP: "91.203.253.239",
		Port:     16105,
	},
	IncrementalClientBOptions: &ClientOptions{
		GroupIP:  "239.195.1.233",
		SourceIP: "91.203.253.239",
		Port:     17105,
	},
}

var AstsUatCurrencyInstrumentsAOptions = &ClientOptions{
	GroupIP:  "239.195.1.89",
	SourceIP: "91.203.253.238",
	Port:     16089,
}

var AstsUatStockInstrumentsAOptions = &ClientOptions{
	GroupIP:  "239.195.1.121",
	SourceIP: "91.203.253.239",
	Port:     16121,
}
