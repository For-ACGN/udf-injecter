package injecter

import (
	"math/rand"
	"time"
)

func RandomStr(size int) string {
	if size < 1 {
		size = 8
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	result := make([]rune, size)
	for i := 0; i < size; i++ {
		// after space
		n := 33 + r.Intn(90)
		switch {
		case n > 47 && n < 58: //  48-57 number
		case n > 64 && n < 91: //  65-90 A-Z
		case n > 96 && n < 123: // 97-122 a-z
		default:
			i--
			continue
		}
		result[i] = rune(n)
	}
	return string(result)
}
