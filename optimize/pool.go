package optimize

import (
	"image/png"
	"sync"
)

type pngBufferPool struct {
	p *sync.Pool
}

func NewBufferPool() *pngBufferPool {
	p := &pngBufferPool{
		p: &sync.Pool{
			New: func() any {
				return &png.EncoderBuffer{}
			},
		},
	}
	for i := 0; i < 8; i++ {
		p.Put(&png.EncoderBuffer{})
	}
	return p
}

func (p *pngBufferPool) Get() *png.EncoderBuffer {
	e, _ := p.p.Get().(*png.EncoderBuffer)
	return e
}

func (p *pngBufferPool) Put(b *png.EncoderBuffer) {
	p.p.Put(b)
}
