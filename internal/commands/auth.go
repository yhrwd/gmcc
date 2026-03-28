package commands

import (
	"strings"
	"sync"
)

type AuthManager struct {
	whitelist map[string]bool
	allowAll  bool
	mu        sync.RWMutex
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		whitelist: make(map[string]bool),
		allowAll:  false,
	}
}

func (a *AuthManager) SetAllowAll(allowAll bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.allowAll = allowAll
}

func (a *AuthManager) SetWhitelist(list []string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.whitelist = make(map[string]bool)
	for _, entry := range list {
		entry = strings.TrimSpace(entry)
		if entry != "" {
			a.whitelist[entry] = true
		}
	}
}

func (a *AuthManager) AddToWhitelist(entries ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry != "" {
			a.whitelist[entry] = true
		}
	}
}

func (a *AuthManager) RemoveFromWhitelist(entries ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		delete(a.whitelist, entry)
	}
}

func (a *AuthManager) Check(msg *Message) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.allowAll {
		return true
	}

	if msg.SenderUUID != "" && a.whitelist[msg.SenderUUID] {
		return true
	}

	if msg.Sender != "" {
		for entry := range a.whitelist {
			if strings.EqualFold(entry, msg.Sender) {
				return true
			}
		}
	}

	return false
}

func (a *AuthManager) IsAllowed(name, uuid string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.allowAll {
		return true
	}

	if uuid != "" && a.whitelist[uuid] {
		return true
	}

	if name != "" {
		for entry := range a.whitelist {
			if strings.EqualFold(entry, name) {
				return true
			}
		}
	}

	return false
}

func (a *AuthManager) GetWhitelist() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	list := make([]string, 0, len(a.whitelist))
	for entry := range a.whitelist {
		list = append(list, entry)
	}
	return list
}

func (a *AuthManager) IsAllowAll() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.allowAll
}
