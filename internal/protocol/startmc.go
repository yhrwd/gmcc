package protocol

import (
	"gmcc/internal/logger"
	"gmcc/internal/protocol/connection"
	"gmcc/internal/protocol/stage"
	"net"
	"strconv"
	"time"
)

func Startmc(target string) {
	// 1. 解析地址 (例如 "mc.example.com:25565")
	host, portStr, err := net.SplitHostPort(target)
	if err != nil {
		// 如果没有端口，默认 25565 (看你需求是否支持 SRV 解析，这里简单处理)
		host = target
		portStr = "25565"
	}
	port, _ := strconv.Atoi(portStr)

	// 2. 建立 TCP 连接
	logger.Infof("连接到 %s...", target)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, portStr), 5*time.Second)
	if err != nil {
		logger.Errorf("连接失败: %v", err)
		return
	}
	
	// 3. 初始化状态机
	// NewStateMachine 内部会调用 NewConn(conn)，从而初始化 bufio
	sm := connection.NewStateMachine(conn)

	// 4. 注册状态 (填入解析后的 IP 和 端口)
	sm.Register(&stage.ClientHandshakeState{
		Addr: host,
		Port: uint16(port),
	})
	sm.Register(&stage.ClientStatusState{})

	// 5. 启动状态机
	if err := sm.Start(); err != nil {
		logger.Errorf("状态机启动失败: %v", err)
		return // 启动失败直接返回
	}

	// 6. 阻塞等待
	// 直到连接断开或 Context 被 Cancel
	<-sm.Context().Done()
	logger.Info("连接已断开")
}
