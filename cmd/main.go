package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	// Launching listeners
	//futureInfoLstnr, _ := futureinfo.New(&config.T0InfoOptions)
	//futureLstnr, _ := futurelistener.New(&config.T0Top5Options, configProvider.ActiveFutures())
	//stockLstnr, _ := stocklistener.New(&config.AstsNextStockOrderOptions, configProvider.ActiveStocks())
	//stockInfoLstnr, _ := stockinfo.New(config.AstsUatStockInstrumentsAOptions)
	//currencyLstnr, _ := currencylistener.New(&config.AstsUatCurrencyOrderOptions, map[string]bool{})
	//currencyInfoLstnr, _ := stockinfo.New(config.AstsUatCurrencyInstrumentsAOptions)

}
