package secret

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/molin0000/secretMaster/rlp"
)

func TestRun(t *testing.T) {
	b := NewSecretBot(3334, 333, "cat")
	fmt.Println(b.Run("@cat 属性", 112, "mm"))
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

	type test1 struct {
		Age  uint64
		Name string
	}

	c := &test1{
		Age:  1000,
		Name: "molin",
	}

	buf, _ = rlp.EncodeToBytes(c)

	type test2 struct {
		Age   uint64
		Name  string
		Luck  uint64
		HP    uint64
		Name2 string
	}

	var d test2
	rlp.DecodeBytes(buf, &d)
	fmt.Printf("%+v", d)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().Unix())

	fmt.Println(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
}

func TestSprinf(t *testing.T) {
	v := "asdf更换途径3sdfasdfaf"
	fmt.Println(v)

	// var s = "MemTotal: 1001332 kB"
	var valid = regexp.MustCompile("[0-9]")
	fmt.Println(fmt.Sprint(valid.FindAllStringSubmatch(v, -1)))
	r := valid.FindAllStringSubmatch(v, -1)

	a, err := strconv.Atoi(r[0][0])
	fmt.Println(a, err)

}

func TestGoRoute(t *testing.T) {

	a := "sdfadsfadsfadf尊名摸摸摸摸摸"

	b := a[strings.Index(a, "尊名")+len("尊名"):]
	fmt.Println(b)

	var eCh chan string
	eCh = make(chan string, 100)
	// eCh <- "hello"
	select {
	case v := <-eCh:
		fmt.Println(v)
	default:
		fmt.Println("default")
	}

	fmt.Println("Finish")
}
