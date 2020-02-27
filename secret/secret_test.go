package secret

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/molin0000/secretMaster/rlp"
	"gopkg.in/ini.v1"
)

func TestRun(t *testing.T) {
	b := NewSecretBot(3334, 333, "cat", false, &debugInteract{})
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

func TestIniRead(t *testing.T) {
	iniPath := "/Users/molin/coolq/data/app/com.qmt.demo/玩家数据.ini"

	cfg, err := ini.Load(iniPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}

	fmt.Println("Data Path:", cfg.Section(cString("67939461")).Key(cString("金币数量")).String())
	amount, err := cfg.Section(cString("67939461")).Key(cString("数量")).Uint64()
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}
	amount += 10000
	cfg.Section(cString("67939461")).Key(cString("金币数量")).SetValue(strconv.FormatUint(amount, 10))
	cfg.SaveTo(iniPath)
}

func TestUint64(t *testing.T) {
	n2, err := strconv.Atoi("3459914053")
	fmt.Println(n2, err)
}
