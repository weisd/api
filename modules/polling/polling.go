package polling

import (
	"math"
)

type Polling struct {
	Idx   int
	Total int
}

func (p *Polling) Index() int {

	p.Idx++
	if p.Idx >= math.MaxInt32 {
		p.Idx = 0
	}

	return p.Idx % p.Total
}

func NewPolling(l int) *Polling {
	if l == 0 {
		return nil
	}
	return &Polling{Idx: 0, Total: l}
}
