package service

import (
	"errors"
	"sync"
)

var (
	ErrLockNotAcquired = errors.New("resource is already locked")
)

type Locker interface {
	TryLock(resource string) (release func(), ok bool)
}

type lockEntry struct {
	mu  sync.Mutex
	ref int
}

type locker struct {
	mu    sync.Mutex
	locks map[string]*lockEntry
}

func NewLocker() Locker {
	return &locker{
		locks: make(map[string]*lockEntry),
	}
}

func (l *locker) TryLock(resource string) (func(), bool) {
	entry := l.acquireEntry(resource)

	if !entry.mu.TryLock() {
		l.releaseEntry(resource, entry)
		return nil, false
	}

	return func() {
		entry.mu.Unlock()
		l.releaseEntry(resource, entry)
	}, true
}

func (l *locker) acquireEntry(resource string) *lockEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.locks[resource]
	if !ok {
		entry = &lockEntry{}
		l.locks[resource] = entry
	}
	entry.ref++

	return entry
}

func (l *locker) releaseEntry(resource string, entry *lockEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry.ref--
	if entry.ref == 0 {
		delete(l.locks, resource)
	}
}
