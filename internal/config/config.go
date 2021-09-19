package config

import (
	"github.com/kdt-wolf/moex-fast/internal/protocol"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Provider struct {
	cfg *GlobalConfig
}

type GlobalConfig struct {
	Futures struct {
		Active []uint64 `yaml:"active"`
	} `yaml:"futures"`
	Stocks struct {
		Active []string `yaml:"active"`
	} `yaml:"stocks"`
}

// ClientOptions Deprecated (use protocol.Options)
type ClientOptions struct {
	GroupIP  string
	SourceIP string
	Port     int
}

type InstrumentOption struct {
	Depth uint32

	IncrementalClientAOptions *ClientOptions
	IncrementalClientBOptions *ClientOptions
	SnapshotClientAOptions    *ClientOptions
	SnapshotClientBOptions    *ClientOptions

	RecoveryOptions *protocol.TcpOptions
}

type FutureInfoOptions struct {
	SnapshotClient *ClientOptions
}

func New() (*Provider, error) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Println("Can't parse config", err)
		return nil, err
	}
	var cfg GlobalConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("Can't parse config", err)
		return nil, err
	}
	return &Provider{
		cfg: &cfg,
	}, nil
}

func (p *Provider) ActiveFutures() map[uint64]bool {
	result := map[uint64]bool{}
	for _, future := range p.cfg.Futures.Active {
		result[future] = true
	}
	return result
}

func (p *Provider) ActiveStocks() map[string]bool {
	result := map[string]bool{}
	for _, stock := range p.cfg.Stocks.Active {
		result[stock] = true
	}
	return result
}
