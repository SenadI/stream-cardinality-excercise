package alg

import (
	"hash/fnv"
)

func getHash(value string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(value))
	return h.Sum32()
}

// Counter that stors unique users per minute. The main structure is a HyperLogLog
// counter and the starting epoh time - The first incoming message. This implementation
// assumes messages are ordered.
type Counter struct {
	Hll *HyperLogLog
	// Now - artifical now of the incoming data reprents the first timestamp of the open interval.
	Now int
}

// Count unique uids from one minute timeframe.
func (c *Counter) Count(uid string, ts int) (bool, int) {
	// this is helpful for first run. Counter is started from 0 time
	// which is a signal that this is a first iteration.
	if c.Now == 0 {
		c.Now = ts
		c.Hll = NewHyperLogLog(0.1)
		c.Hll.Add(getHash(uid))
		return false, 0
		// Now if the ts is in the current minute range
		// Increment the counter
	} else if c.Now <= ts && c.Now+60 > ts {
		c.Hll.Add(getHash(uid))
		return false, 0
		// Now new one minute timeframe is here so report the curent count
	} else {
		uniqueCount := c.Hll.Count()
		c.Now = ts
		c.Hll = NewHyperLogLog(0.01)
		return true, int(uniqueCount)
	}
}
