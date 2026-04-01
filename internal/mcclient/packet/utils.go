package packet

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func ParseAddress(addr string) (string, uint16, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0, fmt.Errorf("server.address 不能为空")
	}

	if !strings.Contains(addr, ":") {
		return addr, 25565, nil
	}

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		idx := strings.LastIndex(addr, ":")
		if idx <= 0 {
			return "", 0, fmt.Errorf("无效地址: %s", addr)
		}
		host = addr[:idx]
		portStr = addr[idx+1:]
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt <= 0 || portInt > 65535 {
		return "", 0, fmt.Errorf("无效端口: %s", portStr)
	}
	return host, uint16(portInt), nil
}

func OfflineUUID(name string) [16]byte {
	hash := md5.Sum([]byte("OfflinePlayer:" + name))
	hash[6] = (hash[6] & 0x0F) | 0x30
	hash[8] = (hash[8] & 0x3F) | 0x80
	return hash
}

func ParseUUID(raw string) ([16]byte, error) {
	clean := strings.ReplaceAll(strings.TrimSpace(raw), "-", "")
	if len(clean) != 32 {
		return [16]byte{}, fmt.Errorf("uuid 长度无效")
	}
	b, err := hex.DecodeString(clean)
	if err != nil {
		return [16]byte{}, err
	}
	var id [16]byte
	copy(id[:], b)
	return id, nil
}

func FormatUUID(id [16]byte) string {
	hexStr := hex.EncodeToString(id[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8],
		hexStr[8:12],
		hexStr[12:16],
		hexStr[16:20],
		hexStr[20:32],
	)
}

func FormatUUIDShort(id [16]byte) string {
	return hex.EncodeToString(id[:4])
}

func MinecraftServerHash(serverID string, sharedSecret []byte, publicKey []byte) string {
	h := sha1.New()
	_, _ = h.Write([]byte(serverID))
	_, _ = h.Write(sharedSecret)
	_, _ = h.Write(publicKey)
	sum := h.Sum(nil)

	if len(sum) == 0 {
		return "0"
	}
	if sum[0]&0x80 != 0 {
		neg := make([]byte, len(sum))
		for i := range sum {
			neg[i] = ^sum[i]
		}
		for i := len(neg) - 1; i >= 0; i-- {
			neg[i]++
			if neg[i] != 0 {
				break
			}
		}
		n := new(big.Int).SetBytes(neg)
		if n.Sign() == 0 {
			return "0"
		}
		return "-" + n.Text(16)
	}

	n := new(big.Int).SetBytes(sum)
	if n.Sign() == 0 {
		return "0"
	}
	return n.Text(16)
}

func RawPreview(b []byte) string {
	if len(b) == 0 {
		return "<empty>"
	}
	const max = 120
	if len(b) <= max {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:max]) + "..."
}

func ReadFloat32FromBytes(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}
