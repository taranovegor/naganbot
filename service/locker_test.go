package service

import (
	"fmt"
	"sync"
	"testing"
)

func TestLockerTryLockSameResourceIsExclusive(t *testing.T) {
	l := NewLocker()

	release, ok := l.TryLock("chat-1")
	if !ok {
		t.Fatal("expected to acquire the lock")
	}

	if _, ok := l.TryLock("chat-1"); ok {
		t.Fatal("expected a second lock on the same resource to fail while held")
	}

	release()

	release, ok = l.TryLock("chat-1")
	if !ok {
		t.Fatal("expected the lock to be acquirable again after release")
	}
	release()
}

func TestLockerTryLockDifferentResourcesAreIndependent(t *testing.T) {
	l := NewLocker()

	releaseFirst, ok := l.TryLock("chat-1")
	if !ok {
		t.Fatal("expected to acquire the lock for chat-1")
	}
	defer releaseFirst()

	releaseSecond, ok := l.TryLock("chat-2")
	if !ok {
		t.Fatal("expected chat-2 to be lockable while chat-1 is held")
	}
	releaseSecond()
}

func TestLockerConcurrentTryLockIsExclusivePerResource(t *testing.T) {
	l := NewLocker()

	const attempts = 500

	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		holders    int
		maxHolders int
		acquired   int
	)

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			release, ok := l.TryLock("shared")
			if !ok {
				return
			}
			defer release()

			mu.Lock()
			holders++
			acquired++
			if holders > maxHolders {
				maxHolders = holders
			}
			mu.Unlock()

			mu.Lock()
			holders--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if maxHolders > 1 {
		t.Fatalf("expected at most one concurrent holder of the same resource, got %d", maxHolders)
	}
	if acquired == 0 {
		t.Fatal("expected at least one goroutine to acquire the lock")
	}
}

func TestLockerConcurrentTryLockOnDifferentResourcesAllSucceed(t *testing.T) {
	l := NewLocker()

	const resources = 20

	results := make([]bool, resources)
	var wg sync.WaitGroup

	for i := 0; i < resources; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			release, ok := l.TryLock(fmt.Sprintf("resource-%d", n))
			results[n] = ok
			if ok {
				release()
			}
		}(i)
	}
	wg.Wait()

	for n, ok := range results {
		if !ok {
			t.Fatalf("expected to acquire lock for resource-%d concurrently with the others", n)
		}
	}
}

func TestLockerMapIsEmptyAfterAllLocksReleased(t *testing.T) {
	impl := NewLocker().(*locker)

	const (
		workers   = 50
		resources = 5
	)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			resource := fmt.Sprintf("resource-%d", n%resources)
			for {
				release, ok := impl.TryLock(resource)
				if ok {
					release()
					return
				}
			}
		}(i)
	}
	wg.Wait()

	impl.mu.Lock()
	defer impl.mu.Unlock()
	if len(impl.locks) != 0 {
		t.Fatalf("expected the internal map to be empty, got %d entries", len(impl.locks))
	}
}
