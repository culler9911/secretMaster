package main

import (
	"fmt"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
	"github.com/molin0000/secretMaster/secret"
)

//go:generate cqcfg -c .
// cqp: 名称: SecretMaster
// cqp: 版本: 0.0.3:1
// cqp: 作者: molin
// cqp: 简介: 专为诡秘之主粉丝序列群开发的小游戏
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

	selfQQ := cqp.GetLoginQQ()

	selfInfo := cqp.GetGroupMemberInfo(fromGroup, selfQQ, false)

	bot := secret.NewSecretBot(uint64(cqp.GetLoginQQ()), uint64(fromGroup), selfInfo.Name)

	ret := bot.Run(msg, uint64(fromQQ), info.Name)

	if len(ret) > 0 {
		fmt.Printf("\nSend group msg:%d, %s\n", fromGroup, ret)
		time.Sleep(1000)
		id := cqp.SendGroupMsg(fromGroup, "@"+info.Name+" "+ret)
		fmt.Printf("\nSend finish id:%d\n", id)
		// fmt.Println("private ret:", cqp.SendPrivateMsg(fromQQ, ret))
	}

	return 0
}
