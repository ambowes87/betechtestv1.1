package notifications

import "sync"

type Subscriber struct {
	id            string
	mutex         sync.RWMutex
	topics        map[string]bool
	notifications chan notification
}

type subscribers map[string]*Subscriber

func (s *Subscriber) Push(n notification) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.notifications <- n
}
