package decoder

import (
	"github.com/shopspring/decimal"
	"math"
)

type Pmap struct {
	bytes  []byte
	cursor int
}

func newPmap(b []byte) *Pmap {
	return &Pmap{
		bytes:  b,
		cursor: 0,
	}
}

func (p *Pmap) HasNext() bool {
	// Remove first bit from every bite in PMap
	byteNumber := p.cursor / 7
	if len(p.bytes) <= byteNumber {
		return false
	}
	inByte := p.cursor % 7
	pow := math.Pow(2, float64(6-inByte))
	p.cursor++
	return p.bytes[byteNumber]&byte(pow) > 0
}

func (p *Pmap) NextUInt32(f func() uint32, defaultValue uint32) uint32 {
	if p.HasNext() {
		return f()
	}
	return defaultValue
}

func (p *Pmap) NextInt32(f func() int32, defaultValue int32) int32 {
	if p.HasNext() {
		return f()
	}
	return defaultValue
}

func (p *Pmap) NextString(f func() string, defaultValue string) string {
	if p.HasNext() {
		return f()
	}
	return defaultValue
}

func (p *Pmap) NextDecimal(f func() decimal.Decimal, defaultValue decimal.Decimal) decimal.Decimal {
	if p.HasNext() {
		return f()
	}
	return defaultValue
}
