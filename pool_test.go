package gomemory

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	pool := NewTypedPool[int](1024)

	if pool.itemSize != SizeOfAligned[int](1) {
		t.Errorf("itemSize=%d but should be equal to alignedSize=%d for type [%T]", pool.itemSize, SizeOfAligned[int](1), int(0))
	}
}

func TestNewObjectInPool(t *testing.T) {
	pool := NewTypedPool[int](1024)

	x := pool.NewAt(0)
	if x == nil {
		t.Errorf("failed to allocate new object in pool")
	}
}

// @Cleanup: Cleanup this test.
func TestGetObjectFromPool(t *testing.T) {
	pool := NewTypedPool[int](1024)

	x := pool.NewAt(0)
	if x == nil {
		// @Incomplete: Maybe check this in .NewAt()? If object is not allocated, it either a cast issue with memory, or malloc issue, or ovewflow.
		t.Error("failed to allocate new object in pool")
	} else {
		*x = 123
	}

	xFromPool, _ := pool.GetAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 123 {
		t.Errorf("expected %d but got %d", 123, *xFromPool)
	}

	*x = 234
	xFromPool, _ = pool.GetAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}

	y := pool.NewAt(1)
	if y == nil {
		// @Incomplete: Maybe check this in .NewAt()? If object is not allocated, it either a cast issue with memory, or malloc issue, or ovewflow
		t.Error("failed to allocate new object in pool")
	} else {
		*y = 102
	}

	yFromPool, _ := pool.GetAt(1)
	if yFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *yFromPool != 102 {
		t.Errorf("expected %d but got %d", 102, *yFromPool)
	}

	xFromPool, _ = pool.GetAt(0)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 234 {
		t.Errorf("expected %d but got %d", 234, *xFromPool)
	}
}

func TestPool(t *testing.T) {
	untypedPool := NewPool[int](16)
	
	x := (*int)(untypedPool.NewAt(1))
	*x = 123
	
	ptr, _ := untypedPool.GetAt(1)
	xFromPool := (*int)(ptr)
	if xFromPool == nil {
		t.Error("failed to get object from pool by index")
	} else if *xFromPool != 123 {
		t.Errorf("expected %d but got %d", 123, *xFromPool)
	}
}

func TestToTypedPool(t *testing.T) {
	untypedPool := NewPool[int](16)
	
	typedPool := ToTypedPool[int](untypedPool)
	
	if typedPool == nil {
		t.Fail()
	}
}