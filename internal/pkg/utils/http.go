package utils

func HttpOk(code int) bool {
	return code >= 200 && code < 300
}
