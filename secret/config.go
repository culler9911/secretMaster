package secret

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
)

func isGroupAdmin(group, qq uint64) bool {
	member := cqp.GetGroupMemberInfo(int64(group), int64(qq), true)
	return member.Auth > 1
}

func isGroupMaster(group, qq uint64) bool {
	member := cqp.GetGroupMemberInfo(int64(group), int64(qq), true)
	return member.Auth > 2
}

func (b *Bot) setMaster(fromQQ uint64, msg string) string {
	cfg := b.getGroupValue("Config", &Config{})
	if cfg.(*Config).HaveMaster && !isGroupMaster(b.Group, fromQQ) {
		return "Master已经配置，只有群主可以修改"
	}

	if !isGroupAdmin(b.Group, fromQQ) {
		return "只有群主或管理员可以设置.master"
	}

	str0 := strings.Split(msg, "@")
	msg = str0[0]

	strs := strings.Split(msg, ";")
	if len(strs) != 2 {
		return fmt.Sprintf("当前master:%+v", cfg)
	}

	masterQQ, _ := strconv.ParseUint(strs[1], 10, 64)
	cfg.(*Config).MasterQQ = masterQQ
	cfg.(*Config).HaveMaster = true
	if masterQQ == 0 {
		cfg.(*Config).HaveMaster = false
	}
	b.setGroupValue("Config", cfg)
	return fmt.Sprintf("成功设置%d为插件master", masterQQ)
}

func (b *Bot) isMaster(fromQQ uint64) bool {
	cfg := b.getGroupValue("Config", &Config{})
	return cfg.(*Config).HaveMaster && (fromQQ == cfg.(*Config).MasterQQ)
}

func (b *Bot) notGM() string {
	return "对不起，你不是GM，别想欺骗机器人"
}

func (b *Bot) moneyUpdate(fromQQ uint64, update bool) string {
	if !b.isMaster(fromQQ) {
		return b.notGM()
	}

	bind := b.getMoneyBind()
	if update {
		if bind.HasUpdate {
			return "已经升级过了，请不要重复升级"
		}
		bind.HasUpdate = true
		b.setMoneyBind(bind)
		return "升级成功"
	}

	bind.HasUpdate = false
	b.setMoneyBind(bind)
	return "降级成功，配置保留，但不再读取ini文件"
}

func (b *Bot) moneyMap(fromQQ uint64, msg string) string {
	if !b.isMaster(fromQQ) {
		return b.notGM()
	}

	str0 := strings.Split(msg, "@")
	msg = str0[0]

	strs := strings.Split(msg, ";")
	bind := &MoneyBind{}
	bind.IniPath = strs[1]
	bind.IniSection = strs[2]
	bind.IniKey = strs[3]
	b.setMoneyBind(bind)
	return fmt.Sprintf("映射成功, Path:%s, Section:%s, Key:%s\n", strs[1], strs[2], strs[3])
}

func (b *Bot) gmCmd(fromQQ uint64, msg string) string {
	if !b.isMaster(fromQQ) {
		return b.notGM()
	}

	str0 := strings.Split(msg, "@")
	msg = str0[0]

	strs := strings.Split(msg, ";")

	n1, err1 := strconv.Atoi(strs[2])
	n2, err2 := strconv.ParseUint(strs[3], 10, 64)

	if err1 != nil || err2 != nil {
		return fmt.Sprintf("参数解析错误: 0:%s, 1:%s, 2:%s, 3:%s, %+v, %+v", strs[0], strs[1], strs[2], strs[3], err1, err2)
	}
	switch strs[1] {
	case "money":
		m := b.getMoneyFromDb(uint64(n2), 0)
		if n1 > 0 {
			m.Money += uint64(n1)
		} else {
			m.Money -= uint64(-1 * n1)
		}
		b.setMoneyToDb(uint64(n2), m)
		return fmt.Sprintf("%d 金镑：%d", n2, n1)
	case "exp":
		m := b.getPersonFromDb(uint64(n2))
		if n1 > 0 {
			m.ChatCount += uint64(n1)
		} else {
			m.ChatCount -= uint64(-1 * n1)
		}
		b.setPersonToDb(uint64(n2), m)
		return fmt.Sprintf("%d 经验：%d", n2, n1)
	case "magic":
		m := b.getExternFromDb(uint64(n2))
		if n1 > 0 {
			m.Magic += uint64(n1)
		} else {
			m.Magic -= uint64(-1 * n1)
		}
		b.setExternToDb(uint64(n2), m)
		return fmt.Sprintf("%d 灵性：%d", n2, n1)
	default:
		return "参数解析错误"
	}
}

func (b *Bot) botSwitch(fromQQ uint64, enable bool) string {
	if !b.isMaster(fromQQ) {
		return b.notGM()
	}

	b.setSwitch(enable)
	if enable {
		return fmt.Sprintf("已在群%d开启《序列战争》诡秘之主背景小游戏插件。", b.Group)
	}

	return fmt.Sprintf("已在群%d关闭《序列战争》诡秘之主背景小游戏插件。", b.Group)
}
