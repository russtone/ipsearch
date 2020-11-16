// Package ipsearch provides function to search
// IPv4 and IPv6 addresses in string.
package ipsearch

import "net"

const (
	maxIPv4OctetLen = 3
	maxIPv6OctetLen = 4
)

// Find returns positions of all found IPv4 and IPv6 addreses
// in the given string.
func Find(s string) [][2]int {
	res := make([][2]int, 0)
	last := 0

outer:
	for i := 1; i < len(s); {
		switch s[i] {
		case '.':
			// Go back to find possible start of
			// the IPv4 address.
			n := 0
			for j := 1; j <= min(maxIPv4OctetLen, i-last); j++ {
				if !isdec(s[i-j]) {
					break
				}
				n++
			}

			// There are no decimal digits
			// before ".".
			if n == 0 {
				break
			}

			// Try all possible starts.
			for k := n; k > 0; k-- {
				cc, ok := checkIPv4(s[i-k:])
				if ok {
					res = append(res, [2]int{i - k, i - k + cc})
					last = i - k + cc
					i += cc - k
					continue outer
				}
			}
		case ':':
			// Go back to find possible start of
			// the IPv6 address.
			n := 0
			for j := 1; j <= min(maxIPv6OctetLen, i-last); j++ {
				if !ishex(s[i-j]) {
					break
				}
				n++
			}

			// n can be equal to 0 here,
			// because IPv6 might start with "::".
			for k := n; k >= 0; k-- {
				cc, ok := checkIPv6(s[i-k:])
				if ok {
					res = append(res, [2]int{i - k, i - k + cc})
					last = i - k + cc
					i += cc - k
					continue outer
				}
			}
		}
		i++
	}

	return res
}

// checkIPv4 checks if the given string starts with
// valid IPv4 address.
// Returns characters consumed, success.
func checkIPv4(s string) (int, bool) {
	cc := 0
	for i := 0; i < net.IPv4len; i++ {
		if len(s) == 0 {
			// End of string, missing octets.
			return cc, false
		}

		if i > 0 {
			if s[0] != '.' {
				return cc, false
			}
			// Consume "."
			s = s[1:]
			cc++
		}

		c, ok := dtoi(s)
		if !ok {
			return cc, false
		}

		// Consume decimal digits.
		s = s[c:]
		cc += c
	}

	return cc, true
}

// checkIPv6 checks if the given string starts with
// valid IPv6 address.
// Returns characters consumed, success.
func checkIPv6(s string) (int, bool) {
	ellipsis := -1 // position of ellipsis in ip
	cc := 0        // characters consumed

	// Might have leading ellipsis
	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0

		// Consume "::"
		s = s[2:]
		cc += 2

		// Might be only ellipsis
		if len(s) == 0 || !ishex(s[0]) {
			return cc, true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i := 0
	for i < net.IPv6len {
		// Hex number.
		c, ok := xtoi(s)
		if !ok {
			// Handle case with semicolon
			// at the and "1:2::3:xxxxx".
			cc-- // do not include last ":"
			break
		}

		// If followed by dot, might be in trailing IPv4.
		if c < len(s) && s[c] == '.' {
			if ellipsis < 0 && i != net.IPv6len-net.IPv4len {
				// Not the right place.
				return cc, false
			}
			if n, ok := checkIPv4(s); ok {
				i += net.IPv4len
				cc += n
				break
			} else {
				return cc, false
			}
		}

		// Consume hex.
		s = s[c:]
		cc += c
		i += 2

		// Stop at end of string
		// or if there are enough bytes.
		if len(s) == 0 ||
			i == net.IPv6len ||
			ellipsis >= 0 && i == net.IPv6len-2 {
			break
		}

		// Otherwise must be followed by colon and more.
		if s[0] != ':' || len(s) == 1 {
			break
		}

		// Consume ":"
		s = s[1:]
		cc++

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				// Handle case "1::1::1".
				cc-- // do not include last ":"
				break
			}
			ellipsis = i
			// Consume ":"
			s = s[1:]
			cc++

			// End of string of enough bytes.
			if len(s) == 0 || i == net.IPv6len-2 {
				break
			}
		}
	}

	// If didn't parse enough, expand ellipsis.
	if i < net.IPv6len {
		if ellipsis < 0 {
			return cc, false
		}
	}
	return cc, true
}
