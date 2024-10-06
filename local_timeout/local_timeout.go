package local_timeout

import (
	"banyan/config"
	"sync"
	"time"
)

type LocalTimeout struct {
	curHeight     int
	newHeightChan chan int
	mu            sync.Mutex
}

func NewLocalTimeout() *LocalTimeout {
	lt := new(LocalTimeout)
	lt.curHeight = 1
	lt.newHeightChan = make(chan int, 100)
	return lt
}

func (lt *LocalTimeout) HeightIncreased(block_production_height int) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	if block_production_height <= lt.curHeight {
		return
	}
	lt.curHeight = block_production_height
	lt.newHeightChan <- block_production_height // reset timer for the next view
}

func (lt *LocalTimeout) GetNewHeight() chan int {
	return lt.newHeightChan
}

func (lt *LocalTimeout) GetTimeoutDuration() time.Duration {
	return time.Duration(config.GetConfig().Timeout) * time.Millisecond
}
