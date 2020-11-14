package ipsearch

// Decimal to integer.
// Returns characters consumed, success.
func dtoi(s string) (i int, ok bool) {
	n := 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n > 0xFF {
			// Stop here and return true to be able to
			// handle cases like "192.168.1.355".
			return i, true
		}
	}
	if i == 0 {
		// There are no digits in string.
		return 0, false
	}
	return i, true
}

// Hexadecimal to integer.
// Returns characters consumed, success.
func xtoi(s string) (i int, ok bool) {
	n := 0
	for i = 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n > 0xFFFF {
			// Stop here and return true to be able to
			// handle cases like "1:2::3:4:abcde".
			return i, true
		}
	}
	if i == 0 {
		// There are no digits in string.
		return i, false
	}
	return i, true
}

func isdec(c byte) bool {
	return c >= '0' && c <= '9'
}

func ishex(c byte) bool {
	return isdec(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
