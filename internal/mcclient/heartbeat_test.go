package mcclient

import (
	"net"
	"testing"
	"time"

	"gmcc/internal/config"
	"gmcc/internal/mcclient/packet"
	"gmcc/internal/mcclient/protocol"
)

func TestSendAFKHeartbeatUsesMovePlayerPosAndTickEnd(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client := New(&config.Config{
		Account: config.AccountConfig{PlayerID: "TestPlayer"},
	})
	client.state = protocol.StatePlay
	client.conn = packet.NewPacketConn(clientConn)
	client.lastAFKPacket = time.Time{}
	client.Player.UpdatePosition(12.5, 64, -3.25, 0, 0, 0)

	done := make(chan error, 1)
	go func() {
		done <- client.sendAFKHeartbeatIfNeeded()
	}()

	serverPacketConn := packet.NewPacketConn(serverConn)

	pkt, err := serverPacketConn.ReadPacket()
	if err != nil {
		t.Fatalf("ReadPacket(move_player_pos) error = %v", err)
	}
	if pkt.ID != protocol.PlayServerMovePlayerPos {
		t.Fatalf("move packet id = 0x%02X, want 0x%02X", pkt.ID, protocol.PlayServerMovePlayerPos)
	}
	if len(pkt.Data) != 25 {
		t.Fatalf("move payload len = %d, want 25", len(pkt.Data))
	}

	pkt, err = serverPacketConn.ReadPacket()
	if err != nil {
		t.Fatalf("ReadPacket(client_tick_end) error = %v", err)
	}
	if pkt.ID != protocol.PlayServerClientTickEnd {
		t.Fatalf("tick_end id = 0x%02X, want 0x%02X", pkt.ID, protocol.PlayServerClientTickEnd)
	}

	if err := <-done; err != nil {
		t.Fatalf("sendAFKHeartbeatIfNeeded() error = %v", err)
	}
}
