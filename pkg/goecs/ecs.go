package goecs

import (
	"reflect"

	"github.com/jejikeh/invaders/pkg/gomemory"
)

// @Incomplete: For now, there are no way to delete entities in the layer.
// For deleting, it might be needed to implement some sort of free-list tracking
// in arena allocators inside gomemory.
// @Incomplete: For now, there are no way to resize component and entity pool, so the

// Max count of components and entities are actually limited by bit size of ComponentID and EntityID

const MaxEntityCount = 1024

type ComponentID int64

type EntityID int64

type EntityInfo struct {
	ID            EntityID
	ComponentMask gomemory.BitSet[ComponentID]
}

type System = func(*Layer)

type Layer struct {
	componentPool []*gomemory.Pool[int, any]
	componentIDs  map[string]ComponentID

	// @Incomplete: I use TypedPool here because in the future the entities can be removed, though it can be replaced by Arena in the future
	entities *gomemory.Pool[int, EntityInfo]
	// @Incomplete: For now, systems cannot be removed or disabled at runtime
	systems []System
}

func NewLayer(systems ...System) *Layer {
	return &Layer{
		// @Incomplete: We could store this data manually to give complete ownership of memory to Layer?
		componentPool: make([]*gomemory.Pool[int, any], gomemory.Sizeof(ComponentID(0))),
		componentIDs:  make(map[string]ComponentID, gomemory.Sizeof(ComponentID(0))),
		entities:      gomemory.NewPool[int, EntityInfo](MaxEntityCount),
		systems:       systems,
	}
}

func (l *Layer) NewEntity() EntityID {
	id := l.entities.Length()
	entity := l.entities.StoreAt(id)
	entity.ID = EntityID(id)

	return EntityID(id)
}

func (l *Layer) AddSystems(systems ...System) {
	l.systems = append(l.systems, systems...)
}

func Attach[T any](layer *Layer, entityID EntityID) *T {
	entity, ok := layer.entities.LoadAt(int(entityID))
	if !ok {
		return nil
	}

	componentID := GetComponentID[T](layer)
	if layer.componentPool[componentID] == nil {
		// @Cleanup: https://www.david-colson.com/2020/02/09/making-a-simple-ecs.html
		// Potential improvements and alternatives
		layer.componentPool[componentID] = gomemory.NewPool[int, any](MaxEntityCount, new(T))
	}

	componentPool := layer.componentPool[componentID]
	component := componentPool.StoreAt(int(entityID), func(a *any) {
		*a = new(T)
	})

	entity.ComponentMask.Set(componentID)

	return (*component).(*T)
}

// @Cleanup: In componentPools the are still allocated memory for this entityID.
func Detach[T any](layer *Layer, entityID EntityID) {
	entity, ok := layer.entities.LoadAt(int(entityID))
	if !ok {
		return
	}

	componentID := GetComponentID[T](layer)
	entity.ComponentMask.Unset(componentID)
}

func GetComponent[T any](layer *Layer, entityID EntityID) (*T, bool) {
	entity, entityExists := layer.entities.LoadAt(int(entityID))
	if !entityExists {
		return nil, false
	}

	componentID := GetComponentID[T](layer)
	if !entity.ComponentMask.Has(componentID) {
		return nil, false
	}

	componentPool := layer.componentPool[componentID]

	memPtr, _ := componentPool.LoadAt(int(entityID))

	return (*memPtr).(*T), true
}

func HasComponent[T any](layer *Layer, entityID EntityID) bool {
	entity, entityExists := layer.entities.LoadAt(int(entityID))
	if !entityExists {
		return false
	}

	return entity.ComponentMask.Has(GetComponentID[T](layer))
}

func GetComponentID[T any](layer *Layer) ComponentID {
	// @Incomplete: Maybe this could be better? Maybe let user assign their own ID?
	// So, the problem is what this allocate new object every time. And in ideal world we
	// could do that once at startup or even at compile time with some hardcoded ID. Why not?
	componentType := reflect.TypeOf(new(T)).String()
	if id, ok := layer.componentIDs[componentType]; ok {
		return id
	}

	layer.componentIDs[componentType] = ComponentID(len(layer.componentIDs))

	return layer.componentIDs[componentType]
}

func (l *Layer) Update() {
	for _, system := range l.systems {
		system(l)
	}
}

// @Cleanup: When Go 1.23 comes out, rewrite this all to iterators.
// @Incomplete: For now it is very stupid straightforward implementation.

// @Incomplete: We could have some state in systems, for cache component IDs and etc.
func (l *Layer) Request(components ...ComponentID) []EntityID {
	var entities []EntityID

	mask := gomemory.NewBitSet[ComponentID](components...)

	for i := range l.entities.Length() {
		if entity, ok := l.entities.LoadAt(i); ok && entity.ComponentMask.Check(mask) {
			entities = append(entities, entity.ID)
		}
	}

	return entities
}
