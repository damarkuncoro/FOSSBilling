package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	low = "abcdefghijklmnopqrstuvwxyz"
	upp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	dig = "0123456789"
	spc = "!@#$%^&*()-_=+"
)

func GeneratePassword(length int, spec bool) (string, error) {
	if length < 8 { length = 8 }
	set := low + upp + dig
	if spec { set += spc }
	req := []string{low, upp, dig}; if spec { req = append(req, spc) }

	res := make([]byte, length)
	for i, c := range req {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(c))))
		res[i] = c[n.Int64()]
	}
	for i := len(req); i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		res[i] = set[n.Int64()]
	}
	for i := len(res) - 1; i > 0; i-- {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := n.Int64(); res[i], res[j] = res[j], res[i]
	}
	return string(res), nil
}

func GenerateRandomString(l int) string {
	cs := low + upp + dig; b := make([]byte, l)
	for i := range b { n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(cs)))); b[i] = cs[n.Int64()] }
	return string(b)
}

func FormatBytes(b int64) string {
	if b < 1024 { return fmt.Sprintf("%d B", b) }
	d, e := int64(1024), 0
	for n := b / 1024; n >= 1024; n /= 1024 { d *= 1024; e++ }
	return fmt.Sprintf("%.2f %s", float64(b)/float64(d), []string{"KB", "MB", "GB", "TB"}[e])
}

func NormalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }
