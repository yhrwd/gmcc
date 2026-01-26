package protocol

import (
	"errors"
	"net"
	"strconv"
	"fmt"
	"gmcc/internal/logger"
)

// 默认端口
const DefaultMCPort = "25565"

func SRVCheck(name string) (targetHost string, targetPort string, err error) {
	service := "minecraft"
	proto := "tcp"

	_, addrs, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return "", "", err
	}

	if len(addrs) == 0 {
		return "", "", errors.New("no SRV records found")
	}
	bestSRV := addrs[0]
	for _, srv := range addrs {
		if srv.Priority < bestSRV.Priority || (srv.Priority == bestSRV.Priority && srv.Weight > bestSRV.Weight) {
			bestSRV = srv
		}
	}

	return bestSRV.Target, strconv.Itoa(int(bestSRV.Port)), nil
}

func CheckPortInAddr(addr string) (host string, port string, hasPort bool, err error) {
	host, port, err = net.SplitHostPort(addr)
	if err != nil {
		var addrErr *net.AddrError
		if ok := errors.As(err, &addrErr); ok && addrErr.Err == "missing port in address" {
			return addr, "", false, nil
		}
		return "", "", false, err
	}
	return host, port, true, nil
}

func CheckAddr(addr string) (address string, err error) {
	host, port, hasPort, err := CheckPortInAddr(addr)
	if err != nil {
		return "", fmt.Errorf("服务器地址格式错误: %w", err)
	}

	if hasPort {
		return net.JoinHostPort(host, port), nil
	}

	srvHost, srvPort, err := SRVCheck(host)
	if err != nil {
		logger.Info("SRV解析失败，将使用默认端口25565: %v\n", err)
		return net.JoinHostPort(host, DefaultMCPort), nil
	}

	return net.JoinHostPort(srvHost, srvPort), nil
}
