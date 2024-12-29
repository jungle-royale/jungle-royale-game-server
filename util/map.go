package util

import (
	"sync"
)

type Map[Key comparable, Value any] struct {
	internal_map      map[Key]Value
	internal_sync_map sync.Map
	sync              bool
}

// NewMap: 동기화되지 않은 Map 생성
func NewMap[Key comparable, Value any]() *Map[Key, Value] {
	return &Map[Key, Value]{
		internal_map: make(map[Key]Value),
		sync:         false,
	}
}

// NewSyncMap: 동기화된 Map 생성
func NewSyncMap[Key comparable, Value any]() *Map[Key, Value] {
	return &Map[Key, Value]{
		sync: true,
	}
}

// Get: 키에 해당하는 값을 반환
func (m *Map[Key, Value]) Get(key Key) (*Value, bool) {
	if m.sync {
		value, ok := m.internal_sync_map.Load(key)
		if ok {
			v := value.(Value) // 타입 단언 필요
			return &v, true
		}
		return nil, false
	} else {
		value, ok := m.internal_map[key]
		if ok {
			return &value, true
		}
		return nil, false
	}
}

// Store: 키-값 저장
func (m *Map[Key, Value]) Store(key Key, value Value) {
	if m.sync {
		m.internal_sync_map.Store(key, value)
	} else {
		m.internal_map[key] = value
	}
}

// Update: 기존 키의 값을 업데이트
func (m *Map[Key, Value]) Update(key Key, value Value) bool {
	if m.sync {
		if _, ok := m.internal_sync_map.Load(key); ok {
			m.internal_sync_map.Store(key, value)
			return true
		}
		return false
	} else {
		if _, ok := m.internal_map[key]; ok {
			m.internal_map[key] = value
			return true
		}
		return false
	}
}

// Delete: 키 삭제
func (m *Map[Key, Value]) Delete(key Key) {
	if m.sync {
		m.internal_sync_map.Delete(key)
	} else {
		delete(m.internal_map, key)
	}
}

// Range: Map 순회
func (m *Map[Key, Value]) Range(f func(Key, Value) bool) {
	if m.sync {
		m.internal_sync_map.Range(func(key, value any) bool {
			return f(key.(Key), value.(Value)) // 타입 단언 필요
		})
	} else {
		for k, v := range m.internal_map {
			if !f(k, v) {
				break
			}
		}
	}
}

// KeyList: 모든 키를 슬라이스로 반환
func (m *Map[Key, Value]) KeyList() []Key {
	keys := []Key{}
	if m.sync {
		m.internal_sync_map.Range(func(key, _ any) bool {
			keys = append(keys, key.(Key)) // 타입 단언 필요
			return true
		})
	} else {
		for k := range m.internal_map {
			keys = append(keys, k)
		}
	}
	return keys
}

// ValueList: 모든 값을 슬라이스로 반환
func (m *Map[Key, Value]) ValueList() []Value {
	values := []Value{}
	if m.sync {
		m.internal_sync_map.Range(func(_, value any) bool {
			values = append(values, value.(Value)) // 타입 단언 필요
			return true
		})
	} else {
		for _, v := range m.internal_map {
			values = append(values, v)
		}
	}
	return values
}