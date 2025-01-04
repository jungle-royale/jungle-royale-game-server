package util

import (
	"math/rand"
	"sync"
)

type Set[T comparable] struct {
	internal_set      map[T]struct{}
	internal_sync_set sync.Map
	sync              bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		internal_set: make(map[T]struct{}),
		sync:         false,
	}
}

func NewSyncSet[T comparable]() *Set[T] {
	return &Set[T]{
		sync: true,
	}
}

func (set *Set[T]) Add(v T) {
	if set.sync {
		set.internal_sync_set.Store(v, struct{}{})
	} else {
		set.internal_set[v] = struct{}{}
	}
}

func (set *Set[T]) Remove(v T) {
	if set.sync {
		set.internal_sync_set.Delete(v)
	} else {
		delete(set.internal_set, v)
	}
}

func (set *Set[T]) Contain(v T) bool {
	if set.sync {
		_, exists := set.internal_sync_set.Load(v)
		return exists
	} else {
		_, exists := set.internal_set[v]
		return exists
	}
}

func (set *Set[T]) Range(f func(T) bool) {
	if set.sync {
		set.internal_sync_set.Range(func(key, value any) bool {
			return f(key.(T))
		})
	} else {
		for k := range set.internal_set {
			if !f(k) {
				break
			}
		}
	}
}

func (set *Set[T]) Difference(other *Set[T]) *Set[T] {
	ret := NewSet[T]()
	set.Range(func(t T) bool {
		if !other.Contain(t) {
			ret.Add(t)
		}
		return true
	})
	return ret
}

func (set *Set[T]) Length() int {
	if set.sync {
		len := 0
		set.internal_sync_set.Range(func(key, value any) bool {
			len++
			return true
		})
		return len
	} else {
		return len(set.internal_set)
	}
}

func (set *Set[T]) KeyList(f func(T) bool) []T {
	if f == nil {
		f = func(v T) bool {
			return true
		}
	}
	keys := []T{}
	if set.sync {
		set.internal_sync_set.Range(func(key, value any) bool {
			if f(value.(T)) {
				keys = append(keys, key.(T)) // 타입 단언 필요
			}
			return true
		})
	} else {
		for k := range set.internal_set {
			keys = append(keys, k)
		}
	}
	return keys
}

func (set *Set[T]) SelectRandom(f func(T) bool) (T, bool) {
	var zeroKey T

	// 모든 키 목록을 가져옴
	keys := set.KeyList(f)
	if len(keys) == 0 {
		// 맵이 비어있으면 false 반환
		return zeroKey, false
	}

	// 무작위 인덱스 선택
	idx := rand.Intn(len(keys))
	chosenKey := keys[idx]

	return chosenKey, true
}
