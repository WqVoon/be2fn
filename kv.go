package be2fn

// 所有的入参
type Kv map[string]interface{}

// 尝试获取字符串类型
func (vars Kv) GetString(key string) (string, bool) {
	val, ok := vars[key].(string)
	return val, ok
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
func (vars Kv) GetInt64(key string) (int64, bool) {
	val, ok := vars[key].(int64)
	return val, ok
}

// 尝试获取整数类型，如果失败返回 defaultVal
func (vars Kv) GetInt64OrDefault(key string, defaultVal int64) int64 {
	val, ok := vars[key].(int64)
	if !ok {
		return defaultVal
	}
	return val
}
