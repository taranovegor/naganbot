package service

import (
	"errors"
	"sync"
)

var (
	ErrLockFailed = errors.New("lock timeout")
)

type Locker interface {
	LockFor(resource string) *sync.Mutex
}

type locker struct {
	locks sync.Map
}

func NewLocker() Locker {
	return &locker{}
}

func (l *locker) LockFor(resource string) *sync.Mutex {
	mutex, _ := l.locks.LoadOrStore(resource, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}
