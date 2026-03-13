package crypto

import "crypto/cipher"

type cfb8 struct {
	b        cipher.Block
	next     []byte
	tmp      []byte
	encrypt  bool
	blockLen int
}

func NewCFB8(block cipher.Block, iv []byte, encrypt bool) cipher.Stream {
	bs := block.BlockSize()
	next := make([]byte, bs)
	copy(next, iv)
	return &cfb8{
		b:        block,
		next:     next,
		tmp:      make([]byte, bs),
		encrypt:  encrypt,
		blockLen: bs,
	}
}

func (x *cfb8) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("cfb8 output smaller than input")
	}

	for i := 0; i < len(src); i++ {
		in := src[i]
		x.b.Encrypt(x.tmp, x.next)
		out := in ^ x.tmp[0]
		dst[i] = out

		copy(x.next, x.next[1:])
		if x.encrypt {
			x.next[x.blockLen-1] = out
		} else {
			x.next[x.blockLen-1] = in
		}
	}
}
