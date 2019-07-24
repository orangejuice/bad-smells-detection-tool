package util

import "regexp"

func IsIpAddress(host string) bool {
	if host == "localhost" {
		return true
	} else if matched, _ := regexp.MatchString(`((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}`, host); matched {
		return true
	}
	return false
}
