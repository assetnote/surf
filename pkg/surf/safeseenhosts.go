package surf

import (
	"errors"
	"sync"
)

type SafeSeenHosts struct {
	sync.RWMutex
	hosts map[string]bool
}

func (sn *SafeSeenHosts) Add(host string) {
	sn.Lock()
	defer sn.Unlock()
	sn.hosts[host] = true
}

func (sn *SafeSeenHosts) Get(host string) (bool, error) {
	sn.RLock()
	defer sn.RUnlock()
	if seen, ok := sn.hosts[host]; ok {
		return seen, nil
	}
	return false, errors.New("Host does not exist")
}
