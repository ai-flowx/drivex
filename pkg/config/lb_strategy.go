package config

import (
	"crypto/rand"
	"hash/fnv"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-flowx/drivex/pkg/mycomdef"
)

var (
	rrIndices = make(map[string]*uint32)
	randLock  = &sync.Mutex{}
	modelLock = &sync.RWMutex{}
)

func getRandomIndex(n int) int {
	randLock.Lock()
	defer randLock.Unlock()

	_max := big.NewInt(int64(n))
	num, _ := rand.Int(rand.Reader, _max)

	return int(num.Int64())
}

func getRoundRobinIndex(modelName string, n int) int {
	modelLock.RLock()
	idx, exists := rrIndices[modelName]
	modelLock.RUnlock()

	if !exists {
		modelLock.Lock()
		if idx, exists = rrIndices[modelName]; !exists { // double check locking
			var newIndex uint32 = 0
			rrIndices[modelName] = &newIndex
			idx = &newIndex
		}
		modelLock.Unlock()
	}

	newIdx := atomic.AddUint32(idx, 1)

	return int(newIdx) % n
}

func getHashIndex(key string, n int) int {
	timestamp := time.Now().Format("2006-01-02 15:04:05.999")

	h := fnv.New32a()
	_, _ = h.Write([]byte(key + timestamp))

	return int(h.Sum32()) % n
}

func GetLBIndex(lbStrategy, key string, length int) int {
	lbs := strings.ToLower(lbStrategy)

	switch lbs {
	case mycomdef.KeynameFirst:
		return 0
	case mycomdef.KeynameRandom, mycomdef.KeynameRand:
		return getRandomIndex(length)
	case mycomdef.KeynameRoundRobin, mycomdef.KeynameRr:
		return getRoundRobinIndex(key, length)
	case mycomdef.KeynameHash:
		return getHashIndex(key, length)
	default:
		return getRandomIndex(length)
	}
}
