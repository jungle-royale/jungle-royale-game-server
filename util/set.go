package util

import "sync"

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
		for k, _ := range set.internal_set {
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
