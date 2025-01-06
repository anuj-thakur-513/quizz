package utils

func ConvertToString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
