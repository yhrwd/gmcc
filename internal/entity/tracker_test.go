package entity

import (
	"testing"
	"time"
)

func TestEntity_IsPlayer(t *testing.T) {
	e := &Entity{
		ID:   1,
		Type: "minecraft:player",
	}
	if !e.IsPlayer() {
		t.Error("预期玩家类型返回true")
	}

	e2 := &Entity{
		ID:   2,
		Type: "minecraft:zombie",
	}
	if e2.IsPlayer() {
		t.Error("预期非玩家类型返回false")
	}
}

func TestPosition_DistanceTo(t *testing.T) {
	p1 := Position{X: 0, Y: 0, Z: 0}
	p2 := Position{X: 3, Y: 4, Z: 0}

	dist := p1.DistanceTo(p2)
	if dist != 25 {
		t.Errorf("预期距离25，实际%f", dist)
	}
}

func TestTracker_SpawnEntity(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	var spawnCalled bool
	var spawnedEntity *Entity
	tracker.SetCallbacks(Callbacks{
		OnSpawn: func(e *Entity) {
			spawnCalled = true
			spawnedEntity = e
		},
	})

	uuid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pos := Position{X: 1.0, Y: 2.0, Z: 3.0}
	vel := Vector3{X: 0, Y: 0, Z: 0}

	e := tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	// 等待回调
	time.Sleep(50 * time.Millisecond)

	if e == nil {
		t.Fatal("期望实体不为nil")
	}

	if e.ID != 1 {
		t.Errorf("期望ID=1，实际%d", e.ID)
	}

	if e.Type != "minecraft:player" {
		t.Errorf("期望Type=minecraft:player，实际%s", e.Type)
	}

	if e.Position.X != 1.0 || e.Position.Y != 2.0 || e.Position.Z != 3.0 {
		t.Errorf("期望位置(1,2,3)，实际(%f,%f,%f)", e.Position.X, e.Position.Y, e.Position.Z)
	}

	if !spawnCalled {
		t.Error("期望OnSpawn回调被调用")
	}

	if spawnedEntity == nil || spawnedEntity.ID != 1 {
		t.Error("期望回调中实体ID正确")
	}

	// 测试通过ID获取
	e2, ok := tracker.Get(1)
	if !ok || e2 == nil {
		t.Error("期望能通过ID获取实体")
	}

	// 测试通过UUID获取
	e3, ok := tracker.GetByUUID(uuid)
	if !ok || e3 == nil {
		t.Error("期望能通过UUID获取实体")
	}
}

func TestTracker_UpdatePosition(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	var moveCalled bool
	var movedEntity *Entity
	var oldPos Position
	tracker.SetCallbacks(Callbacks{
		OnMove: func(e *Entity, old Position) {
			moveCalled = true
			movedEntity = e
			oldPos = old
		},
	})

	// 先创建实体
	uuid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pos := Position{X: 0, Y: 0, Z: 0}
	vel := Vector3{X: 0, Y: 0, Z: 0}
	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond) // 等待spawn回调

	// 更新位置
	newPos := Position{X: 10.0, Y: 20.0, Z: 30.0}
	tracker.UpdatePosition(1, newPos)

	// 等待100ms去重后回调
	time.Sleep(150 * time.Millisecond)

	if !moveCalled {
		t.Error("期望OnMove回调被调用")
	}

	if movedEntity == nil || movedEntity.ID != 1 {
		t.Error("期望回调中实体正确")
	}

	if oldPos.X != 0 || oldPos.Y != 0 || oldPos.Z != 0 {
		t.Errorf("期望旧位置(0,0,0)，实际(%f,%f,%f)", oldPos.X, oldPos.Y, oldPos.Z)
	}

	// 验证实体位置已更新
	e, _ := tracker.Get(1)
	if e.Position.X != 10.0 || e.Position.Y != 20.0 || e.Position.Z != 30.0 {
		t.Errorf("期望新位置(10,20,30)，实际(%f,%f,%f)", e.Position.X, e.Position.Y, e.Position.Z)
	}
}

