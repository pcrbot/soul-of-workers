package game

import (
	"sync"
)

type namedLock struct {
	lock sync.Mutex
	m    map[int64]sync.Locker
}

func (l *namedLock) Get(key int64) sync.Locker {
	l.lock.Lock()
	defer l.lock.Unlock()
	v, ok := l.m[key]
	if !ok {
		v = &sync.Mutex{}
		l.m[key] = v
	}
	return v
}

func (l *namedLock) Forget(key int64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.m, key)
}

var roleLock = namedLock{}
var companyLock = namedLock{}
