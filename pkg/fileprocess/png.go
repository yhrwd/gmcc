package fileprocess

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	cache   = make(map[string][]byte)
	cacheMu sync.RWMutex
)

func ClearCache() {
	cacheMu.Lock()
	cache = make(map[string][]byte)
	cacheMu.Unlock()
}

func ParsePNGRegion(path string) []byte {
	cacheMu.RLock()
	if d, ok := cache[path]; ok {
		cacheMu.RUnlock()
		return d
	}
	cacheMu.RUnlock()

	d, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read err:%v\n", err)
		return nil
	}

	w, h := int(binary.BigEndian.Uint32(d[16:20])), int(binary.BigEndian.Uint32(d[20:24]))
	var idat bytes.Buffer

	for i := 33; i < len(d); {
		l := int(binary.BigEndian.Uint32(d[i : i+4]))
		if string(d[i+4:i+8]) == "IDAT" {
			idat.Write(d[i+8 : i+8+l])
		}
		i += l + 12
	}

	zr, err := zlib.NewReader(&idat)
	if err != nil {
		fmt.Printf("zlib err:%v\n", err)
		return nil
	}
	defer zr.Close()

	raw, err := io.ReadAll(zr)
	if err != nil {
		fmt.Printf("read zlib err:%v\n", err)
		return nil
	}

	stride := w * 4
	pixels := make([]byte, h*stride)

	for y, i := 0, 0; y < h; y++ {
		f := raw[i]
		i++
		for x := 0; x < stride; x++ {
			var a, b, c byte
			if x >= 4 {
				a = pixels[y*stride+x-4]
			}
			if y > 0 {
				b = pixels[(y-1)*stride+x]
			}
			if x >= 4 && y > 0 {
				c = pixels[(y-1)*stride+x-4]
			}
			v := raw[i]
			i++
			switch f {
			case 1:
				v += a
			case 2:
				v += b
			case 3:
				v += byte((int(a) + int(b)) >> 1)
			case 4:
				p := int(a) + int(b) - int(c)
				pa, pb, pc := abs(p-int(a)), abs(p-int(b)), abs(p-int(c))
				if pa <= pb && pa <= pc {
					v += a
				} else if pb <= pc {
					v += b
				} else {
					v += c
				}
			}
			pixels[y*stride+x] = v
		}
	}

	out := make([]byte, 0, 8*8*4)
	for y := 8; y <= 15; y++ {
		for x := 8; x <= 15; x++ {
			i := (y*w + x) * 4
			out = append(out, pixels[i], pixels[i+1], pixels[i+2], pixels[i+3])
		}
	}

	cacheMu.Lock()
	cache[path] = out
	cacheMu.Unlock()

	return out
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func RenderRegion(path string) string {
	r := ParsePNGRegion(path)
	if r == nil {
		return "err"
	}

	var b strings.Builder
	for y := 0; y < 8; y += 2 {
		for x := 0; x < 8; x++ {
			i1 := (y*8 + x) * 4
			i2 := ((y+1)*8 + x) * 4
			if i2 >= len(r) {
				continue
			}
			fmt.Fprintf(&b, "\033[38;2;%d;%d;%d;48;2;%d;%d;%dm▀\033[0m",
				r[i1], r[i1+1], r[i1+2], r[i2], r[i2+1], r[i2+2])
		}
		b.WriteByte('\n')
	}
	return b.String()
}
