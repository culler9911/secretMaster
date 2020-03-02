package secret

import (
	"fmt"
	"testing"
)

func TestGetProperty(t *testing.T) {
	fromQQ := uint64(121)
	b := NewSecretBot(1234, 4567, "aa", false, &debugInteract{})

	fmt.Println(b.getProperty(fromQQ))
}
