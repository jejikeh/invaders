package gomemory

import (
	"testing"
)

func TestNewObjectInPool(t *testing.T) {
	t.Parallel()

	pool := NewTypedPool[int](1024)
	x := pool.NewAt(0)

	if x == nil {
		t.Errorf("failed to allocate new object in pool")
	}
}

// @Cleanup: Cleanup this test.
func TestGetObjectFromPool(t *testing.T) {
	t.Parallel()

	pool := NewTypedPool[int](1024)
	x := pool.NewAt(0)
	*x = 123

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
	*y = 102

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
	t.Parallel()

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
	t.Parallel()

	untypedPool := NewPool[int](16)
	typedPool := ToTypedPool[int](untypedPool)

	if typedPool == nil {
		t.Fail()
	}
}
