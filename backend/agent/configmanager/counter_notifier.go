package cm

import "sync"

type ShutdownController struct {
	isActive    bool
	allchannels map[int]chan struct{}
	counter     int
	C           chan struct{}
	lock        sync.Mutex
}

func (s *ShutdownController) IsActive() bool {
	return s.isActive
}

func (s *ShutdownController) Shutdown() {
	s.isActive = false
	for _, v := range s.allchannels {
		v <- struct{}{}
	}
}

func (s *ShutdownController) Subsribe() (<-chan struct{}, func()) {
	s.lock.Lock()
	k := s.counter
	s.counter++
	s.allchannels[k] = make(chan struct{}, 1)
	s.lock.Unlock()
	return s.allchannels[k], func() {
		delete(s.allchannels, k)
		if len(s.allchannels) == 0 {
			s.C <- struct{}{}
		}
	}
}

func NewShutdownController() *ShutdownController {
	return &ShutdownController{
		isActive:    true,
		allchannels: make(map[int]chan struct{}),
		C:           make(chan struct{}, 1),
	}
}
