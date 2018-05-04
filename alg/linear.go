package alg

import (
	"math"

	"hash/fnv"

	"github.com/damnever/bitarray"
	log "github.com/sirupsen/logrus"
)

// LpCounter - Linear Probabilistic counter
type LpCounter struct {
	Bits *bitarray.BitArray
}

// NewLpCounter - Creates new LpCounter
func NewLpCounter(maxKb int) *LpCounter {
	return &LpCounter{Bits: bitarray.New(8 * 1024 * maxKb)}
}

// Count - Gets the current value of the bitmap: -size * ln(unset_bits/size)
func (lp *LpCounter) Count() int {
	ratio := float64(lp.Bits.Len()-lp.Bits.Count()) / float64(lp.Bits.Len())
	if ratio == float64(0.0) {
		return lp.Bits.Len()
	}
	return int(-float64(lp.Bits.Len()) * math.Log(ratio))

}

// Inc - Increment the counter for the item
func (lp *LpCounter) Inc(item string) error {
	hasher := fnv.New64a()
	_, err := hasher.Write([]byte(item))
	if err != nil {
		log.WithField("Item", item).Error("failed to hash item")
		return err
	}
	hashed := hasher.Sum64()
	offset := hashed * uint64(lp.Bits.Len())
	_, err = lp.Bits.Put(int(offset), 1)
	if err != nil {
		log.WithField("Offset", offset).Error("failed to increment the counter for offset")
		return err
	}
	return nil
}
