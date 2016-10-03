package broadcaster

import (
	"crypto/rand"
	"fmt"
)

func MustUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("err: rand.Read failed for reason %s", err.Error()))
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
