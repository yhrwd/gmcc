package mcclient

import (
	"testing"

	"gmcc/internal/config"
	"gmcc/internal/mcclient/packet"
)

func TestGetOnlinePlayers(t *testing.T) {
	cfg := &config.Config{
		Account: config.AccountConfig{
			PlayerID: "TestPlayer",
		},
		Log: config.LogConfig{},
	}
	client := New(cfg)

	if !client.IsReady() {
		t.Log("客户端未连接")
	}

	players := client.GetOnlinePlayers()
	if players != nil {
		t.Logf("在线玩家数量: %d", len(players))
		for i, player := range players {
			t.Logf("玩家 %d: %s", i+1, player)
		}
	}
}

func TestPlayerInfoCache(t *testing.T) {
	cfg := &config.Config{
		Account: config.AccountConfig{
			PlayerID: "TestPlayer",
		},
		Log: config.LogConfig{},
	}
	client := New(cfg)

	client.playersMu.Lock()
	client.players["testPlayer1"] = playerInfo{uuid: [16]byte{1, 2, 3, 4}}
	client.players["testPlayer2"] = playerInfo{uuid: [16]byte{5, 6, 7, 8}}
	client.playersMu.Unlock()

	players := client.GetOnlinePlayers()
	if len(players) != 2 {
		t.Errorf("期望 2 个玩家，实际得到 %d 个", len(players))
	}

	playerSet := make(map[string]bool)
	for _, player := range players {
		playerSet[player] = true
	}
	if !playerSet["testPlayer1"] || !playerSet["testPlayer2"] {
		t.Error("玩家列表不完整")
	}
}

func TestFormatUUIDShort(t *testing.T) {
	uuid := [16]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	result := formatUUIDShort(uuid)
	expected := "12345678"

	if result != expected {
		t.Errorf("期望 '%s'，实际得到 '%s'", expected, result)
	}
}

func TestOfflineUUID(t *testing.T) {
	name := "TestPlayer"
	uuid := packet.OfflineUUID(name)

	if uuid != [16]byte{} {
		t.Logf("离线 UUID: %x", uuid)
	}
}
