package player

import (
	"testing"
	"time"

	"gmcc/internal/entity"
)

func TestNewNearbyTracker(t *testing.T) {
	// 创建模拟的查找函数
	lookup := func(uuid [16]byte) (string, bool) {
		return "TestPlayer", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	if nt == nil {
		t.Fatal("期望NearbyTracker不为nil")
	}

	if nt.entityTracker != tracker {
		t.Error("期望entityTracker正确设置")
	}

	if nt.lookupPlayer == nil {
		t.Error("期望lookupPlayer正确设置")
	}
}

func TestNearbyTracker_HandlePlayerSpawn(t *testing.T) {
	// 创建模拟的查找函数
	lookup := func(uuid [16]byte) (string, bool) {
		return "TestPlayer", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	var enterCalled bool
	var enteredPlayer *NearbyPlayer
	nt.SetCallbacks(PlayerCallbacks{
		OnPlayerEnter: func(p *NearbyPlayer) {
			enterCalled = true
			enteredPlayer = p
		},
	})

	// 生成玩家实体
	uuid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pos := entity.Position{X: 100.0, Y: 64.0, Z: -50.0}
	vel := entity.Vector3{X: 0, Y: 0, Z: 0}

	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	if !enterCalled {
		t.Error("期望OnPlayerEnter回调被调用")
	}

	if enteredPlayer == nil {
		t.Fatal("期望enteredPlayer不为nil")
	}

	if enteredPlayer.Username != "TestPlayer" {
		t.Errorf("期望Username=TestPlayer，实际%s", enteredPlayer.Username)
	}

	if enteredPlayer.Entity.ID != 1 {
		t.Errorf("期望Entity.ID=1，实际%d", enteredPlayer.Entity.ID)
	}

	// 测试能否通过NearbyTracker获取
	players := nt.GetNearbyPlayers()
	if len(players) != 1 {
		t.Errorf("期望找到1个玩家，实际%d", len(players))
	}

	// 测试通过UUID查找
	found, ok := nt.GetNearbyPlayer(uuid)
	if !ok || found == nil {
		t.Error("期望能通过UUID找到玩家")
	}
}

func TestNearbyTracker_HandleNonPlayerSpawn(t *testing.T) {
	// 创建模拟的查找函数
	lookup := func(uuid [16]byte) (string, bool) {
		return "", false
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	var enterCalled bool
	nt.SetCallbacks(PlayerCallbacks{
		OnPlayerEnter: func(p *NearbyPlayer) {
			enterCalled = true
		},
	})

	// 生成非玩家实体（僵尸）
	uuid := [16]byte{1}
	pos := entity.Position{X: 0, Y: 0, Z: 0}
	vel := entity.Vector3{X: 0, Y: 0, Z: 0}

	tracker.SpawnEntity(1, "minecraft:zombie", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	if enterCalled {
		t.Error("期望OnPlayerEnter回调不被调用（非玩家实体）")
	}

	// 验证玩家列表为空
	players := nt.GetNearbyPlayers()
	if len(players) != 0 {
		t.Errorf("期望玩家列表为空，实际%d", len(players))
	}
}

func TestNearbyTracker_HandlePlayerMove(t *testing.T) {
	lookup := func(uuid [16]byte) (string, bool) {
		return "TestPlayer", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	var moveCalled bool
	var movedPlayer *NearbyPlayer
	var oldPos entity.Position
	nt.SetCallbacks(PlayerCallbacks{
		OnPlayerMove: func(p *NearbyPlayer, old entity.Position) {
			moveCalled = true
			movedPlayer = p
			oldPos = old
		},
	})

	// 生成玩家
	uuid := [16]byte{1}
	pos := entity.Position{X: 0, Y: 64, Z: 0}
	vel := entity.Vector3{X: 0, Y: 0, Z: 0}
	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	// 更新位置
	newPos := entity.Position{X: 10, Y: 64, Z: 5}
	tracker.UpdatePosition(1, newPos)

	time.Sleep(150 * time.Millisecond)

	if !moveCalled {
		t.Error("期望OnPlayerMove回调被调用")
	}

	if movedPlayer == nil {
		t.Fatal("期望movedPlayer不为nil")
	}

	if oldPos.X != 0 || oldPos.Z != 0 {
		t.Errorf("期望旧位置为(0,64,0)，实际(%f,%f,%f)", oldPos.X, oldPos.Y, oldPos.Z)
	}
}

func TestNearbyTracker_HandlePlayerLeave(t *testing.T) {
	lookup := func(uuid [16]byte) (string, bool) {
		return "TestPlayer", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	var leaveCalled bool
	var leftPlayer *NearbyPlayer
	nt.SetCallbacks(PlayerCallbacks{
		OnPlayerLeave: func(p *NearbyPlayer) {
			leaveCalled = true
			leftPlayer = p
		},
	})

	// 生成玩家
	uuid := [16]byte{1}
	pos := entity.Position{X: 0, Y: 0, Z: 0}
	vel := entity.Vector3{X: 0, Y: 0, Z: 0}
	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	// 移除玩家
	tracker.RemoveEntity(1)

	time.Sleep(50 * time.Millisecond)

	if !leaveCalled {
		t.Error("期望OnPlayerLeave回调被调用")
	}

	if leftPlayer == nil {
		t.Fatal("期望leftPlayer不为nil")
	}

	if leftPlayer.Username != "TestPlayer" {
		t.Errorf("期望Username=TestPlayer，实际%s", leftPlayer.Username)
	}

	// 验证玩家已移除
	if nt.Count() != 0 {
		t.Errorf("期望玩家数量为0，实际%d", nt.Count())
	}
}

func TestNearbyTracker_PlayersWithinDistance(t *testing.T) {
	lookup := func(uuid [16]byte) (string, bool) {
		return "Player", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	// 生成3个玩家在不同位置
	uuid1 := [16]byte{1}
	uuid2 := [16]byte{2}
	uuid3 := [16]byte{3}

	tracker.SpawnEntity(1, "minecraft:player", uuid1, entity.Position{X: 0, Y: 64, Z: 0}, entity.Vector3{})
	tracker.SpawnEntity(2, "minecraft:player", uuid2, entity.Position{X: 5, Y: 64, Z: 0}, entity.Vector3{})  // 距离5
	tracker.SpawnEntity(3, "minecraft:player", uuid3, entity.Position{X: 20, Y: 64, Z: 0}, entity.Vector3{}) // 距离20

	time.Sleep(50 * time.Millisecond)

	// 查找距离10以内的玩家
	center := entity.Position{X: 0, Y: 64, Z: 0}
	nearby := nt.PlayersWithinDistance(center, 10.0)

	if len(nearby) != 2 {
		t.Errorf("期望找到2个玩家（距离<=10），实际%d", len(nearby))
	}

	// 查找距离3以内的玩家
	nearby2 := nt.PlayersWithinDistance(center, 3.0)
	if len(nearby2) != 1 {
		t.Errorf("期望找到1个玩家（距离<=3），实际%d", len(nearby2))
	}
}

func TestNearbyTracker_Clear(t *testing.T) {
	lookup := func(uuid [16]byte) (string, bool) {
		return "Player", true
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	// 生成玩家
	uuid := [16]byte{1}
	tracker.SpawnEntity(1, "minecraft:player", uuid, entity.Position{}, entity.Vector3{})

	time.Sleep(50 * time.Millisecond)

	if nt.Count() != 1 {
		t.Errorf("期望初始玩家数量为1，实际%d", nt.Count())
	}

	// 清空
	nt.Clear()

	if nt.Count() != 0 {
		t.Errorf("期望清空后玩家数量为0，实际%d", nt.Count())
	}

	// 验证GetNearbyPlayers返回空
	players := nt.GetNearbyPlayers()
	if len(players) != 0 {
		t.Errorf("期望GetNearbyPlayers为空，实际%d", len(players))
	}
}

func TestNearbyTracker_UnknownPlayer(t *testing.T) {
	// 查找返回未找到
	lookup := func(uuid [16]byte) (string, bool) {
		return "", false
	}

	tracker := entity.NewTracker()
	nt := NewNearbyTracker(tracker, lookup)

	var enterCalled bool
	nt.SetCallbacks(PlayerCallbacks{
		OnPlayerEnter: func(p *NearbyPlayer) {
			enterCalled = true
		},
	})

	// 生成玩家实体（但lookup找不到玩家信息）
	uuid := [16]byte{1}
	tracker.SpawnEntity(1, "minecraft:player", uuid, entity.Position{}, entity.Vector3{})

	time.Sleep(50 * time.Millisecond)

	if !enterCalled {
		t.Error("期望OnPlayerEnter回调被调用（即使找不到玩家信息）")
	}

	// 验证玩家被添加但用户名为空
	players := nt.GetNearbyPlayers()
	if len(players) != 1 {
		t.Fatalf("期望找到1个玩家，实际%d", len(players))
	}

	if players[0].Username != "" {
		t.Errorf("期望Username为空字符串，实际%s", players[0].Username)
	}
}