func TestTracker_UpdatePositionDelta(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	// 先创建实体
	uuid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pos := Position{X: 100.0, Y: 64.0, Z: 200.0}
	vel := Vector3{X: 0, Y: 0, Z: 0}
	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	// deltaX = 4096 * 1.0 = 4096
	// deltaY = 4096 * 0.5 = 2048
	// deltaZ = 4096 * (-1.0) = -4096
	tracker.UpdatePositionDelta(1, 4096, 2048, -4096)

	time.Sleep(150 * time.Millisecond)

	// 验证位置
	e, _ := tracker.Get(1)
	expectedX := 101.0
	expectedY := 64.5
	expectedZ := 199.0

	if e.Position.X != expectedX {
		t.Errorf("期望X=%f，实际%f", expectedX, e.Position.X)
	}
	if e.Position.Y != expectedY {
		t.Errorf("期望Y=%f，实际%f", expectedY, e.Position.Y)
	}
	if e.Position.Z != expectedZ {
		t.Errorf("期望Z=%f，实际%f", expectedZ, e.Position.Z)
	}
}

func TestTracker_RemoveEntity(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	var removeCalled bool
	var removedEntity *Entity
	tracker.SetCallbacks(Callbacks{
		OnRemove: func(e *Entity) {
			removeCalled = true
			removedEntity = e
		},
	})

	// 创建实体
	uuid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pos := Position{X: 0, Y: 0, Z: 0}
	vel := Vector3{X: 0, Y: 0, Z: 0}
	tracker.SpawnEntity(1, "minecraft:player", uuid, pos, vel)

	time.Sleep(50 * time.Millisecond)

	// 移除实体
	tracker.RemoveEntity(1)

	time.Sleep(50 * time.Millisecond)

	if !removeCalled {
		t.Error("期望OnRemove回调被调用")
	}

	if removedEntity == nil || removedEntity.ID != 1 {
		t.Error("期望回调中实体正确")
	}

	// 验证已移除
	_, ok := tracker.Get(1)
	if ok {
		t.Error("期望实体已被移除")
	}

	if tracker.Count() != 0 {
		t.Errorf("期望实体数量为0，实际%d", tracker.Count())
	}
}

func TestTracker_ByType(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	// 创建不同类型的实体
	uuid1 := [16]byte{1}
	uuid2 := [16]byte{2}
	uuid3 := [16]byte{3}
	pos := Position{X: 0, Y: 0, Z: 0}
	vel := Vector3{X: 0, Y: 0, Z: 0}

	tracker.SpawnEntity(1, "minecraft:player", uuid1, pos, vel)
	tracker.SpawnEntity(2, "minecraft:zombie", uuid2, pos, vel)
	tracker.SpawnEntity(3, "minecraft:player", uuid3, pos, vel)

	players := tracker.ByType("minecraft:player")
	if len(players) != 2 {
		t.Errorf("期望找到2个玩家，实际%d", len(players))
	}

	zombies := tracker.ByType("minecraft:zombie")
	if len(zombies) != 1 {
		t.Errorf("期望找到1个僵尸，实际%d", len(zombies))
	}
}

func TestTracker_All(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	uuid1 := [16]byte{1}
	uuid2 := [16]byte{2}
	pos := Position{X: 0, Y: 0, Z: 0}
	vel := Vector3{X: 0, Y: 0, Z: 0}

	tracker.SpawnEntity(1, "minecraft:player", uuid1, pos, vel)
	tracker.SpawnEntity(2, "minecraft:zombie", uuid2, pos, vel)

	all := tracker.All()
	if len(all) != 2 {
		t.Errorf("期望所有实体数量为2，实际%d", len(all))
	}
}

func TestTracker_RemoveEntities(t *testing.T) {
	tracker := NewTracker()
	defer tracker.Stop()

	uuid1 := [16]byte{1}
	uuid2 := [16]byte{2}
	uuid3 := [16]byte{3}
	pos := Position{X: 0, Y: 0, Z: 0}
	vel := Vector3{X: 0, Y: 0, Z: 0}

	tracker.SpawnEntity(1, "minecraft:player", uuid1, pos, vel)
	tracker.SpawnEntity(2, "minecraft:zombie", uuid2, pos, vel)
	tracker.SpawnEntity(3, "minecraft:creeper", uuid3, pos, vel)

	// 批量移除
	tracker.RemoveEntities([]int32{1, 2})

	if tracker.Count() != 1 {
		t.Errorf("期望剩余1个实体，实际%d", tracker.Count())
	}

	_, ok := tracker.Get(3)
	if !ok {
		t.Error("期望实体3仍然存在")
	}
}
