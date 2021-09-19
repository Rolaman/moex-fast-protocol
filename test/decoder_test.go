package test

import (
	"github.com/kdt-wolf/moex-fast/internal/decoder"
	"testing"
)

func TestPmap(t *testing.T) {
	iterator := decoder.NewIterator([]byte{158})
	pmap := iterator.NextPmap()
	next1 := pmap.HasNext()
	next2 := pmap.HasNext()
	next3 := pmap.HasNext()
	if !next1 && !next2 && next3 {
		return
	}
	t.Fatal("Wrong pmap")
}
