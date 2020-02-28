package secret

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Bot) createChurch(fromQQ uint64, msg string) string {
	strs := strings.Split(msg, ";")
	if len(strs) != 5 {
		return "\n指令格式不正确！别想欺骗机器人！" + fmt.Sprintf("%+v", strs)
	}

	enterPrice, err := strconv.Atoi(strs[4])

	if err != nil {
		return "\n指令格式不正确！别想欺骗机器人！" + fmt.Sprintf("%+v %+v", strs, err)
	}

	p := b.getPersonFromDb(fromQQ)

	if p.SecretLevel < 5 {
		return "\n你的能力和声望不足以服众，没人愿意听你说什么。"
	}

	if b.getMoney(fromQQ) < 1000 {
		return "\n你查了查贝克兰德的地价，发现钱似乎还不够。"
	}

	newChurch := &ChurchInfo{strs[1], strs[2], nil, fromQQ, uint64(enterPrice), b.getMoney(fromQQ) / 200, 1, b.getMoney(fromQQ)}

	info := "\n"

	churchs := b.getGroupValue("Churchs", &Churchs{}).(*Churchs)

	for _, church := range churchs.ChurchList {
		if church.CreatorQQ == newChurch.CreatorQQ || church.Name == newChurch.Name {
			return "请不要重复创建。"
		}
	}

	findSkill := false
	for _, skill := range skillList {
		if strs[3] == skill.Name {
			newChurch.Skills = append(newChurch.Skills, &skill)
			findSkill = true
			break
		}
	}

	if !findSkill {
		return "\n抱歉，没听说过这个技能。"
	}

	if !b.useItem(fromQQ, "至高权杖") {
		return "\n你似乎缺少了什么信物，没法吸引信徒。"
	}

	b.setMoney(fromQQ, -1000)

	churchs.ChurchList = append(churchs.ChurchList, newChurch)
	b.setGroupValue("Churchs", churchs)

	info += fmt.Sprintf("创建成功！\n名称:%s\n介绍:%s\n尊神/教主:%s\n技能:%s\n入会费:%d\n最大人数:%d\n等级:%d级",
		newChurch.Name, newChurch.Commit, b.CurrentNick, strs[3], newChurch.Money, b.getMoney(fromQQ)/200, newChurch.Level)

	return info
}
