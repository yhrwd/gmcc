package key

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/crypto/scrypt"
)

var (
	machineFingerprintOnce sync.Once
	machineFingerprint     []byte
	machineFingerprintErr  error
)

// Manager 密钥管理器
type Manager struct {
	cachedFingerprint []byte
}

// NewManager 创建密钥管理器
func NewManager() *Manager {
	return &Manager{}
}

// GetMachineFingerprint 获取机器指纹
// 组合多个硬件标识生成唯一指纹
func (m *Manager) GetMachineFingerprint() ([]byte, error) {
	machineFingerprintOnce.Do(func() {
		m.cachedFingerprint, machineFingerprintErr = m.computeFingerprint()
	})
	return machineFingerprint, machineFingerprintErr
}

// computeFingerprint 计算机器指纹
func (m *Manager) computeFingerprint() ([]byte, error) {
	var components []string

	// 获取操作系统信息
	components = append(components, runtime.GOOS)
	components = append(components, runtime.GOARCH)

	// 尝试获取CPU信息
	if cpuInfo, err := m.getCPUInfo(); err == nil && cpuInfo != "" {
		components = append(components, cpuInfo)
	}

	// 尝试获取主机名
	if hostname, err := m.getHostname(); err == nil && hostname != "" {
		components = append(components, hostname)
	}

	// 组合并哈希
	combined := strings.Join(components, "|")
	hash := sha256.Sum256([]byte(combined))
	return hash[:], nil
}

// getCPUInfo 获取CPU信息（平台特定）
func (m *Manager) getCPUInfo() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return m.getWindowsCPUInfo()
	case "linux":
		return m.getLinuxCPUInfo()
	case "darwin":
		return m.getDarwinCPUInfo()
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// getWindowsCPUInfo 获取Windows CPU信息
func (m *Manager) getWindowsCPUInfo() (string, error) {
	// 使用wmic命令获取CPU信息
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorId", "/value")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ProcessorId=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "ProcessorId=")), nil
		}
	}
	return "", fmt.Errorf("ProcessorId not found")
}

// getLinuxCPUInfo 获取Linux CPU信息
func (m *Manager) getLinuxCPUInfo() (string, error) {
	// 读取/proc/cpuinfo
	cmd := exec.Command("grep", "^Serial", "/proc/cpuinfo")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		parts := strings.Split(string(output), ":")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	// 尝试使用lscpu
	cmd = exec.Command("lscpu")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "CPU(s):") || strings.HasPrefix(line, "Model name:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("CPU info not found")
}

// getDarwinCPUInfo 获取macOS CPU信息
func (m *Manager) getDarwinCPUInfo() (string, error) {
	cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getHostname 获取主机名
func (m *Manager) getHostname() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// DeriveKey 派生加密密钥
// 从机器指纹和密码派生密钥
func (m *Manager) DeriveKey(password string, salt []byte, n, r, p, keyLen int) ([]byte, error) {
	fingerprint, err := m.GetMachineFingerprint()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine fingerprint: %w", err)
	}

	// 组合机器指纹和密码
	combined := append(fingerprint, []byte(password)...)

	// 使用scrypt派生密钥
	key, err := scrypt.Key(combined, salt, n, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("scrypt key derivation failed: %w", err)
	}

	return key, nil
}

// GetKeyGetter 返回一个KeyGetter函数
func (m *Manager) GetKeyGetter(defaultN, defaultR, defaultP, defaultKeyLen int) func(password string) ([]byte, error) {
	return func(password string) ([]byte, error) {
		// 使用固定的盐值（实际应该每次存储不同的盐值）
		// 这里返回原始组合的哈希，在Vault中再用scrypt派生
		fingerprint, err := m.GetMachineFingerprint()
		if err != nil {
			return nil, err
		}

		// 组合机器指纹和密码，计算SHA256
		combined := append(fingerprint, []byte(password)...)
		hash := sha256.Sum256(combined)
		return hash[:], nil
	}
}

// String 返回机器指纹的十六进制字符串表示
func (m *Manager) String() string {
	fingerprint, err := m.GetMachineFingerprint()
	if err != nil {
		return ""
	}
	return hex.EncodeToString(fingerprint)
}
