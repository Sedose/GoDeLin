package godelin

type Map[K comparable, V any] struct {
	data map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		data: make(map[K]V),
	}
}

func (m *Map[K, V]) GetOrPut(key K, defaultValue func() V) V {
	val, ok := m.data[key]
	if ok {
		return val
	}
	val = defaultValue()
	m.data[key] = val
	return val
}
