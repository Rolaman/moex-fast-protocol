package test

import (
	"github.com/kdt-wolf/moex-fast/internal/decoder"
	"github.com/shopspring/decimal"
	"testing"
)

func TestBytesToInt32OptionalPositive(t *testing.T) {
	if 942755 != decoder.Int32Optional([]byte{0x39, 0x45, 0xa4}, true) {
		t.Error("Not 942755")
	}
}

func TestBytesToInt32MandatoryPositive(t *testing.T) {
	if 942755 != decoder.Int32Optional([]byte{0x39, 0x45, 0xa3}, false) {
		t.Error("Not 942755")
	}
}

func TestBytesToInt32OptionalNegative(t *testing.T) {
	if -942755 != decoder.Int32Optional([]byte{0x46, 0x3a, 0xdd}, true) {
		t.Error("Not -942755")
	}
}

func TestBytesToInt32MandatoryNegative(t *testing.T) {
	if -7942755 != decoder.Int32Optional([]byte{0x7c, 0x1b, 0x1b, 0x9d}, false) {
		t.Error("Not -7942755")
	}
}

func TestBytesToInt32MandatoryPositiveSignBit(t *testing.T) {
	if 8193 != decoder.Int32Optional([]byte{0x00, 0x40, 0x81}, false) {
		t.Error("Not 8193")
	}
}

func TestBytesToInt32MandatoryNegativeSignBit(t *testing.T) {
	if -8193 != decoder.Int32Optional([]byte{0x7f, 0x3f, 0xff}, false) {
		t.Error("Not -8193")
	}
}

func TestBytesToUInt32MandatoryPositiveSignBit(t *testing.T) {
	if 942755 != decoder.Int32Optional([]byte{0x39, 0x45, 0xa3}, false) {
		t.Error("Not 942755")
	}
}

func TestBytesToUInt32OptionalNegativeSignBit(t *testing.T) {
	if 942755 != decoder.UInt32Optional([]byte{0x39, 0x45, 0xa4}, true) {
		t.Error("Not 942755")
	}
}

func TestBytesToStringOptional(t *testing.T) {
	if "ABC" != decoder.String([]byte{0x41, 0x42, 0xc3}) {
		t.Error("Not ABC")
	}
}

func TestBytesToByteVectorOptional(t *testing.T) {
	iterator := decoder.NewIterator([]byte{0x84, 0x41, 0x42, 0x43})
	vector := iterator.NextByteVector(true)
	if vector != 1073475 {
		t.Error("Not 1073475")
	}
}

func TestBytesToByteVectorMandatory(t *testing.T) {
	iterator := decoder.NewIterator([]byte{0x83, 0x41, 0x42, 0x43})
	vector := iterator.NextByteVector(false)
	if vector != 1073475 {
		t.Error("Not 1073475")
	}
}

func TestBytesToDecimalPositiveMandatory(t *testing.T) {
	expected := decimal.NewFromInt(94275500)
	iterator := decoder.NewIterator([]byte{0x82, 0x39, 0x45, 0xa3})
	dec := iterator.NextDecimal(false)
	if !dec.Equal(expected) {
		t.Error("Not 94275500")
	}
}

func TestBytesToDecimalPositiveMandatoryScaled(t *testing.T) {
	expected := decimal.NewFromInt(94275500)
	iterator := decoder.NewIterator([]byte{0x81, 0x04, 0x3f, 0x34, 0xde})
	dec := iterator.NextDecimal(false)
	if !dec.Equal(expected) {
		t.Error("Not 94275500")
	}
}

func TestBytesToDecimalPositiveOptional(t *testing.T) {
	expected := decimal.NewFromInt(94275500)
	iterator := decoder.NewIterator([]byte{0x83, 0x39, 0x45, 0xa3})
	dec := iterator.NextDecimal(true)
	if !dec.Equal(expected) {
		t.Error("Not 94275500")
	}
}

func TestBytesToDecimalPositiveMandatoryWithDecimal(t *testing.T) {
	expected := decimal.NewFromFloat(9427.55)
	iterator := decoder.NewIterator([]byte{0xfe, 0x39, 0x45, 0xa3})
	dec := iterator.NextDecimal(true)
	if !dec.Equal(expected) {
		t.Error("Not 9427.55")
	}
}

func TestBytesToDecimalNegativeOptional(t *testing.T) {
	expected := decimal.NewFromFloat(-9427.55)
	iterator := decoder.NewIterator([]byte{0xfe, 0x46, 0x3a, 0xdd})
	dec := iterator.NextDecimal(true)
	if !dec.Equal(expected) {
		t.Error("Not -9427.55")
	}
}
