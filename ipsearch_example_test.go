package ipsearch_test

import (
	"fmt"

	"github.com/russtone/ipsearch"
)

func ExampleFind() {
	s := `
github.com has address 140.82.121.4
google.com has address 173.194.222.101
yandex.ru has address 77.88.55.50
yandex.ru has IPv6 address 2a02:6b8:a::a
`
	rr := ipsearch.Find(s)

	for _, r := range rr {
		fmt.Println(s[r[0]:r[1]])
	}

	// Output:
	// 140.82.121.4
	// 173.194.222.101
	// 77.88.55.50
	// 2a02:6b8:a::a
}
