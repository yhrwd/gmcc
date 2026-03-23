package binutil

// VarInt 是可变长度整数类型
type VarInt int32

// VarLong 是可变长度长整数类型
type VarLong int64

// Position 是三维位置坐标
type Position struct {
	X, Y, Z int32
}

// UUID 表示16字节UUID
type UUID [16]byte

// BitSet 是位集合类型
type BitSet []int64
