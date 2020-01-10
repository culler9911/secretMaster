package main

import (
	"fmt"

	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
	"github.com/syndtr/goleveldb/leveldb"
)

//go:generate cqcfg -c .
// cqp: 名称: SecretMaster
// cqp: 版本: 0.0.1:1
// cqp: 作者: molin
// cqp: 简介: 一个超棒的Go语言插件Demo，它会回复你的私聊消息~专为诡秘之主粉丝序列群开发的小游戏
func main() { /*此处应当留空*/ }

func init() {
	cqp.AppID = "me.cqp.molin.secretmaster" // TODO: 修改为这个插件的ID
	cqp.PrivateMsg = onPrivateMsg
	cqp.GroupMsg = onGroupMsg
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	cqp.SendPrivateMsg(fromQQ, msg) //复读机
	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	info := cqp.GetGroupMemberInfo(fromGroup, fromQQ, true)

	cqp.SendGroupMsg(fromGroup, "@"+info.Name+" "+msg) //复读机
	cqp.SendGroupMsg(fromGroup, "name:"+info.Name)
	cqp.SendGroupMsg(fromGroup, "card:"+info.Card)
	cqp.SendGroupMsg(fromGroup, fmt.Sprintf("%+v", info))

	db, err := leveldb.OpenFile("secret.db", nil)

	db.Put([]byte("key"), []byte("value"), nil)

	data, err := db.Get([]byte("key"), nil)

	cqp.SendGroupMsg(fromGroup, string(data))

	defer db.Close()

	if err != nil {
		return 1
	}

	return 0
}
