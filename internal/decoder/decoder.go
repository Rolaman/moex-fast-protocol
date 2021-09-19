package decoder

import (
	"log"
	"math"

	"github.com/shopspring/decimal"
)

type MessageIterator struct {
	bytes  []byte
	cursor int
}

func NewIterator(b []byte) *MessageIterator {
	return &MessageIterator{bytes: b}
}

func (iter *MessageIterator) Next() []byte {
	for i := iter.cursor; i < len(iter.bytes); i++ {
		if iter.bytes[i]&128 > 0 {
			result := iter.bytes[iter.cursor : i+1]
			iter.cursor = i + 1
			return result
		}
	}
	iter.cursor = len(iter.bytes)
	return []byte{}
}

func (iter *MessageIterator) nextByte() byte {
	b := iter.bytes[iter.cursor]
	iter.cursor++
	return b
}

func (iter *MessageIterator) NextDecimal(optional bool) decimal.Decimal {
	exp := iter.Next()
	if len(exp) == 1 && exp[0] == 128 {
		return decimal.Decimal{}
	}
	return DecimalOptional(exp, iter.Next(), optional)
}

// for copy operator
func (iter *MessageIterator) NextDecimalOptional() decimal.Decimal {
	exp := iter.Next()
	if len(exp) == 1 && exp[0] == 128 {
		return decimal.Decimal{}
	}
	return DecimalOptional(exp, iter.Next(), true)
}

func (iter *MessageIterator) NextUnicode() string {
	unicodeLen := iter.bytes[iter.cursor] & 127
	bytes := iter.bytes[iter.cursor : iter.cursor+int(unicodeLen)]
	iter.cursor += int(unicodeLen)
	return BytesToUnicodeString(bytes)
}

func (iter *MessageIterator) NextString() string {
	return String(iter.Next())
}

func (iter *MessageIterator) NextStringWithLength() string {
	return StringWithLength(iter)
}

func (iter *MessageIterator) NextInt32(optional bool) int32 {
	return Int32Optional(iter.Next(), optional)
}

func (iter *MessageIterator) NextInt32Optional() int32 {
	return Int32Optional(iter.Next(), true)
}

func (iter *MessageIterator) NextUInt32(optional bool) uint32 {
	return UInt32Optional(iter.Next(), optional)
}

func (iter *MessageIterator) NextUInt32Optional() uint32 {
	return UInt32Optional(iter.Next(), true)
}

func (iter *MessageIterator) NextInt64(optional bool) int64 {
	return Int64Optional(iter.Next(), optional)
}

func (iter *MessageIterator) NextUInt64(optional bool) uint64 {
	return UInt64Optional(iter.Next(), optional)
}

func (iter *MessageIterator) NextN(n int) []byte {
	result := iter.bytes[iter.cursor : iter.cursor+n]
	iter.cursor += n
	return result
}

func (iter *MessageIterator) HasNext() bool {
	return iter.cursor < len(iter.bytes) && iter.bytes[iter.cursor] != 0
}

func (iter *MessageIterator) Lasts() []byte {
	return iter.bytes[iter.cursor:]
}

func (iter *MessageIterator) NextPmap() *Pmap {
	return newPmap(iter.Next())
}

func (iter *MessageIterator) NextByteVector(optional bool) uint32 {
	l := iter.nextByte() & 127
	if optional {
		l -= 1
	}
	if l > 0 {
		return UInt32Optional(iter.NextN(int(l)), false)
	}
	return 0
}

func (iter *MessageIterator) NextByteVectorAsString(optional bool) string {
	l := iter.nextByte() & 127
	if l < 1 {
		return ""
	}
	if optional {
		l -= 1
	}
	if l > 0 {
		return String(iter.NextN(int(l)))
	}
	return ""
}

func Number(iter *MessageIterator) uint32 {
	var number uint32
	nbytes := iter.NextN(4)
	for i, b := range nbytes {
		number = number + uint32(b)*uint32(math.Pow(256, float64(i)))
	}
	return number
}

func PMap(iter *MessageIterator) byte {
	bytes := iter.NextN(1)
	return bytes[0] & 127
}

func Int32Optional(bytes []byte, optional bool) int32 {
	bitLen := len(bytes) * 7
	a := -int32(bytes[0]&64) / 64 * int32(math.Pow(2, float64(bitLen-1)))

	var result int32
	for i, b := range bytes {
		var rawResult int32
		if i == 0 {
			rawResult = int32(b & 63)
		} else if i == len(bytes)-1 {
			rawResult = int32(b & 127)
		} else {
			rawResult = int32(b)
		}
		result += rawResult * int32(math.Pow(128, float64(len(bytes)-i-1)))
	}
	// Because of fast spec 1.1
	if optional && a == 0 {
		result -= 1
	}
	return a + result
}

func Int64Optional(bytes []byte, optional bool) int64 {
	bitLen := len(bytes) * 7
	a := -int64(bytes[0]&64) / 64 * int64(math.Pow(2, float64(bitLen-1)))

	var result int64
	for i, b := range bytes {
		var rawResult int64
		if i == 0 {
			rawResult = int64(b & 63)
		} else if i == len(bytes)-1 {
			rawResult = int64(b & 127)
		} else {
			rawResult = int64(b)
		}
		result += rawResult * int64(math.Pow(128, float64(len(bytes)-i-1)))
	}
	// Because of fast spec 1.1
	if optional && a == 0 {
		result -= 1
	}
	return a + result
}

func UInt32Optional(bytes []byte, optional bool) uint32 {
	l := len(bytes)
	var result uint32 = 0
	for i, b := range bytes {
		powClean := math.Pow(128, float64(l-i-1))
		pow := int(powClean)
		if i == l-1 {
			result = result + uint32(int(b&127)*pow)
		} else {
			result = result + uint32(int(b)*pow)
		}
	}
	// because of fast spec 1.1
	if optional && result > 0 {
		result -= 1
	}
	return result
}

func UInt64Optional(bytes []byte, optional bool) uint64 {
	l := len(bytes)
	var result uint64 = 0
	for i, b := range bytes {
		powClean := math.Pow(128, float64(l-i-1))
		pow := int(powClean)
		if i == l-1 {
			result = result + uint64(int(b&127)*pow)
		} else {
			result = result + uint64(int(b)*pow)
		}
	}
	// Because of fast spec 1.1
	if optional && result > 0 {
		result -= 1
	}
	return result
}

func String(bytes []byte) string {
	bytes[len(bytes)-1] = bytes[len(bytes)-1] & 127
	return string(bytes)
}

func StringWithLength(iter *MessageIterator) string {
	length := iter.NextN(1)
	return string(iter.NextN(int(length[0] - 1)))
}

func DecimalOptional(exp []byte, man []byte, optional bool) decimal.Decimal {
	return decimal.New(Int64Optional(man, false), Int32Optional(exp, optional))
}

func (iter *MessageIterator) PrintBytes() {
	log.Printf("%+v", iter.bytes)
}
