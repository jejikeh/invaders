package gomemory

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewPool[int](1024)

	if pool.itemSize != SizeOfAligned[int](1) {
		t.Errorf("itemSize=%d but should be equal to alignedSize=%d for type [%T]", pool.itemSize, SizeOfAligned[int](1), int(0))
	}
}

func TestNewObjectInPool(t *testing.T) {
	pool := NewPool[int](1024)

	x := pool.New(0)
	if x == nil {
		t.Errorf("failed to allocate new object in pool")
	}
}

// @Cleanup: Cleanup this test.
func TestGetObjectFromPool(t *testing.T) {
	pool := NewPool[int](1024)

	x := pool.New(0)
	if x == nil {
		// @Incomplete: Maybe check this in .New()? If object is not allocated, it either a cast issue with memory, or malloc issue, or ovewflow.
		t.Error("failed to allocate new object in pool")
	} else {
		*x = 123
	}

	xFromPool, _ := pool.Get(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 123 {
		t.Errorf("expected %d but got %d", 123, *xFromPool)
	}

	*x = 234
	xFromPool, _ = pool.Get(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}

	y := pool.New(1)
	if y == nil {
		// @Incomplete: Maybe check this in .New()? If object is not allocated, it either a cast issue with memory, or malloc issue, or ovewflow
		t.Error("failed to allocate new object in pool")
	} else {
		*y = 102
	}

	yFromPool, _ := pool.Get(1)
	if yFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *yFromPool != 102 {
		t.Errorf("expected %d but got %d", 102, *yFromPool)
	}

	xFromPool, _ = pool.Get(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}
}
