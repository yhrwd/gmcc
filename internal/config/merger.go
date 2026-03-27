package config

import (
	"fmt"
	"reflect"
	"strings"
)

// ConfigChange 表示配置文件中的一个变更
type ConfigChange struct {
	Path     string // 字段路径，如 "actions.delay_ms"
	OldValue any    // 原值（可能为空）
	NewValue any    // 新值
}

// ConfigMerger 负责合并配置结构体
type ConfigMerger struct{}

// MergeWithDefault 将当前配置与默认配置合并，保留用户自定义值
// 返回合并后的配置和变更记录
func (m *ConfigMerger) MergeWithDefault(current *Config) (*Config, []ConfigChange, error) {
	if current == nil {
		return nil, nil, fmt.Errorf("当前配置为空")
	}

	defaultCfg := Default()
	changes := make([]ConfigChange, 0)

	_, err := m.mergeStruct(current, defaultCfg, "", &changes)
	if err != nil {
		return nil, nil, fmt.Errorf("合并配置失败: %w", err)
	}

	// 直接返回 current 的副本 with changes applied
	// 由于我们修改了 current 的底层结构，直接返回它
	return current, changes, nil
}

// mergeStruct 递归合并两个结构体
func (m *ConfigMerger) mergeStruct(current, defaultVal any, path string, changes *[]ConfigChange) (any, error) {
	currentVal := reflect.ValueOf(current)
	defaultValOf := reflect.ValueOf(defaultVal)

	// 处理指针类型
	if currentVal.Kind() == reflect.Ptr {
		if currentVal.IsNil() {
			return defaultVal, nil
		}
		currentVal = currentVal.Elem()
	}

	if defaultValOf.Kind() == reflect.Ptr {
		if defaultValOf.IsNil() {
			return current, nil
		}
		defaultValOf = defaultValOf.Elem()
	}

	// 如果是接口类型，获取其动态值
	if currentVal.Kind() == reflect.Interface {
		if currentVal.IsNil() {
			return defaultVal, nil
		}
		currentVal = currentVal.Elem()
	}

	if defaultValOf.Kind() == reflect.Interface {
		if defaultValOf.IsNil() {
			return current, nil
		}
		defaultValOf = defaultValOf.Elem()
	}

	// 处理基本类型
	if m.isBasicType(currentVal.Kind()) {
		return m.mergeBasicType(current, defaultVal, path, changes)
	}

	// 处理结构体
	if currentVal.Kind() == reflect.Struct {
		return m.mergeStructFields(currentVal, defaultValOf, path, changes)
	}

	// 处理其他类型（切片、map等）
	return current, nil
}

// isBasicType 检查是否为基本类型
func (m *ConfigMerger) isBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}

// mergeBasicType 合并基本类型字段
func (m *ConfigMerger) mergeBasicType(current, defaultVal any, path string, changes *[]ConfigChange) (any, error) {
	// 如果当前值为零值，使用默认值
	if m.isZeroValue(current) {
		change := ConfigChange{
			Path:     path,
			OldValue: current,
			NewValue: defaultVal,
		}
		*changes = append(*changes, change)
		return defaultVal, nil
	}
	return current, nil
}

// mergeStructFields 合并结构体字段
func (m *ConfigMerger) mergeStructFields(currentVal, defaultVal reflect.Value, path string, changes *[]ConfigChange) (any, error) {
	// 直接修改当前结构体（通过指针）
	if currentVal.CanAddr() {
		_ = currentVal.Addr().Interface()

		// 遍历结构体字段
		for i := 0; i < currentVal.NumField(); i++ {
			field := currentVal.Type().Field(i)
			if !field.IsExported() {
				continue
			}

			yamlTag := field.Tag.Get("yaml")
			if yamlTag == "" {
				continue // 没有yaml标签的字段跳过
			}

			fieldPath := path
			if fieldPath != "" {
				fieldPath += "."
			}
			fieldPath += yamlTag

			currentFieldValue := currentVal.Field(i)
			defaultFieldValue := defaultVal.Field(i)

			mergedValue, err := m.mergeStruct(
				currentFieldValue.Interface(),
				defaultFieldValue.Interface(),
				fieldPath,
				changes,
			)
			if err != nil {
				return nil, fmt.Errorf("合并字段 %s 失败: %w", fieldPath, err)
			}

			// 如果有合并后的值且当前字段可以设置
			if mergedValue != nil && currentFieldValue.CanSet() {
				mergedVal := reflect.ValueOf(mergedValue)
				// 处理基础类型
				if mergedVal.Type().ConvertibleTo(currentFieldValue.Type()) {
					currentFieldValue.Set(mergedVal.Convert(currentFieldValue.Type()))
				}
			}
		}
	}

	return currentVal.Interface(), nil
}

// isZeroValue 检查值是否为零值
func (m *ConfigMerger) isZeroValue(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

// needsUpdate 检查配置是否需要更新
func (m *ConfigMerger) needsUpdate(current *Config) (bool, error) {
	if current == nil {
		return true, nil
	}

	_, foundChanges, err := m.MergeWithDefault(current)
	if err != nil {
		return false, fmt.Errorf("检查更新失败: %w", err)
	}

	needsUpdate := len(foundChanges) > 0
	return needsUpdate, nil
}

// GetFieldPath 获取字段的配置路径
func (m *ConfigMerger) GetFieldPath(structType reflect.Type, fieldIndex int) string {
	field := structType.Field(fieldIndex)
	yamlTag := field.Tag.Get("yaml")
	if yamlTag == "" {
		return ""
	}

	// 处理嵌套结构体的路径
	parts := strings.Split(yamlTag, ",")
	if len(parts) > 0 {
		return parts[0]
	}
	return yamlTag
}
