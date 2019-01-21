package tyut_osc

import (
	"github.com/otiai10/gosseract"
	"go.uber.org/atomic"
	"sync"
)

type OcrEnginePool struct {
	active *atomic.Int32
	size   *atomic.Int32
	p      sync.Pool
}

func (p *OcrEnginePool) SetSize(size int32) {
	p.size.Store(size)
}

func (p *OcrEnginePool) Get() *gosseract.Client {
	client := p.p.Get().(*gosseract.Client)
	p.active.Add(1)
	return client
}

func (p *OcrEnginePool) Put(c *gosseract.Client) {
	go func() {
		if p.active.Load() > p.size.Load() {
			c.Close()
			c = nil //GC
		} else {
			p.p.Put(c)
		}
	}()
}

func NewOcrEnginePool(size int32, initialActive int32) *OcrEnginePool {
	if size < 0 {
		return nil
	}

	sp := sync.Pool{
		New: func() interface{} {
			return gosseract.NewClient()
		},
	}

	ocrpool := OcrEnginePool{
		active: atomic.NewInt32(0),
		size:   atomic.NewInt32(size),
		p:      sp,
	}

	// Create initial engines
	if initialActive > 0 {
		for i := 0; i < int(initialActive); i++ {
			ocrpool.Put(gosseract.NewClient())
		}
	}

	return &ocrpool
}
