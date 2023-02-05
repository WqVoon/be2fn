package internal

import (
	"fmt"
)

// 所有的入参
type Kv map[string]interface{}

// 尝试获取字符串类型
func (vars Kv) GetString(key string) (string, error) {
	val, ok := vars[key].(string)
	if !ok {
		return val, fmt.Errorf("failed to get string by key(%s)", key)
	}
	return val, nil
}

// 尝试获取字符串类型，如果失败返回 defaultVal
func (vars Kv) GetStringOrDefault(key, defaultVal string) string {
	val, ok := vars[key].(string)
	if !ok {
		return defaultVal
	}
	return val
}

// 尝试获取整数类型
func (vars Kv) GetInt(key string) (int, error) {
	val, ok := vars[key].(int)
	if !ok {
		return val, fmt.Errorf("failed to get int by key(%s)", key)
	}
	return val, nil
}

// 尝试获取整数类型，如果失败返回 defaultVal
func (vars Kv) GetInt64OrDefault(key string, defaultVal int) int {
	val, ok := vars[key].(int)
	if !ok {
		return defaultVal
	}
	return val
}

// 尝试获取布尔类型
func (vars Kv) GetBool(key string) (bool, error) {
	val, ok := vars[key].(bool)
	if !ok {
		return val, fmt.Errorf("failed to get bool by key(%s)", key)
	}
	return val, nil
}

// 尝试获取布尔类型，如果失败返回 defaultVal
func (vars Kv) GetBoolOrDefault(key string, defaultVal bool) bool {
	val, ok := vars[key].(bool)
	if !ok {
		return defaultVal
	}
	return val
}
