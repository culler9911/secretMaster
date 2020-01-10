package secret

import "testing"
import "fmt"

// "github.com/molin0000/secretMaster/rlp"
import "github.com/molin0000/secretMaster/rlp"

func TestRun(t *testing.T) {
	b := NewSecretBot(333, 333, "cat")
	fmt.Println(b.Run("@cat 帮助", 111))
}

func TestRlp(t *testing.T) {
	p := &Person{
		Group:       333,
		QQ:          333,
		Name:        "aa",
		JoinTime:    123,
		LastChat:    1234,
		SecretID:    0,
		SecretLevel: 0,
		ChatCount:   100,
	}

	buf, err := rlp.EncodeToBytes(p)
	fmt.Println(buf, err)

	var a Person
	err = rlp.DecodeBytes(buf, &a)

	fmt.Printf("%+v", a)
	fmt.Println(err)
}
