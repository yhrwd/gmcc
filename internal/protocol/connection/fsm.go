package connection

import (
	"context"
	"fmt"
	"net"
	"sync"
	
	"gmcc/internal/logger"
)

// ===== 连接状态枚举 =====

// ConnState 当前状态
type ConnState int

const (
	StateHandshake ConnState = iota
	StateStatus
	StateLogin
	StateTransfer
	StateConfiguration
	StatePlay
)

func (s ConnState) String() string {
	switch s {
	case StateHandshake:
		return "HANDSHAKE"
	case StateStatus:
		return "STATUS"
	case StateLogin:
		return "LOGIN"
	case StateTransfer:
		return "TRANSFER"
	case StateConfiguration:
		return "CONFIGURATION"
	case StatePlay:
		return "PLAY"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// ===== 核心接口定义 =====

// Temporary occupation
type IPlayer interface {
	Name() string
	UUID() string
}

// Temporary occupation
type IWorld interface {
	Name() string
}

// ConnContext 共享状态/信息
type ConnContext struct {
    // 底层连接与状态
    conn   *Conn        // 底层网络连接
    state  ConnState    // 当前连接状态
    sm     *StateMachine // 状态机实例
    encryptStream cipher.Stream

    // 连接元数据
    addr                string // 服务器地址
    protocolVersion     int    // 协议版本号
    compressionThreshold int   // 压缩阈值 (-1 为禁用)
    sharedSecret        []byte // 加密密钥 (登录阶段)

    // 游戏实体引用
    player IPlayer // 当前玩家
    world  IWorld  // 当前世界

    // 生命周期管理
    ctx    context.Context    // 用于控制连接退出
    cancel context.CancelFunc // 取消函数

    // 并发控制锁
    readMu  sync.Mutex // 保护读操作 (SetReadDeadline + Read)
    writeMu sync.Mutex // 保护写操作 (发包必须串行)
    
    // 状态切换保护
    stateMu   sync.Mutex // 保护 state 和 switching 的访问
    switching bool       // 是否正在执行状态切换 (防递归/并发)
}


// 一个状态必须实现的行为
type State interface {
	// Name 返回该状态对应的 ConnState 枚举值
	Name() ConnState

	// Enter 当状态机切换到该状态时调用
	// 注意：Enter 不应阻塞（除非它在内部启动了 goroutine），否则会卡住状态机
	Enter(ctx *ConnContext) error

	// HandlePacket 处理该状态下收到的数据包
	HandlePacket(ctx *ConnContext, packet any) error

	// Exit 当状态机离开该状态时调用
	// 用于清理资源或停止 goroutine
	Exit(ctx *ConnContext) error
}

// ===== 状态机实现 =====

// StateMachine 管理连接的状态流转
type StateMachine struct {
	ctx *ConnContext          // 连接上下文

	mu     sync.RWMutex       // 读写锁：读多写少场景下比 Mutex 性能更好
	states map[ConnState]State // 状态注册表
}

// NewStateMachine 创建一个新的状态机实例
func NewStateMachine(rawConn net.Conn) *StateMachine {
	ctx, cancel := context.WithCancel(context.Background())

	// 这里使用 NewConn 包装原始连接，初始化 bufio.Reader/Writer
	wrappedConn := NewConn(rawConn) 

	connCtx := &ConnContext{
		Conn:   wrappedConn,
		State:  StateHandshake,
		Ctx:    ctx,
		Cancel: cancel,
        ProtocolVersion: 773, // 默认1.20.10
	}

	sm := &StateMachine{
		ctx:    connCtx,
		states: make(map[ConnState]State),
	}
	connCtx.SM = sm
	return sm
}

// Register 注册一个状态到状态机
func (sm *StateMachine) Register(state State) {
	if state == nil {
		panic("Registering nil state is not allowed")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	stateName := state.Name()
	sm.states[stateName] = state
	logger.Infof("Registered state: %s", stateName)
}

// Start 启动状态机，进入初始状态
func (sm *StateMachine) Start() error {
	sm.mu.RLock()
	initialState, ok := sm.states[sm.ctx.State]
	sm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("initial state %s is not registered", sm.ctx.State)
	}

	logger.Infof("Starting state machine, entering initial state: %s", sm.ctx.State)
	return initialState.Enter(sm.ctx)
}

// Switch 切换到下一个状态
// 这是核心方法，负责状态的平滑过渡
func (sm *StateMachine) Switch(nextState ConnState) error {
    sm.ctx.writeMu.Lock()         // 用 writeMu 当全局切换锁（因为发包也用它）
    defer sm.ctx.writeMu.Unlock()

    if sm.ctx.switching {
        return fmt.Errorf("state switch already in progress")
    }
    sm.ctx.switching = true
    defer func() { sm.ctx.switching = false }()

    sm.mu.RLock()
    current := sm.ctx.State
    oldHandler, hasOld := sm.states[current]
    newHandler, hasNew := sm.states[nextState]
    sm.mu.RUnlock()

    if !hasNew || newHandler == nil {
        return fmt.Errorf("target state %s not registered", nextState)
    }

    logger.Infof("Switching: %s → %s", current, nextState)

    // 退出旧状态
    if hasOld && oldHandler != nil {
        if err := oldHandler.Exit(sm.ctx); err != nil {
            logger.Errorf("exit %s error: %v", current, err)
            // 不 return，继续尝试进入新状态
        }
    }

    // 原子更新状态
    sm.mu.Lock()
    sm.ctx.State = nextState
    sm.mu.Unlock()

    // 进入新状态
    if err := newHandler.Enter(sm.ctx); err != nil {
        logger.Infof("enter %s error: %v", nextState, err)
        sm.Close()
        return err
    }

    return nil
}

// HandlePacket 分发数据包给当前状态处理
func (sm *StateMachine) HandlePacket(packet any) error {
	// 快速读取当前状态处理器
	sm.mu.RLock()
	currentState := sm.ctx.State
	handler, ok := sm.states[currentState]
	sm.mu.RUnlock()

	if !ok || handler == nil {
		return fmt.Errorf("no handler found for current state: %s", currentState)
	}

	// 调用具体的处理逻辑
	return handler.HandlePacket(sm.ctx, packet)
}

// Close 关闭状态机及相关资源
func (sm *StateMachine) Close() {
	logger.Infof("Closing state machine (Current state: %s)", sm.ctx.State)

	// 1. 触发 Context 取消，通知所有 goroutine 退出
	if sm.ctx.Cancel != nil {
		sm.ctx.Cancel()
	}

	// 2. 调用当前状态的 Exit 方法进行清理
	sm.mu.RLock()
	currentState := sm.ctx.State
	handler, ok := sm.states[currentState]
	sm.mu.RUnlock()

	if ok && handler != nil {
		if err := handler.Exit(sm.ctx); err != nil {
			logger.Infof("Error during cleanup of state %s: %v", currentState, err)
		}
	}

	// 3. 关闭底层连接
	if sm.ctx.Conn != nil {
		_ = sm.ctx.Conn.Close()
	}
}

// GetCurrentState 返回当前状态
func (sm *StateMachine) GetCurrentState() ConnState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.ctx.State
}

func (sm *StateMachine) Context() context.Context {
    return sm.ctx.Ctx
}
