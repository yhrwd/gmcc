package mcclient

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"gmcc/internal/config"
	"gmcc/internal/entity"
	"gmcc/internal/mcclient/packet"
)

func TestHandleAddEntityReadsLpVec3Velocity(t *testing.T) {
	client := newEntityTestClient()

	data, err := hex.DecodeString(
		"02884f5b21b5ee41dab734409451e3bc78470000000000000000c04d570a3d8000000000000000000000c1f880cee41900c30000",
	)
	if err != nil {
		t.Fatalf("DecodeString() error = %v", err)
	}

	if err := client.handleAddEntity(data); err != nil {
		t.Fatalf("handleAddEntity() error = %v", err)
	}

	got, ok := client.entityTracker.Get(2)
	if !ok {
		t.Fatal("实体未写入跟踪器")
	}

	if !almostEqual(got.Position.X, 0, 1e-12) {
		t.Fatalf("X 坐标错误: got=%f", got.Position.X)
	}
	if !almostEqual(got.Position.Y, -58.68000000715256, 1e-9) {
		t.Fatalf("Y 坐标错误: got=%f", got.Position.Y)
	}
	if !almostEqual(got.Position.Z, 0, 1e-12) {
		t.Fatalf("Z 坐标错误: got=%f", got.Position.Z)
	}

	if !almostEqual(got.Velocity.X, -0.014099981688335483, 1e-12) {
		t.Fatalf("速度 X 错误: got=%f", got.Velocity.X)
	}
	if !almostEqual(got.Velocity.Y, -0.10895440395531952, 1e-12) {
		t.Fatalf("速度 Y 错误: got=%f", got.Velocity.Y)
	}
	if !almostEqual(got.Velocity.Z, 0.006348043703839457, 1e-12) {
		t.Fatalf("速度 Z 错误: got=%f", got.Velocity.Z)
	}
}

func TestHandleTeleportEntityReadsCurrentPacketLayout(t *testing.T) {
	client := newEntityTestClient()
	client.entityTracker.SpawnEntity(7, "minecraft:zombie", [16]byte{}, entity.Position{}, entity.Vector3{})

	data := make([]byte, 0, 64)
	data = append(data, packet.EncodeVarInt(7)...)
	data = append(data, encodeFloat64(1.25)...)
	data = append(data, encodeFloat64(2.5)...)
	data = append(data, encodeFloat64(3.75)...)
	data = append(data, encodeFloat64(0.1)...)
	data = append(data, encodeFloat64(-0.2)...)
	data = append(data, encodeFloat64(0.3)...)
	data = append(data, packet.EncodeFloat32(90)...)
	data = append(data, packet.EncodeFloat32(45)...)
	data = append(data, packet.EncodeBool(true)...)

	if err := client.handleTeleportEntity(data); err != nil {
		t.Fatalf("handleTeleportEntity() error = %v", err)
	}

	time.Sleep(150 * time.Millisecond)

	got, ok := client.entityTracker.Get(7)
	if !ok {
		t.Fatal("实体未写入跟踪器")
	}

	if !almostEqual(got.Position.X, 1.25, 1e-12) {
		t.Fatalf("X 坐标错误: got=%f", got.Position.X)
	}
	if !almostEqual(got.Position.Y, 2.5, 1e-12) {
		t.Fatalf("Y 坐标错误: got=%f", got.Position.Y)
	}
	if !almostEqual(got.Position.Z, 3.75, 1e-12) {
		t.Fatalf("Z 坐标错误: got=%f", got.Position.Z)
	}
}

func TestHandleRemoveEntitiesSupportsByteArrayPayload(t *testing.T) {
	client := newEntityTestClient()
	client.entityTracker.SpawnEntity(1, "minecraft:player", [16]byte{1}, entity.Position{}, entity.Vector3{})
	client.entityTracker.SpawnEntity(300, "minecraft:zombie", [16]byte{2}, entity.Position{}, entity.Vector3{})

	payload := make([]byte, 0, 4)
	payload = append(payload, packet.EncodeVarInt(1)...)
	payload = append(payload, packet.EncodeVarInt(300)...)

	if err := client.handleRemoveEntities(packet.EncodeByteArray(payload)); err != nil {
		t.Fatalf("handleRemoveEntities() error = %v", err)
	}

	if _, ok := client.entityTracker.Get(1); ok {
		t.Fatal("实体 1 未被移除")
	}
	if _, ok := client.entityTracker.Get(300); ok {
		t.Fatal("实体 300 未被移除")
	}
}

func TestHandleRemoveEntitiesFallsBackToCountedLayout(t *testing.T) {
	client := newEntityTestClient()
	client.entityTracker.SpawnEntity(1, "minecraft:player", [16]byte{1}, entity.Position{}, entity.Vector3{})
	client.entityTracker.SpawnEntity(300, "minecraft:zombie", [16]byte{2}, entity.Position{}, entity.Vector3{})

	data := make([]byte, 0, 4)
	data = append(data, packet.EncodeVarInt(2)...)
	data = append(data, packet.EncodeVarInt(1)...)
	data = append(data, packet.EncodeVarInt(300)...)

	if err := client.handleRemoveEntities(data); err != nil {
		t.Fatalf("handleRemoveEntities() error = %v", err)
	}

	if _, ok := client.entityTracker.Get(1); ok {
		t.Fatal("实体 1 未被移除")
	}
	if _, ok := client.entityTracker.Get(300); ok {
		t.Fatal("实体 300 未被移除")
	}
}

func newEntityTestClient() *Client {
	cfg := &config.Config{
		Account: config.AccountConfig{
			PlayerID: "TestPlayer",
		},
		Log: config.LogConfig{},
	}

	client := New(cfg)
	client.entityTracker = entity.NewTracker()
	return client
}

func encodeFloat64(v float64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], math.Float64bits(v))
	return b[:]
}

func almostEqual(got, want, epsilon float64) bool {
	return math.Abs(got-want) <= epsilon
}
