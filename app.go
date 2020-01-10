package main

import "github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
import "strconv"

//go:generate cqcfg -c .
// cqp: 名称: SecretMaster
// cqp: 版本: 1.0.0:0
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
	cqp.SendGroupMsg(fromGroup, "@"+strconv.FormatInt(fromQQ, 10) + " " + msg) //复读机
	return 0
}