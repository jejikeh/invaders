package gomemory

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

func TestBitSetSetBit(t *testing.T) {
	b := NewBitSet[int]()

	tests := make(map[int]bool, 32)

	for range 32 {
		r := rand.IntN(33)
		testSetBit(b, r, t)
		tests[r] = true
	}

	for i := range 32 {
		inSet, _ := tests[i]
		if b.Has(i) != inSet {
			t.Errorf("for %d expected %t, but got %t in bitset[%v]", i, inSet, b.Has(i), strconv.FormatInt(int64(b.bits), 2))
		}
	}
}

func TestBitSetClear(t *testing.T) {
	b := NewBitSet[uint8]()

	testSetBit(b, 2, t)
	testSetBit(b, 3, t)
	testSetBit(b, 0, t)

	testClearBit(b, 3, t)
	testClearBit(b, 2, t)
	testClearBit(b, 0, t)
}

func testSetBit[T Int](b *BitSet[T], v T, t *testing.T) {
	t.Helper()
	b.Set(v)
	if !b.Has(v) {
		t.Errorf("failed to set %d in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}
}

func testClearBit[T Int](b *BitSet[T], v T, t *testing.T) {
	t.Helper()

	if !b.Has(v) {
		t.Errorf("failed to set %d in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}

	b.Clear(v)
	if b.Has(v) {
		t.Errorf("bit %d was not clear in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}
}
