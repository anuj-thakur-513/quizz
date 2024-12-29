package core

func ApiResponse(status int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": status < 400,
		"message": message,
		"data":    data,
	}
}
