/*
 * @Author: Liu Sainan
 * @Date: 2024-01-15 23:38:41
 */

package validateutils

import (
	"myadmin/internal/utils/netutils"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// letterDigits 字符串以字母开头
func Digits(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	// 使用正则表达式匹配以字母开头的字符串
	match, _ := regexp.MatchString("^[0-9]*$", name)
	return match
}

// letterDigits 字符串以字母开头
func LetterDigits(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	// 使用正则表达式匹配以字母开头的字符串
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]*$", name)
	return match
}

// ValidateIPPort 验证IP:PORT格式
func ValidateIPPort(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	ipPort := strings.Split(fl.Field().String(), ":")
	if err := netutils.IsIPv4(ipPort[0]); err != nil {
		if netutils.IsHostName(ipPort[0]) {
			return true
		}
		// logger.Warningf("主机名不可访问")
		return false
	}
	if len(ipPort) > 1 {
		port, err := strconv.Atoi(ipPort[1])
		if err != nil {
			return false
		}
		if port <= 0 || port >= 65536 {
			return false
		}
	}
	return true
}
