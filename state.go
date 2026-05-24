package gofront

import (
	"sync"

	"gofront/runtime"
)

type StateValue[T any] struct {
	mu        sync.RWMutex
	value     T
	listeners []func(T)
}

func State[T any](initial T) *StateValue[T] { return &StateValue[T]{value: initial} }

func (s *StateValue[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *StateValue[T]) Set(value T) {
	s.mu.Lock()
	s.value = value
	listeners := append([]func(T){}, s.listeners...)
	s.mu.Unlock()
	for _, listener := range listeners {
		listener(value)
	}
	runtime.RequestRender()
}

func (s *StateValue[T]) Subscribe(listener func(T)) func() {
	s.mu.Lock()
	s.listeners = append(s.listeners, listener)
	s.mu.Unlock()
	return func() {}
}
