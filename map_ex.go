package godelin

type MapEx[K comparable, V any] struct {
	data map[K]V
}

func NewMapEx[K comparable, V any]() MapEx[K, V] {
	return MapEx[K, V]{
		data: make(map[K]V),
	}
}

func (m *MapEx[K, V]) GetOrPut(key K, defaultValue func() V) V {
	val, ok := m.data[key]
	if ok {
		return val
	}
	val = defaultValue()
	m.data[key] = val
	return val
}
