package ppic

import "image/png"

// encodeBufferPool is a simple png.EncoderBufferPool containing a single png.EncoderBuffer
type encoderBufferPool struct {
	b *png.EncoderBuffer
}

func (p *encoderBufferPool) Get() *png.EncoderBuffer {
	return p.b
}

func (p *encoderBufferPool) Put(b *png.EncoderBuffer) {
	p.b = b
}
