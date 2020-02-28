package secret

import (
	"fmt"
	"testing"
)

func TestCreateChurch(t *testing.T) {
	fromQQ := uint64(67939461000)
	b := NewSecretBot(3334, 333, "aaa", false, &debugInteract{})
	fmt.Println(b.Run("[CQ:at,qq=3334] 自杀", fromQQ, "mm"))

	fmt.Println(b.Run("[CQ:at,qq=3334] 创建", fromQQ, "mm"))

	fmt.Println(b.Run("[CQ:at,qq=3334] 创建;aaa;bbb;ccc;ddd", fromQQ, "mm"))
	fmt.Println(b.Run("[CQ:at,qq=3334] 创建;aaa;bbb;ccc;200", fromQQ, "mm"))
	p := b.getPersonFromDb(fromQQ)
	p.SecretLevel = 5
	b.setPersonToDb(fromQQ, p)
	fmt.Println(b.Run("[CQ:at,qq=3334] 创建;aaa;bbb;ccc;200", fromQQ, "mm"))

	b.setMoney(fromQQ, 2000)
	fmt.Println(b.Run("[CQ:at,qq=3334] 创建;aaa;bbb;ccc;200", fromQQ, "mm"))
	b.buyMace(fromQQ)
	fmt.Println(b.getItems(fromQQ))
	fmt.Println(b.Run("[CQ:at,qq=3334] 创建;aaa;bbb;ccc;100", fromQQ, "mm"))

}
