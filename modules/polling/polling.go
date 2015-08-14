package polling

import (
	"github.com/weisd/log"
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

	log.Info("idx %d t %d", p.Idx, p.Total)

	return p.Idx % p.Total
}

func NewPolling(l int) *Polling {
	if l == 0 {
		return nil
	}
	return &Polling{Idx: 0, Total: l - 1}
}
