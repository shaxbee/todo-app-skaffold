package dbtest

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func containerName(driver string) string {
	suffix, _ := rand.Int(rand.Reader, big.NewInt(100000))
	return fmt.Sprintf("%s-%s", driver, suffix)
}
