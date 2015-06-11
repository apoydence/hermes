package handlers

type Locker struct {
	ch chan int
}

func NewLocker() *Locker {
	l := &Locker{
		ch: make(chan int, 1),
	}

	l.ch <- 1
	return l
}

func (l *Locker) TryLock() bool {
	select {
	case <-l.ch:
		return true
	default:
		return false
	}
}

func (l *Locker) Unlock() {
	select {
	case <-l.ch:
		panic("Lock already unlocked")
	default:
		l.ch <- 1
	}
}
