package service

import (
	"strconv"
	"strings"
)

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate") || strings.Contains(message, "unique constraint") || strings.Contains(message, "constraint failed")
}

func ParseUintParam(value string) (uint, error) {
	id64, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id64 == 0 {
		return 0, BadRequest("无效的 ID")
	}
	return uint(id64), nil
}
