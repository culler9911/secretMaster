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

	newChurch := &ChurchInfo{strs[1], strs[2], nil, fromQQ, b.CurrentNick, uint64(enterPrice), b.getMoney(fromQQ) / 200, 1, b.getMoney(fromQQ), 1}

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
	fmt.Println(b.getMoney(fromQQ))

	churchs.ChurchList = append(churchs.ChurchList, newChurch)
	b.setGroupValue("Churchs", churchs)

	info += fmt.Sprintf("创建成功！\n名称:%s\n介绍:%s\n尊神/教主:%s\n技能:%s\n入会费:%d\n最大人数:%d\n等级:%d级\n注册资本:%d金镑",
		newChurch.Name, newChurch.Commit, newChurch.CreatorNick, strs[3], newChurch.Money, b.getMoney(fromQQ)/200, newChurch.Level, b.getMoney(fromQQ))

	return info
}

func (b *Bot) deleteChurch(fromQQ uint64, msg string) string {
	strs := strings.Split(msg, ";")
	if len(strs) != 2 {
		return "格式不正确" + fmt.Sprintf("%+v", strs)
	}

	churchs := b.getGroupValue("Churchs", &Churchs{}).(*Churchs)
	fmt.Printf("churchs: %+v", churchs)
	for i, c := range churchs.ChurchList {
		if c == nil || (c.Name == strs[1] && c.CreatorQQ == fromQQ) {
			if len(churchs.ChurchList) > 1 {
				churchs.ChurchList[i] = churchs.ChurchList[len(churchs.ChurchList)-1]
				churchs.ChurchList = churchs.ChurchList[:len(churchs.ChurchList)-1]
			} else {
				churchs.ChurchList = nil
			}
			b.setGroupValue("Churchs", churchs)
		}
	}

	return "解散成功"
}

func (b *Bot) listChurch() string {
	info := "\n"
	churchs := b.getGroupValue("Churchs", &Churchs{}).(*Churchs)

	for i, c := range churchs.ChurchList {
		update := false

		money := b.getMoney(c.CreatorQQ)
		if int64(money)-int64(c.CreateMoney) > 0 {
			c.Level = (money - c.CreateMoney) / 10000
			if c.Level > 3 {
				c.Level = 3
			}

			if c.Level > uint64(len(c.Skills)) {
				add := c.Level - uint64(len(c.Skills))
				for i, s := range skillList {
					have := false
					for _, skill := range c.Skills {
						if s.ID == skill.ID {
							have = true
						}
					}
					if !have {
						c.Skills = append(c.Skills, &skillList[i])
						add--
						if add == 0 {
							break
						}
					}
				}
				update = true
			}
		}
		p := b.getPersonFromDb(c.CreatorQQ)
		skillStr := ""
		for i := 0; i < len(c.Skills); i++ {
			skillStr += fmt.Sprintf("%s lv%d; ", c.Skills[i].Name, p.SecretLevel-4)
		}

		info += fmt.Sprintf("\n名称:%s\n介绍:%s\n尊神/教主:%s\n技能:%s\n入会费:%d\n人数:%d/%d\n等级:%d级\n注册资本:%d金镑\n",
			c.Name, c.Commit, c.CreatorNick, skillStr, c.Money, c.Members, b.getMoney(c.CreatorQQ)/200, c.Level, c.CreateMoney)

		if update {
			churchs.ChurchList[i] = c
			b.setGroupValue("Churchs", churchs)
		}
	}

	return info
}
