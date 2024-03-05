package config

import (
	"fmt"

	"try-on/internal/pkg/utils"
)

var JsonLogFormat = func() string {
	values := []string{"time", "status", "latency", "ip", "method", "path", "error"}

	result := utils.Reduce(values, func(first, second string) string {
		return first + fmt.Sprintf(`%s: "${%s}", `, second, second)
	})

	return "{" + result[:len(result)-2] + "}\n"
}()
