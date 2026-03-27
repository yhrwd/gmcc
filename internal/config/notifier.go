package config

import (
	"gmcc/internal/i18n"
	"gmcc/internal/logx"
)

// UpdateNotifier 配置更新通知接口
type UpdateNotifier interface {
	NotifyConfigUpdate(changes []ConfigChange) error
}

// TUIUpdateNotifier TUI模式的通知实现
type TUIUpdateNotifier struct{}

// NotifyConfigUpdate 在TUI模式下通知配置更新
func (n *TUIUpdateNotifier) NotifyConfigUpdate(changes []ConfigChange) error {
	i18n := i18n.GetI18n()

	// 获取翻译文本
	title := i18n.Translate("config.updated.title", len(changes))
	detail := i18n.Translate("config.updated.detail", len(changes))

	// 记录通知信息
	logx.Infof(title)
	logx.Infof(detail)

	// 记录详细变更
	if len(changes) > 0 {
		logx.Infof(i18n.Translate("config.updated.fields_added"))
		for _, change := range changes {
			logx.Infof("  - %s: %v", change.Path, change.NewValue)
		}
	}

	return nil
}

// HeadlessUpdateNotifier Headless模式的通知实现
type HeadlessUpdateNotifier struct{}

// NotifyConfigUpdate 在Headless模式下通知配置更新
func (n *HeadlessUpdateNotifier) NotifyConfigUpdate(changes []ConfigChange) error {
	i18n := i18n.GetI18n()

	// 使用logx系统记录通知
	logx.Infof(i18n.Translate("config.updated.headless", len(changes)))

	// 记录添加的配置项路径
	for _, change := range changes {
		logx.Infof(i18n.Translate("config.field.added", change.Path, change.NewValue))
	}

	return nil
}

// GetNotifier 根据运行模式获取对应的通知器
func GetNotifier(isHeadless bool) UpdateNotifier {
	if isHeadless {
		return &HeadlessUpdateNotifier{}
	}
	return &TUIUpdateNotifier{}
}
