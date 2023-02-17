package util

func GetValueInt(value int, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func GetValueInt64(value int64, defaultValue int64) int64 {
	if value == 0 {
		return defaultValue
	}
	return value
}
