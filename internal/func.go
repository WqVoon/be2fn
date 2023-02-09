package internal

// x 在 s 代表的整数切片中
func InIntSlice(x string, s []int) Unit {
	return func(vals Kv) (bool, error) {
		xVal, err := vals.GetInt(x)
		if err != nil {
			return false, err
		}

		for _, val := range s {
			if val == xVal {
				return true, nil
			}
		}
		return false, nil
	}
}

// x 在 s 代表的字符串切片中
func InStrSlice(x string, s []string) Unit {
	return func(vals Kv) (bool, error) {
		xVal, err := vals.GetString(x)
		if err != nil {
			return false, err
		}

		for _, val := range s {
			if val == xVal {
				return true, nil
			}
		}
		return false, nil
	}
}
